// 🎯 完整示例脚本 - 生成带完整注释和验证的 TypeScript 接口
console.log("开始生成 TypeScript 代码...");

/**
 * 输入对象结构：
 * input: {
 *   tableName: string,           // 表名
 *   fields: Array<{             // 字段数组
 *     name: string,            // 字段名
 *     type: string,            // 数据库类型
 *     tsType: string,          // TypeScript 类型
 *     isPrimary: boolean,      // 是否主键
 *     isNullable: boolean,     // 是否可为空
 *     comment: string,         // 字段注释
 *     defaultValue: string     // 默认值
 *   }>,
 *   comment: string,            // 表注释
 *   indexes: Array<{           // 索引数组
 *     name: string,            // 索引名
 *     columns: string[],       // 索引包含的列
 *     isUnique: boolean        // 是否唯一索引
 *   }>,
 *   foreignKeys: Array<{       // 外键数组
 *     name: string,            // 外键名
 *     columns: string[],       // 外键列
 *     refTable: string,        // 引用表
 *     refColumns: string[]     // 引用列
 *   }>
 * }
 */

// 生成完整的 TypeScript 接口
output = generateCompleteInterface(input);

function generateCompleteInterface(data) {
    let code = '';
    
    // 1. 添加文件头注释
    code += generateFileHeader(data);
    
    // 2. 生成导入语句
    code += generateImports(data);
    
    // 3. 生成主接口
    code += generateMainInterface(data);
    
    // 4. 生成创建DTO
    code += generateCreateDTO(data);
    
    // 5. 生成更新DTO
    code += generateUpdateDTO(data);
    
    // 6. 生成查询参数接口
    code += generateQueryParams(data);
    
    console.log("代码生成完成！");
    return code;
}

// 生成文件头注释
function generateFileHeader(data) {
    let header = '';
    header += `/**\n`;
    header += ` * ${data.tableName} 表对应的 TypeScript 定义\n`;
    if (data.comment) {
        header += ` * ${data.comment}\n`;
    }
    header += ` * \n`;
    header += ` * 自动生成时间: ${new Date().toLocaleString('zh-CN')}\n`;
    header += ` * 表名: ${data.tableName}\n`;
    header += ` * 字段数量: ${data.fields.length}\n`;
    header += ` * 主键字段: ${data.fields.filter(f => f.isPrimary).map(f => f.name).join(', ') || '无'}\n`;
    header += ` * \n`;
    header += ` * ⚠️ 注意: 此文件为自动生成，请勿手动修改\n`;
    header += ` * 如需自定义生成逻辑，请使用脚本管理功能\n`;
    header += ` */\n\n`;
    return header;
}

// 生成导入语句
function generateImports(data) {
    let imports = '';
    const types = new Set();
    
    // 检查需要导入的类型
    data.fields.forEach(field => {
        const type = field.tsType.toLowerCase();
        if (type.includes('date') || type.includes('time')) {
            types.add('Date');
        }
    });
    
    if (types.size > 0) {
        imports += `// 类型导入\n`;
        types.forEach(type => {
            imports += `// import { ${type} } from './types';\n`;
        });
        imports += `\n`;
    }
    
    return imports;
}

// 生成主接口
function generateMainInterface(data) {
    let interfaceCode = '';
    interfaceCode += `/**\n`;
    interfaceCode += ` * ${data.tableName} 实体接口\n`;
    if (data.comment) {
        interfaceCode += ` * ${data.comment}\n`;
    }
    interfaceCode += ` */\n`;
    interfaceCode += `export interface ${data.tableName} {\n`;
    
    data.fields.forEach(field => {
        // 字段注释
        interfaceCode += `  /**\n`;
        if (field.comment) {
            interfaceCode += `   * ${field.comment}\n`;
        }
        interfaceCode += `   * 类型: ${field.type} -> ${field.tsType}\n`;
        if (field.isPrimary) {
            interfaceCode += `   * 🔑 主键字段\n`;
        }
        if (field.isNullable) {
            interfaceCode += `   * 📌 允许为空\n`;
        }
        if (field.defaultValue) {
            interfaceCode += `   * ⚡ 默认值: ${field.defaultValue}\n`;
        }
        interfaceCode += `   */\n`;
        
        // 字段定义
        interfaceCode += `  ${field.name}${field.isNullable ? '?' : ''}: ${field.tsType};\n\n`;
    });
    
    interfaceCode += `}\n\n`;
    return interfaceCode;
}

// 生成创建DTO
function generateCreateDTO(data) {
    let dtoCode = '';
    dtoCode += `/**\n`;
    dtoCode += ` * 创建 ${data.tableName} 的 DTO\n`;
    dtoCode += ` * 用于创建新记录时的数据验证\n`;
    dtoCode += ` */\n`;
    dtoCode += `export interface Create${data.tableName}Dto {\n`;
    
    data.fields.forEach(field => {
        if (!field.isPrimary) { // 主键通常由数据库自动生成
            dtoCode += `  /** ${field.comment || field.name} */\n`;
            dtoCode += `  ${field.name}${field.isNullable ? '?' : ''}: ${field.tsType};\n\n`;
        }
    });
    
    dtoCode += `}\n\n`;
    return dtoCode;
}

// 生成更新DTO
function generateUpdateDTO(data) {
    let dtoCode = '';
    dtoCode += `/**\n`;
    dtoCode += ` * 更新 ${data.tableName} 的 DTO\n`;
    dtoCode += ` * 用于更新记录时的数据验证\n`;
    dtoCode += ` */\n`;
    dtoCode += `export interface Update${data.tableName}Dto {\n`;
    
    data.fields.forEach(field => {
        if (!field.isPrimary) { // 主键通常不更新
            dtoCode += `  /** ${field.comment || field.name} */\n`;
            dtoCode += `  ${field.name}?: ${field.tsType}; // 可选字段\n\n`;
        }
    });
    
    dtoCode += `}\n\n`;
    return dtoCode;
}

// 生成查询参数接口
function generateQueryParams(data) {
    let queryCode = '';
    queryCode += `/**\n`;
    queryCode += ` * 查询 ${data.tableName} 的参数接口\n`;
    queryCode += ` * 用于分页查询和条件过滤\n`;
    queryCode += ` */\n`;
    queryCode += `export interface ${data.tableName}QueryParams {\n`;
    queryCode += `  /** 页码 */\n`;
    queryCode += `  page?: number;\n\n`;
    queryCode += `  /** 每页数量 */\n`;
    queryCode += `  pageSize?: number;\n\n`;
    queryCode += `  /** 排序字段 */\n`;
    queryCode += `  sortBy?: string;\n\n`;
    queryCode += `  /** 排序方向 */\n`;
    queryCode += `  sortOrder?: 'asc' | 'desc';\n\n`;
    
    // 为每个字段添加查询条件
    data.fields.forEach(field => {
        queryCode += `  /** 根据 ${field.name} 查询 */\n`;
        queryCode += `  ${field.name}?: ${field.tsType};\n\n`;
        
        if (field.tsType === 'string') {
            queryCode += `  /** ${field.name} 模糊查询 */\n`;
            queryCode += `  ${field.name}Like?: string;\n\n`;
        }
        
        if (field.tsType === 'number' || field.tsType === 'Date') {
            queryCode += `  /** ${field.name} 最小值 */\n`;
            queryCode += `  ${field.name}Min?: ${field.tsType};\n\n`;
            queryCode += `  /** ${field.name} 最大值 */\n`;
            queryCode += `  ${field.name}Max?: ${field.tsType};\n\n`;
        }
    });
    
    queryCode += `}\n`;
    return queryCode;
}

// 类型映射函数
function mapType(dbType) {
    const typeMap = {
        // 整数类型
        'int': 'number',
        'integer': 'number',
        'bigint': 'number',
        'smallint': 'number',
        'tinyint': 'number',
        'mediumint': 'number',
        
        // 浮点类型
        'float': 'number',
        'double': 'number',
        'decimal': 'number',
        'numeric': 'number',
        'real': 'number',
        
        // 字符串类型
        'varchar': 'string',
        'char': 'string',
        'text': 'string',
        'longtext': 'string',
        'mediumtext': 'string',
        'tinytext': 'string',
        'nvarchar': 'string',
        'nchar': 'string',
        
        // 日期时间类型
        'datetime': 'Date',
        'timestamp': 'Date',
        'date': 'Date',
        'time': 'Date',
        'year': 'number',
        
        // 布尔类型
        'boolean': 'boolean',
        'bool': 'boolean',
        'bit': 'boolean',
        
        // 其他类型
        'json': 'any',
        'jsonb': 'any',
        'blob': 'any',
        'longblob': 'any',
        'mediumblob': 'any',
        'tinyblob': 'any',
        'binary': 'any',
        'varbinary': 'any',
        'geometry': 'any',
        'point': 'any',
        'linestring': 'any',
        'polygon': 'any',
        'multipoint': 'any',
        'multilinestring': 'any',
        'multipolygon': 'any',
        'geometrycollection': 'any'
    };
    
    return typeMap[dbType.toLowerCase()] || 'any';
}