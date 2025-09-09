package connector

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// PostgreSQLConnector 实现PostgreSQL数据库连接器
type PostgreSQLConnector struct {
	config *ConnectionConfig
	db     *sql.DB
}

// NewPostgreSQLConnector 创建一个新的PostgreSQL连接器
func NewPostgreSQLConnector(config *ConnectionConfig) *PostgreSQLConnector {
	return &PostgreSQLConnector{
		config: config,
	}
}

// Connect 连接到PostgreSQL数据库
func (c *PostgreSQLConnector) Connect() (*sql.DB, error) {
	// 构建连接字符串
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.config.Host,
		c.config.Port,
		c.config.Username,
		c.config.Password,
		c.config.Database)

	// 如果未指定数据库，连接到默认的postgres数据库
	if c.config.Database == "" {
		connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
			c.config.Host,
			c.config.Port,
			c.config.Username,
			c.config.Password)
	}

	// 连接数据库
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("连接PostgreSQL失败: %v", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("PostgreSQL连接测试失败: %v", err)
	}

	c.db = db
	return db, nil
}

// GetDatabases 获取所有数据库
func (c *PostgreSQLConnector) GetDatabases() ([]string, error) {
	if c.db == nil {
		return nil, fmt.Errorf("数据库未连接")
	}

	// 查询所有数据库
	query := `
		SELECT datname 
		FROM pg_database 
		WHERE datistemplate = false AND datname != 'postgres'
		ORDER BY datname
	`

	rows, err := c.db.Query(query)
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
		databases = append(databases, dbName)
	}

	return databases, nil
}

// GetTables 获取指定数据库中的所有表
func (c *PostgreSQLConnector) GetTables(database string) ([]string, error) {
	if c.db == nil {
		return nil, fmt.Errorf("数据库未连接")
	}

	// 在PostgreSQL中，需要重新连接到指定的数据库
	if c.config.Database != database {
		// 关闭当前连接
		c.Close()

		// 更新配置
		newConfig := *c.config
		newConfig.Database = database
		c.config = &newConfig

		// 重新连接
		_, err := c.Connect()
		if err != nil {
			return nil, err
		}
	}

	// 查询所有表
	query := `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		ORDER BY table_name
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
func (c *PostgreSQLConnector) GetTableMetadata(database, table string) (*TableMetadata, error) {
	if c.db == nil {
		return nil, fmt.Errorf("数据库未连接")
	}

	// 确保连接到正确的数据库
	if c.config.Database != database {
		// 关闭当前连接
		c.Close()

		// 更新配置
		newConfig := *c.config
		newConfig.Database = database
		c.config = &newConfig

		// 重新连接
		_, err := c.Connect()
		if err != nil {
			return nil, err
		}
	}

	// 创建表元数据
	metadata := &TableMetadata{
		Name:   table,
		Fields: make([]FieldInfo, 0),
	}

	// 获取表字段信息
	query := `
		SELECT 
			c.column_name, 
			c.data_type, 
			COALESCE(c.character_maximum_length, 0) as length,
			c.is_nullable, 
			c.column_default,
			pgd.description as column_comment,
			(
				SELECT 
					COUNT(*) 
				FROM 
					information_schema.table_constraints tc
					JOIN information_schema.constraint_column_usage ccu 
					ON tc.constraint_name = ccu.constraint_name
				WHERE 
					tc.constraint_type = 'PRIMARY KEY' 
					AND tc.table_name = c.table_name 
					AND ccu.column_name = c.column_name
			) > 0 as is_primary,
			(
				SELECT 
					COUNT(*) 
				FROM 
					information_schema.table_constraints tc
					JOIN information_schema.constraint_column_usage ccu 
					ON tc.constraint_name = ccu.constraint_name
				WHERE 
					tc.constraint_type = 'UNIQUE' 
					AND tc.table_name = c.table_name 
					AND ccu.column_name = c.column_name
			) > 0 as is_unique
		FROM 
			information_schema.columns c
			LEFT JOIN pg_catalog.pg_statio_all_tables st ON (c.table_schema = st.schemaname AND c.table_name = st.relname)
			LEFT JOIN pg_catalog.pg_description pgd ON (pgd.objoid = st.relid AND pgd.objsubid = c.ordinal_position)
		WHERE 
			c.table_name = $1
			AND c.table_schema = 'public'
		ORDER BY 
			c.ordinal_position
	`

	rows, err := c.db.Query(query, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var field FieldInfo
		var isNullable, columnDefault, columnComment sql.NullString
		var isPrimary, isUnique bool

		if err := rows.Scan(
			&field.Name,
			&field.Type,
			&field.Length,
			&isNullable,
			&columnDefault,
			&columnComment,
			&isPrimary,
			&isUnique,
		); err != nil {
			return nil, err
		}

		// 处理是否可为空
		field.IsNullable = isNullable.String == "YES"

		// 处理主键和唯一键
		field.IsPrimary = isPrimary
		field.IsUnique = isUnique

		// 处理默认值
		if columnDefault.Valid {
			field.Default = columnDefault.String
		}

		// 处理注释
		if columnComment.Valid {
			field.Comment = columnComment.String
		}

		metadata.Fields = append(metadata.Fields, field)
	}

	// 获取索引信息
	indexQuery := `
		SELECT
			i.relname as index_name,
			am.amname as index_type,
			array_agg(a.attname) as column_names
		FROM
			pg_index x
			JOIN pg_class c ON c.oid = x.indrelid
			JOIN pg_class i ON i.oid = x.indexrelid
			JOIN pg_am am ON i.relam = am.oid
			JOIN pg_attribute a ON a.attrelid = c.oid AND a.attnum = ANY(x.indkey)
		WHERE
			c.relkind = 'r' AND
			c.relname = $1
		GROUP BY
			i.relname,
			am.amname
	`

	indexRows, err := c.db.Query(indexQuery, table)
	if err != nil {
		return nil, err
	}
	defer indexRows.Close()

	for indexRows.Next() {
		var index IndexInfo
		var columnNames []string

		if err := indexRows.Scan(&index.Name, &index.Type, &columnNames); err != nil {
			return nil, err
		}

		index.Columns = columnNames
		metadata.Indexes = append(metadata.Indexes, index)
	}

	return metadata, nil
}

// Close 关闭数据库连接
func (c *PostgreSQLConnector) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}
