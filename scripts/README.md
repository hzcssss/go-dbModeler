# JavaScript 脚本目录

此目录用于存放导入的 JavaScript 脚本文件。

## 目录结构

```
scripts/
├── imported/          # 存放导入的 JavaScript 脚本文件
│   ├── camelcase_script.js     # 小驼峰命名转换脚本示例
│   ├── example_script.js       # 示例脚本
└── README.md          # 说明文档
```

## 脚本使用说明

### 脚本格式要求
1. **输入变量**: `input` - 已解析的 JSON 对象，包含表结构数据
2. **代码变量**: `tsCode` - 当前生成的 TypeScript 代码
3. **输出变量**: `output` - 处理后的结果

### 输入数据结构
```javascript
{
  "tableName": "表名",
  "fields": [
    {
      "name": "字段名",
      "tsType": "TypeScript类型",
      "comment": "字段注释"
    }
  ]
}
```

### 示例脚本
```javascript
// 小驼峰命名转换示例
function toCamelCase(str) {
    return str.replace(/_([a-z])/g, (g) => g[1].toUpperCase())
              .replace(/^[A-Z]/, (g) => g.toLowerCase());
}

// 处理字段名
let result = "export interface " + input.tableName + " {\n";
for (const field of input.fields) {
    const camelCaseName = toCamelCase(field.name);
    result += "  " + camelCaseName + ": " + field.tsType + ";\n";
}
result += "}\n";

// 设置输出
output = result;
```

## 导入方法

1. 将脚本文件放入 `scripts/imported/` 目录
2. 在程序中使用文件导入功能加载脚本
3. 或者在程序界面中复制粘贴脚本内容

## 注意事项

- 脚本文件使用 UTF-8 编码
- 确保脚本语法正确，可以在 JavaScript 环境中运行
- 脚本应处理可能的错误情况