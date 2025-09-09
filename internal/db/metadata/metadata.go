package metadata

import (
	"go-DBmodeler/internal/db/connector"
)

// Processor 表示元数据处理器
type Processor struct {
	connector connector.Connector
}

// NewProcessor 创建一个新的元数据处理器
func NewProcessor(connector connector.Connector) *Processor {
	return &Processor{
		connector: connector,
	}
}

// GetDatabases 获取所有数据库
func (p *Processor) GetDatabases() ([]string, error) {
	return p.connector.GetDatabases()
}

// GetTables 获取指定数据库中的所有表
func (p *Processor) GetTables(database string) ([]string, error) {
	return p.connector.GetTables(database)
}

// GetTableMetadata 获取表的元数据信息
func (p *Processor) GetTableMetadata(database, table string) (*connector.TableMetadata, error) {
	return p.connector.GetTableMetadata(database, table)
}

// Close 关闭数据库连接
func (p *Processor) Close() error {
	return p.connector.Close()
}
