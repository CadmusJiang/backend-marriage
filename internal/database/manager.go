package database

import (
	"fmt"

	"github.com/xinliangnote/go-gin-api/internal/proposal/tablesqls"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
	"go.uber.org/zap"
)

// Manager 数据库管理器
type Manager struct {
	db     mysql.Repo
	logger *zap.Logger
}

// New 创建数据库管理器
func New(db mysql.Repo, logger *zap.Logger) *Manager {
	return &Manager{
		db:     db,
		logger: logger,
	}
}

// TestConnection 测试数据库连接
func (m *Manager) TestConnection() error {
	// 测试读连接
	if err := m.db.GetDbR().Raw("SELECT 1").Error; err != nil {
		return fmt.Errorf("读数据库连接测试失败: %v", err)
	}

	// 测试写连接
	if err := m.db.GetDbW().Raw("SELECT 1").Error; err != nil {
		return fmt.Errorf("写数据库连接测试失败: %v", err)
	}

	return nil
}

// EnsureTables 确保核心业务相关的数据表存在
func (m *Manager) EnsureTables() error {
	type table struct {
		name string
		sql  string
	}
	tables := []table{
		{name: "org", sql: tablesqls.CreateOrgTableSql()},
		{name: "org_history", sql: tablesqls.CreateOrgHistoryTableSql()},
		{name: "account", sql: tablesqls.CreateAccountTableSql()},
		{name: "account_history", sql: tablesqls.CreateAccountHistoryTableSql()},
		{name: "account_org_relation", sql: tablesqls.CreateAccountOrgRelationTableSql()},
		{name: "customer_authorization_record", sql: tablesqls.CreateCustomerAuthorizationRecordTableSql()},
		{name: "customer_authorization_record_history", sql: tablesqls.CreateCustomerAuthorizationRecordHistoryTableSql()},
		{name: "cooperation_store", sql: tablesqls.CreateCooperationStoreTableSql()},
		{name: "cooperation_store_history", sql: tablesqls.CreateCooperationStoreHistoryTableSql()},
		{name: "outbox_events", sql: tablesqls.CreateOutboxTableSql()},
	}

	for _, t := range tables {
		m.logger.Info("检查并创建表", zap.String("table", t.name))

		// 检查表是否存在
		var exists bool
		if err := m.db.GetDbW().Raw("SELECT COUNT(*) > 0 FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?", t.name).Scan(&exists).Error; err != nil {
			return fmt.Errorf("检查表 %s 是否存在失败: %v", t.name, err)
		}

		if !exists {
			// 创建表
			m.logger.Info("正在创建表...", zap.String("table", t.name))
			if err := m.db.GetDbW().Exec(t.sql).Error; err != nil {
				return fmt.Errorf("创建表 %s 失败: %v", t.name, err)
			}
			m.logger.Info("表创建完成", zap.String("table", t.name))
		} else {
			m.logger.Debug("表已存在", zap.String("table", t.name))
		}
	}

	return nil
}

// RebuildTables 强制重建数据库表（删除所有表并重新创建）
func (m *Manager) RebuildTables() error {
	// 先删除所有表
	m.logger.Info("正在删除所有现有表...")
	tablesToDrop := []string{
		"cooperation_store_history",
		"customer_authorization_record_history",
		"account_history",
		"org_history",
		"account_org_relation",
		"customer_authorization_record",
		"cooperation_store",
		"account",
		"org",
		"outbox_events",
	}

	for _, tableName := range tablesToDrop {
		m.logger.Info("删除表", zap.String("table", tableName))
		if err := m.db.GetDbW().Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", tableName)).Error; err != nil {
			m.logger.Error("删除表失败", zap.String("table", tableName), zap.Error(err))
		}
	}
	m.logger.Info("所有表删除完成")

	// 重新创建表
	m.logger.Info("正在重新创建表...")
	if err := m.EnsureTables(); err != nil {
		return fmt.Errorf("重新创建表失败: %v", err)
	}

	// 插入测试数据
	m.logger.Info("正在插入测试数据...")
	if err := m.ReinsertMockData(); err != nil {
		return fmt.Errorf("插入测试数据失败: %v", err)
	}

	m.logger.Info("数据库重建完成")
	return nil
}

// ReinsertMockData 重新插入mock数据
func (m *Manager) ReinsertMockData() error {
	// 定义需要重新插入的表和对应的SQL（按依赖顺序执行）
	type tdata struct {
		name   string
		insert string
	}
	items := []tdata{
		{name: "org", insert: tablesqls.CreateOrgTableDataSql()},
		{name: "org_history", insert: tablesqls.CreateOrgHistoryTableDataSql()},
		{name: "account", insert: tablesqls.CreateAccountTableDataSql()},
		{name: "account_history", insert: tablesqls.CreateAccountHistoryTableDataSql()},
		{name: "account_org_relation", insert: tablesqls.CreateAccountOrgRelationTableDataSql()},
		{name: "customer_authorization_record", insert: tablesqls.CreateCustomerAuthorizationRecordTableDataSql()},
		{name: "customer_authorization_record_history", insert: tablesqls.CreateCustomerAuthorizationRecordHistoryTableDataSql()},
		{name: "cooperation_store", insert: tablesqls.CreateCooperationStoreTableDataSql()},
		{name: "cooperation_store_history", insert: tablesqls.CreateCooperationStoreHistoryTableDataSql()},
	}

	for _, it := range items {
		m.logger.Info("清空表数据", zap.String("table", it.name))
		res := m.db.GetDbW().Exec(fmt.Sprintf("DELETE FROM %s", it.name))
		if err := res.Error; err != nil {
			return fmt.Errorf("清空表 %s 失败: %v", it.name, err)
		}
		m.logger.Info("已清空", zap.String("table", it.name), zap.Int64("rows_affected", res.RowsAffected))

		if it.insert == "" {
			m.logger.Info("无插入SQL，跳过", zap.String("table", it.name))
			continue
		}

		m.logger.Info("插入 mock 数据", zap.String("table", it.name))
		res = m.db.GetDbW().Exec(it.insert)
		if err := res.Error; err != nil {
			return fmt.Errorf("插入表 %s 的mock数据失败: %v", it.name, err)
		}
		m.logger.Info("插入完成", zap.String("table", it.name), zap.Int64("rows_affected", res.RowsAffected))
	}

	return nil
}
