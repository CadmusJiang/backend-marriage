package outbox

import (
	"time"
)

// Event represents an outbox event row
type Event struct {
	ID          uint64     `gorm:"column:id;primaryKey;autoIncrement"`
	Topic       string     `gorm:"column:topic;type:varchar(128);not null"`
	Payload     []byte     `gorm:"column:payload;type:json;not null"`
	Status      int8       `gorm:"column:status;type:tinyint;not null;default:0"`
	RetryCount  int        `gorm:"column:retry_count;type:int;not null;default:0"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime"`
	PublishedAt *time.Time `gorm:"column:published_at"`
}

func (Event) TableName() string {
	return "outbox_events"
}
