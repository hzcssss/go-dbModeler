package connector

import (
	"database/sql"
)

// Connector 定义数据库连接器接口
type Connector interface {
	// Connect 连接到数据库
	Connect() (*sql.DB, error)

	// GetDatabases 获取所有数据库
	GetDatabases() ([]string, error)

	// GetTables 获取指定数据库中的所有表
	GetTables(database string) ([]string, error)

	// GetTableMetadata 获取表的元数据信息
	GetTableMetadata(database, table string) (*TableMetadata, error)

	// Close 关闭数据库连接
	Close() error
}

// ConnectionConfig 表示数据库连接配置
type ConnectionConfig struct {
	Type     string // 数据库类型：MySQL, PostgreSQL, SQLite
	Host     string // 主机名或IP地址
	Port     string // 端口号
	Username string // 用户名
	Password string // 密码
	Database string // 数据库名（可选）
}

// TableMetadata 表示表的元数据
type TableMetadata struct {
	Name    string      // 表名
	Fields  []FieldInfo // 字段信息
	Indexes []IndexInfo // 索引信息
}

// FieldInfo 表示字段信息
type FieldInfo struct {
	Name       string // 字段名
	Type       string // 数据库类型
	Length     int    // 长度（如果适用）
	IsNullable bool   // 是否可为空
	IsPrimary  bool   // 是否为主键
	IsUnique   bool   // 是否唯一
	Default    string // 默认值
	Comment    string // 注释
}

// IndexInfo 表示索引信息
type IndexInfo struct {
	Name    string   // 索引名
	Type    string   // 索引类型（PRIMARY, UNIQUE, INDEX, FULLTEXT等）
	Columns []string // 包含的列
}

// NewConnector 根据配置创建对应的数据库连接器
func NewConnector(config *ConnectionConfig) (Connector, error) {
	switch config.Type {
	case "MySQL":
		return NewMySQLConnector(config), nil
	case "PostgreSQL":
		return NewPostgreSQLConnector(config), nil
	case "SQLite":
		return NewSQLiteConnector(config), nil
	default:
		return nil, nil
	}
}
