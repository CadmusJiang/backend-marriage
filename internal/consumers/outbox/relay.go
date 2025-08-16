package outbox

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/xinliangnote/go-gin-api/configs"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
	repoOutbox "github.com/xinliangnote/go-gin-api/internal/repository/mysql/outbox"
	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// StartRelay publishes outbox events to Redis Streams reliably.
func StartRelay(ctx context.Context, logger *zap.Logger, db mysql.Repo) {
	cfg := configs.Get().Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Pass,
		DB:           cfg.Db,
		MaxRetries:   cfg.MaxRetries,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	go func() {
		defer func() { _ = rdb.Close() }()
		ticker := time.NewTicker(300 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}

			// Pull a small batch with SKIP LOCKED to avoid thundering herd.
			var batch []repoOutbox.Event
			err := db.GetDbW().Session(&gorm.Session{Logger: gormlogger.Default.LogMode(gormlogger.Silent)}).Transaction(func(tx *gorm.DB) error {
				// lock rows for this transaction
				if err := tx.Raw("SELECT id, topic, payload, status, retry_count, created_at, published_at FROM outbox_events WHERE status = 0 ORDER BY id ASC LIMIT 100 FOR UPDATE SKIP LOCKED").Scan(&batch).Error; err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				logger.Warn("outbox.select", zap.Error(err))
				continue
			}
			if len(batch) == 0 {
				continue
			}

			// Publish one by one; update status in a transaction batch
			_ = db.GetDbW().Session(&gorm.Session{Logger: gormlogger.Default.LogMode(gormlogger.Silent)}).Transaction(func(tx *gorm.DB) error {
				now := time.Now()
				for i := range batch {
					evt := &batch[i]
					// payload must be JSON object containing at least type/recordId
					var m map[string]interface{}
					if err := json.Unmarshal([]byte(evt.Payload), &m); err != nil {
						// mark as published to avoid poison
						_ = tx.Model(&repoOutbox.Event{}).Where("id = ?", evt.ID).Updates(map[string]interface{}{"status": 1, "published_at": now}).Error
						continue
					}
					topic := evt.Topic
					// XADD
					if _, err := rdb.XAdd(&redis.XAddArgs{Stream: topic, MaxLenApprox: 1_000_000, Values: map[string]interface{}{"event": string(evt.Payload)}}).Result(); err != nil {
						_ = tx.Model(&repoOutbox.Event{}).Where("id = ?", evt.ID).UpdateColumn("retry_count", gorm.Expr("retry_count + 1")).Error
						logger.Warn("outbox.publish", zap.Uint64("id", evt.ID), zap.Error(err))
						continue
					}
					if err := tx.Model(&repoOutbox.Event{}).Where("id = ?", evt.ID).Updates(map[string]interface{}{"status": 1, "published_at": now}).Error; err != nil {
						logger.Warn("outbox.ack", zap.Uint64("id", evt.ID), zap.Error(err))
						continue
					}
				}
				return nil
			})
		}
	}()
}
