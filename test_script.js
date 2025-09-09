// 测试脚本 - 生成带注释的 TypeScript 接口
console.log("开始生成 TypeScript 代码...");

// 输入对象包含表结构信息
// input: { tableName, columns, primaryKey, indexes, foreignKeys }

// 生成带注释的接口
output = generateInterface(input);

function generateInterface(data) {
    let code = '';
    
    // 添加文件头注释
    code += `/**\n * ${data.tableName} 表对应的 TypeScript 接口\n`;
    code += ` * 生成时间: ${new Date().toLocaleString()}\n */\n\n`;
    
    // 生成接口定义
    code += `export interface ${data.tableName} {\n`;
    
    data.columns.forEach(column => {
        // 添加字段注释
        code += `  /**\n   * ${column.name} - ${column.type}\n`;
        if (column.isPrimary) {
            code += `   * 主键字段\n`;
        }
        if (column.isNullable) {
            code += `   * 允许为空\n`;
        }
        code += `   */\n`;
        
        code += `  ${column.name}${column.isNullable ? '?' : ''}: ${mapType(column.type)};\n\n`;
    });
    
    code += '}\n';
    
    // 添加工具函数
    code += '\n// 类型映射工具函数\n';
    code += 'function mapType(dbType: string): string {\n';
    code += '  const typeMap: Record<string, string> = {\n';
    code += "    'int': 'number',\n";
    code += "    'varchar': 'string',\n"; 
    code += "    'text': 'string',\n";
    code += "    'datetime': 'Date',\n";
    code += "    'boolean': 'boolean'\n";
    code += '  };\n';
    code += "  return typeMap[dbType.toLowerCase()] || 'any';\n";
    code += '}\n';
    
    console.log("代码生成完成！");
    return code;
}

function mapType(dbType) {
    const typeMap = {
        'int': 'number',
        'integer': 'number',
        'bigint': 'number',
        'smallint': 'number',
        'tinyint': 'number',
        'float': 'number',
        'double': 'number',
        'decimal': 'number',
        'numeric': 'number',
        'varchar': 'string',
        'char': 'string',
        'text': 'string',
        'longtext': 'string',
        'mediumtext': 'string',
        'datetime': 'Date',
        'timestamp': 'Date',
        'date': 'Date',
        'time': 'Date',
        'boolean': 'boolean',
        'bool': 'boolean',
        'json': 'any',
        'blob': 'any'
    };
    return typeMap[dbType.toLowerCase()] || 'any';
}