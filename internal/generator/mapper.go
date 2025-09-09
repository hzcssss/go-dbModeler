package generator

// TypeMapper 定义类型映射器接口
type TypeMapper interface {
	// Map 将数据库类型映射为TypeScript类型
	Map(dbType string) string
}

// MySQLMapper 实现MySQL类型到TypeScript类型的映射
type MySQLMapper struct {
	// 类型映射表
	typeMap map[string]string
}

// NewMySQLMapper 创建一个新的MySQL类型映射器
func NewMySQLMapper() *MySQLMapper {
	mapper := &MySQLMapper{
		typeMap: make(map[string]string),
	}

	// 初始化默认映射
	mapper.typeMap["int"] = "number"
	mapper.typeMap["tinyint"] = "number"
	mapper.typeMap["smallint"] = "number"
	mapper.typeMap["mediumint"] = "number"
	mapper.typeMap["bigint"] = "number"
	mapper.typeMap["float"] = "number"
	mapper.typeMap["double"] = "number"
	mapper.typeMap["decimal"] = "number"

	mapper.typeMap["char"] = "string"
	mapper.typeMap["varchar"] = "string"
	mapper.typeMap["tinytext"] = "string"
	mapper.typeMap["text"] = "string"
	mapper.typeMap["mediumtext"] = "string"
	mapper.typeMap["longtext"] = "string"

	mapper.typeMap["date"] = "Date"
	mapper.typeMap["datetime"] = "Date"
	mapper.typeMap["timestamp"] = "Date"
	mapper.typeMap["time"] = "string"
	mapper.typeMap["year"] = "number"

	mapper.typeMap["tinyint(1)"] = "boolean"
	mapper.typeMap["bit"] = "boolean"

	mapper.typeMap["json"] = "any"
	mapper.typeMap["enum"] = "string"
	mapper.typeMap["set"] = "string[]"

	mapper.typeMap["binary"] = "Buffer"
	mapper.typeMap["varbinary"] = "Buffer"
	mapper.typeMap["blob"] = "Buffer"

	return mapper
}

// Map 将MySQL类型映射为TypeScript类型
func (m *MySQLMapper) Map(dbType string) string {
	// 检查是否为tinyint(1)，这通常表示布尔值
	if dbType == "tinyint(1)" {
		return "boolean"
	}

	// 提取基本类型（去掉长度等信息）
	baseType := dbType
	for i, c := range dbType {
		if c == '(' || c == ' ' {
			baseType = dbType[:i]
			break
		}
	}

	// 查找映射
	if tsType, ok := m.typeMap[baseType]; ok {
		return tsType
	}

	// 默认为any类型
	return "any"
}

// PostgreSQLMapper 实现PostgreSQL类型到TypeScript类型的映射
type PostgreSQLMapper struct {
	// 类型映射表
	typeMap map[string]string
}

// NewPostgreSQLMapper 创建一个新的PostgreSQL类型映射器
func NewPostgreSQLMapper() *PostgreSQLMapper {
	mapper := &PostgreSQLMapper{
		typeMap: make(map[string]string),
	}

	// 初始化默认映射
	mapper.typeMap["smallint"] = "number"
	mapper.typeMap["integer"] = "number"
	mapper.typeMap["bigint"] = "number"
	mapper.typeMap["decimal"] = "number"
	mapper.typeMap["numeric"] = "number"
	mapper.typeMap["real"] = "number"
	mapper.typeMap["double precision"] = "number"
	mapper.typeMap["serial"] = "number"
	mapper.typeMap["bigserial"] = "number"

	mapper.typeMap["varchar"] = "string"
	mapper.typeMap["character varying"] = "string"
	mapper.typeMap["character"] = "string"
	mapper.typeMap["text"] = "string"

	mapper.typeMap["timestamp"] = "Date"
	mapper.typeMap["timestamp with time zone"] = "Date"
	mapper.typeMap["timestamp without time zone"] = "Date"
	mapper.typeMap["date"] = "Date"
	mapper.typeMap["time"] = "string"
	mapper.typeMap["time with time zone"] = "string"
	mapper.typeMap["time without time zone"] = "string"
	mapper.typeMap["interval"] = "string"

	mapper.typeMap["boolean"] = "boolean"

	mapper.typeMap["json"] = "any"
	mapper.typeMap["jsonb"] = "any"
	mapper.typeMap["uuid"] = "string"
	mapper.typeMap["inet"] = "string"
	mapper.typeMap["cidr"] = "string"
	mapper.typeMap["macaddr"] = "string"

	mapper.typeMap["bytea"] = "Buffer"

	return mapper
}

// Map 将PostgreSQL类型映射为TypeScript类型
func (m *PostgreSQLMapper) Map(dbType string) string {
	// 查找映射
	if tsType, ok := m.typeMap[dbType]; ok {
		return tsType
	}

	// 默认为any类型
	return "any"
}

// SQLiteMapper 实现SQLite类型到TypeScript类型的映射
type SQLiteMapper struct {
	// 类型映射表
	typeMap map[string]string
}

// NewSQLiteMapper 创建一个新的SQLite类型映射器
func NewSQLiteMapper() *SQLiteMapper {
	mapper := &SQLiteMapper{
		typeMap: make(map[string]string),
	}

	// 初始化默认映射
	mapper.typeMap["integer"] = "number"
	mapper.typeMap["int"] = "number"
	mapper.typeMap["tinyint"] = "number"
	mapper.typeMap["smallint"] = "number"
	mapper.typeMap["mediumint"] = "number"
	mapper.typeMap["bigint"] = "number"
	mapper.typeMap["real"] = "number"
	mapper.typeMap["double"] = "number"
	mapper.typeMap["float"] = "number"
	mapper.typeMap["numeric"] = "number"

	mapper.typeMap["text"] = "string"
	mapper.typeMap["char"] = "string"
	mapper.typeMap["varchar"] = "string"
	mapper.typeMap["varying character"] = "string"
	mapper.typeMap["nchar"] = "string"
	mapper.typeMap["native character"] = "string"
	mapper.typeMap["nvarchar"] = "string"

	mapper.typeMap["date"] = "Date"
	mapper.typeMap["datetime"] = "Date"
	mapper.typeMap["timestamp"] = "Date"

	mapper.typeMap["boolean"] = "boolean"

	mapper.typeMap["blob"] = "Buffer"

	return mapper
}

// Map 将SQLite类型映射为TypeScript类型
func (m *SQLiteMapper) Map(dbType string) string {
	// 提取基本类型（去掉长度等信息）
	baseType := dbType
	for i, c := range dbType {
		if c == '(' || c == ' ' {
			baseType = dbType[:i]
			break
		}
	}

	// 检查是否为布尔值
	if baseType == "tinyint" && dbType == "tinyint(1)" {
		return "boolean"
	}

	// 查找映射
	if tsType, ok := m.typeMap[baseType]; ok {
		return tsType
	}

	// 默认为any类型
	return "any"
}

// NewTypeMapper 根据数据库类型创建对应的类型映射器
func NewTypeMapper(dbType string) TypeMapper {
	switch dbType {
	case "MySQL":
		return NewMySQLMapper()
	case "PostgreSQL":
		return NewPostgreSQLMapper()
	case "SQLite":
		return NewSQLiteMapper()
	default:
		return NewMySQLMapper() // 默认使用MySQL映射器
	}
}
