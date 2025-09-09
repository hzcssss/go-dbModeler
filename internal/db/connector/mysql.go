package connector

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// MySQLConnector 实现MySQL数据库连接器
type MySQLConnector struct {
	config *ConnectionConfig
	db     *sql.DB
}

// NewMySQLConnector 创建一个新的MySQL连接器
func NewMySQLConnector(config *ConnectionConfig) *MySQLConnector {
	return &MySQLConnector{
		config: config,
	}
}

// Connect 连接到MySQL数据库
func (c *MySQLConnector) Connect() (*sql.DB, error) {
	// 构建DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		c.config.Username,
		c.config.Password,
		c.config.Host,
		c.config.Port,
		c.config.Database)

	// 连接数据库
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("连接MySQL失败: %v", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("MySQL连接测试失败: %v", err)
	}

	c.db = db
	return db, nil
}

// GetDatabases 获取所有数据库
func (c *MySQLConnector) GetDatabases() ([]string, error) {
	if c.db == nil {
		return nil, fmt.Errorf("数据库未连接")
	}

	// 查询所有数据库
	rows, err := c.db.Query("SHOW DATABASES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return nil, err
		}
		// 过滤系统数据库
		if dbName != "information_schema" && dbName != "mysql" && dbName != "performance_schema" && dbName != "sys" {
			databases = append(databases, dbName)
		}
	}

	return databases, nil
}

// GetTables 获取指定数据库中的所有表
func (c *MySQLConnector) GetTables(database string) ([]string, error) {
	if c.db == nil {
		return nil, fmt.Errorf("数据库未连接")
	}

	// 切换到指定数据库
	if _, err := c.db.Exec("USE " + database); err != nil {
		return nil, err
	}

	// 查询所有表
	rows, err := c.db.Query("SHOW TABLES")
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
func (c *MySQLConnector) GetTableMetadata(database, table string) (*TableMetadata, error) {
	if c.db == nil {
		return nil, fmt.Errorf("数据库未连接")
	}

	// 切换到指定数据库
	if _, err := c.db.Exec("USE " + database); err != nil {
		return nil, err
	}

	// 创建表元数据
	metadata := &TableMetadata{
		Name:   table,
		Fields: make([]FieldInfo, 0),
	}

	// 获取表字段信息
	query := `
		SELECT 
			COLUMN_NAME, 
			DATA_TYPE, 
			IFNULL(CHARACTER_MAXIMUM_LENGTH, 0) as LENGTH,
			IS_NULLABLE, 
			COLUMN_KEY, 
			COLUMN_DEFAULT, 
			COLUMN_COMMENT
		FROM 
			INFORMATION_SCHEMA.COLUMNS
		WHERE 
			TABLE_SCHEMA = ? AND TABLE_NAME = ?
		ORDER BY 
			ORDINAL_POSITION
	`

	rows, err := c.db.Query(query, database, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var field FieldInfo
		var isNullable, columnKey, columnDefault sql.NullString

		if err := rows.Scan(
			&field.Name,
			&field.Type,
			&field.Length,
			&isNullable,
			&columnKey,
			&columnDefault,
			&field.Comment,
		); err != nil {
			return nil, err
		}

		// 处理是否可为空
		field.IsNullable = isNullable.String == "YES"

		// 处理主键和唯一键
		if columnKey.Valid {
			field.IsPrimary = columnKey.String == "PRI"
			field.IsUnique = columnKey.String == "UNI"
		}

		// 处理默认值
		if columnDefault.Valid {
			field.Default = columnDefault.String
		}

		metadata.Fields = append(metadata.Fields, field)
	}

	// 获取索引信息
	indexQuery := `
		SELECT 
			INDEX_NAME,
			INDEX_TYPE,
			COLUMN_NAME
		FROM 
			INFORMATION_SCHEMA.STATISTICS
		WHERE 
			TABLE_SCHEMA = ? AND TABLE_NAME = ?
		ORDER BY 
			INDEX_NAME, SEQ_IN_INDEX
	`

	indexRows, err := c.db.Query(indexQuery, database, table)
	if err != nil {
		return nil, err
	}
	defer indexRows.Close()

	indexMap := make(map[string]*IndexInfo)
	for indexRows.Next() {
		var indexName, indexType, columnName string

		if err := indexRows.Scan(&indexName, &indexType, &columnName); err != nil {
			return nil, err
		}

		// 如果索引不存在，创建它
		if _, exists := indexMap[indexName]; !exists {
			indexMap[indexName] = &IndexInfo{
				Name:    indexName,
				Type:    indexType,
				Columns: make([]string, 0),
			}
		}

		// 添加列到索引
		indexMap[indexName].Columns = append(indexMap[indexName].Columns, columnName)
	}

	// 将索引映射转换为切片
	for _, index := range indexMap {
		metadata.Indexes = append(metadata.Indexes, *index)
	}

	return metadata, nil
}

// Close 关闭数据库连接
func (c *MySQLConnector) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}
