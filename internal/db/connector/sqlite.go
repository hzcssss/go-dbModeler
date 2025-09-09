package connector

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
)

// SQLiteConnector 实现SQLite数据库连接器
type SQLiteConnector struct {
	config *ConnectionConfig
	db     *sql.DB
}

// NewSQLiteConnector 创建一个新的SQLite连接器
func NewSQLiteConnector(config *ConnectionConfig) *SQLiteConnector {
	return &SQLiteConnector{
		config: config,
	}
}

// Connect 连接到SQLite数据库
func (c *SQLiteConnector) Connect() (*sql.DB, error) {
	// 对于SQLite，Host字段实际上是文件路径
	dbPath := c.config.Host

	// 检查文件是否存在
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("SQLite数据库文件不存在: %s", dbPath)
	}

	// 连接数据库
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("连接SQLite失败: %v", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("SQLite连接测试失败: %v", err)
	}

	c.db = db
	return db, nil
}

// GetDatabases 获取所有数据库
// 注意：SQLite不支持多数据库，返回文件名作为数据库名
func (c *SQLiteConnector) GetDatabases() ([]string, error) {
	if c.db == nil {
		return nil, fmt.Errorf("数据库未连接")
	}

	// 返回文件名作为数据库名
	dbName := filepath.Base(c.config.Host)
	return []string{dbName}, nil
}

// GetTables 获取指定数据库中的所有表
func (c *SQLiteConnector) GetTables(database string) ([]string, error) {
	if c.db == nil {
		return nil, fmt.Errorf("数据库未连接")
	}

	// 查询所有表
	query := `
		SELECT name 
		FROM sqlite_master 
		WHERE type='table' AND name NOT LIKE 'sqlite_%'
		ORDER BY name
	`

	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}

	return tables, nil
}

// GetTableMetadata 获取表的元数据信息
func (c *SQLiteConnector) GetTableMetadata(database, table string) (*TableMetadata, error) {
	if c.db == nil {
		return nil, fmt.Errorf("数据库未连接")
	}

	// 创建表元数据
	metadata := &TableMetadata{
		Name:   table,
		Fields: make([]FieldInfo, 0),
	}

	// 获取表结构信息
	query := fmt.Sprintf("PRAGMA table_info(%s)", table)
	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var field FieldInfo
		var notNull, pk int

		if err := rows.Scan(&cid, &field.Name, &field.Type, &notNull, &field.Default, &pk); err != nil {
			return nil, err
		}

		// 处理是否可为空
		field.IsNullable = notNull == 0

		// 处理主键
		field.IsPrimary = pk > 0

		// SQLite不直接支持字段注释，可以通过其他方式获取

		metadata.Fields = append(metadata.Fields, field)
	}

	// 获取索引信息
	indexListQuery := fmt.Sprintf("PRAGMA index_list(%s)", table)
	indexRows, err := c.db.Query(indexListQuery)
	if err != nil {
		return nil, err
	}
	defer indexRows.Close()

	for indexRows.Next() {
		var seq int
		var indexName string
		var unique int
		var origin, partial string

		if err := indexRows.Scan(&seq, &indexName, &unique, &origin, &partial); err != nil {
			return nil, err
		}

		// 获取索引列
		indexInfoQuery := fmt.Sprintf("PRAGMA index_info(%s)", indexName)
		infoRows, err := c.db.Query(indexInfoQuery)
		if err != nil {
			return nil, err
		}

		var columns []string
		for infoRows.Next() {
			var seqno, cid int
			var columnName string

			if err := infoRows.Scan(&seqno, &cid, &columnName); err != nil {
				infoRows.Close()
				return nil, err
			}

			columns = append(columns, columnName)
		}
		infoRows.Close()

		// 确定索引类型
		indexType := "INDEX"
		if unique == 1 {
			indexType = "UNIQUE"
		}

		index := IndexInfo{
			Name:    indexName,
			Type:    indexType,
			Columns: columns,
		}

		metadata.Indexes = append(metadata.Indexes, index)
	}

	return metadata, nil
}

// Close 关闭数据库连接
func (c *SQLiteConnector) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}
