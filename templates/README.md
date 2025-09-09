# 模板管理操作指南

## 目录结构
```
templates/
├── imported/          # 用户自定义模板文件目录
│   ├── default.tpl    # 默认模板（系统自动创建）
│   └── custom.tpl     # 用户自定义模板
└── README.md         # 操作说明文档
```

## 模板文件格式
模板文件使用 `.tpl` 扩展名，内容为 Go 模板语法。

## 可用参数和变量

### 根对象参数
模板接收以下根对象参数：

```go
type TemplateData struct {
    PackageName string    // 包名
    StructName  string    // 结构体名称
    TableName   string    // 数据库表名
    Fields      []Field   // 字段列表
    Imports     []string  // 需要导入的包
}
```

### Field 结构体
```go
type Field struct {
    Name     string // 字段名（Go风格，首字母大写）
    Type     string // 字段类型（Go类型）
    JsonName string // JSON字段名（小写蛇形命名）
    DbName   string // 数据库字段名（小写蛇形命名）
    Comment  string // 字段注释
    IsPrimary bool  // 是否为主键
    IsNullable bool // 是否可为空
    DefaultValue string // 默认值
    MaxLength    int    // 最大长度（字符串类型）
}
```

### 模板函数
除了标准Go模板函数外，还提供以下自定义函数：

- `camelCase str` - 转换为驼峰命名
- `pascalCase str` - 转换为帕斯卡命名
- `snakeCase str` - 转换为蛇形命名
- `lower str` - 转换为小写
- `upper str` - 转换为大写
- `title str` - 首字母大写
- `plural str` - 转换为复数形式
- `singular str` - 转换为单数形式

### 示例模板
```go
// 基础结构体模板
package {{.PackageName}}

import (
    {{range .Imports}}
    "{{.}}"{{end}}
    "time"
    "gorm.io/gorm"
)

// {{.StructName}} 对应数据库表 {{.TableName}}
type {{.StructName}} struct {
    {{range .Fields}}
    {{.Name}} {{.Type}} `json:"{{.JsonName}}" gorm:"column:{{.DbName}}{{if .IsPrimary}};primaryKey{{end}}{{if .IsNullable}};null{{end}}"`{{end}}
    
    // 公共字段
    CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func ({{.StructName}}) TableName() string {
    return "{{.TableName}}"
}

// 带方法的模板示例
package {{.PackageName}}

import (
    {{range .Imports}}
    "{{.}}"{{end}}
    "context"
    "gorm.io/gorm"
)

type {{.StructName}} struct {
    {{range .Fields}}
    {{.Name}} {{.Type}} `json:"{{.JsonName}}" gorm:"column:{{.DbName}}"`{{end}}
}

// Create 创建记录
func (m *{{.StructName}}) Create(ctx context.Context, db *gorm.DB) error {
    return db.WithContext(ctx).Create(m).Error
}

// Update 更新记录
func (m *{{.StructName}}) Update(ctx context.Context, db *gorm.DB) error {
    return db.WithContext(ctx).Save(m).Error
}

// Delete 删除记录
func (m *{{.StructName}}) Delete(ctx context.Context, db *gorm.DB) error {
    return db.WithContext(ctx).Delete(m).Error
}

// FindByID 根据ID查找
func Find{{.StructName}}ByID(ctx context.Context, db *gorm.DB, id uint) (*{{.StructName}}, error) {
    var result {{.StructName}}
    err := db.WithContext(ctx).First(&result, id).Error
    return &result, err
}
```

### 条件判断示例
```go
package {{.PackageName}}

type {{.StructName}} struct {
    {{range .Fields}}
    {{.Name}} {{.Type}} `json:"{{.JsonName}}"`{{if .Comment}} // {{.Comment}}{{end}}{{end}}
}

// 根据字段类型添加验证方法
{{$hasString := false}}
{{range .Fields}}
{{if eq .Type "string"}}{{$hasString = true}}{{end}}
{{end}}

{{if $hasString}}
// Validate 验证字段
func (m *{{.StructName}}) Validate() error {
    {{range .Fields}}
    {{if eq .Type "string"}}
    if len(m.{{.Name}}) > {{.MaxLength}} {
        return fmt.Errorf("{{.Name}} 长度不能超过 %d", {{.MaxLength}})
    }
    {{end}}
    {{end}}
    return nil
}
{{end}}
```

### 使用自定义函数
```go
package {{.PackageName}}

// {{.StructName | camelCase}} 对应表 {{.TableName | snakeCase}}
type {{.StructName}} struct {
    {{range .Fields}}
    // {{.Comment}}
    {{.Name}} {{.Type}} `json:"{{.JsonName | snakeCase}}" db:"{{.DbName | snakeCase}}"`{{end}}
}

// 生成DTO对象
type {{.StructName}}DTO struct {
    {{range .Fields}}
    {{.Name}} {{.Type}} `json:"{{.JsonName}}"`{{end}}
}

// ToDTO 转换为DTO
func (m *{{.StructName}}) ToDTO() {{.StructName}}DTO {
    return {{.StructName}}DTO{
        {{range .Fields}}
        {{.Name}}: m.{{.Name}},{{end}}
    }
}
```

## 高级模板技巧

### 1. 条件导入包
```go
package {{.PackageName}}

import (
    "fmt"
    "time"{{if .HasTimeFields}}
    "time"{{end}}{{if .HasJSONFields}}
    "encoding/json"{{end}}{{if .HasDatabaseFields}}
    "gorm.io/gorm"{{end}}
)
```

### 2. 字段分组处理
```go
type {{.StructName}} struct {
    // 基础字段
    {{range .Fields}}
    {{if not .IsPrimary}}
    {{.Name}} {{.Type}} `json:"{{.JsonName}}"`{{end}}
    {{end}}
    
    // 主键字段
    {{range .Fields}}
    {{if .IsPrimary}}
    {{.Name}} {{.Type}} `json:"{{.JsonName}}" gorm:"primaryKey"`{{end}}
    {{end}}
}
```

### 3. 自定义方法生成
```go
{{$structName := .StructName}}
// 为每个字段生成Getter方法
{{range .Fields}}
func (m *{{$structName}}) Get{{.Name}}() {{.Type}} {
    return m.{{.Name}}
}

func (m *{{$structName}}) Set{{.Name}}(value {{.Type}}) {
    m.{{.Name}} = value
}
{{end}}
```

### 4. 枚举字段处理
```go
type {{.StructName}} struct {
    {{range .Fields}}
    {{.Name}} {{.Type}} `json:"{{.JsonName}}"`{{end}}
}

// 状态常量
const (
    {{.StructName}}StatusActive = "active"
    {{.StructName}}StatusInactive = "inactive"
)

// 验证状态
func (m *{{.StructName}}) IsValidStatus() bool {
    switch m.Status {
    case {{.StructName}}StatusActive, {{.StructName}}StatusInactive:
        return true
    default:
        return false
    }
}
```

## 模板变量参考表

| 变量名 | 类型 | 描述 | 示例 |
|--------|------|------|------|
| `.PackageName` | string | Go包名 | `model` |
| `.StructName` | string | 结构体名称 | `User` |
| `.TableName` | string | 数据库表名 | `users` |
| `.Fields` | []Field | 字段列表 | - |
| `.Fields[].Name` | string | 字段名（Go风格） | `UserName` |
| `.Fields[].Type` | string | 字段类型 | `string` |
| `.Fields[].JsonName` | string | JSON字段名 | `user_name` |
| `.Fields[].DbName` | string | 数据库字段名 | `user_name` |
| `.Fields[].Comment` | string | 字段注释 | `用户姓名` |
| `.Fields[].IsPrimary` | bool | 是否为主键 | `true` |
| `.Fields[].IsNullable` | bool | 是否可为空 | `false` |
| `.Fields[].DefaultValue` | string | 默认值 | `""` |
| `.Fields[].MaxLength` | int | 最大长度 | `255` |

## 最佳实践

1. **模板设计原则**
   - 保持模板简洁易读
   - 使用有意义的模板变量名
   - 添加适当的注释说明

2. **错误处理**
   - 在模板中使用条件判断处理可能为空的值
   - 为关键操作添加验证方法

3. **性能优化**
   - 避免在模板中进行复杂的计算
   - 使用预定义的函数进行处理

4. **可维护性**
   - 将复杂模板拆分为多个简单模板
   - 使用一致的命名约定
   - 定期备份重要模板

## 常见问题

### Q: 模板变量不显示怎么办？
A: 检查变量名拼写是否正确，确保使用 `.VariableName` 格式

### Q: 如何添加自定义函数？
A: 在模板中使用 `{{/* 注释 */}}` 添加注释，或联系开发人员添加新的模板函数

### Q: 模板语法错误如何调试？
A: 检查花括号 `{{}}` 是否匹配，确保所有模板指令正确闭合

### Q: 如何创建嵌套结构体？
A: 使用 range 循环和条件判断来生成嵌套结构
```go
type {{.StructName}} struct {
    {{range .Fields}}
    {{if .IsComplexType}}
    {{.Name}} struct {
        // 嵌套字段
    } `json:"{{.JsonName}}"`
    {{else}}
    {{.Name}} {{.Type}} `json:"{{.JsonName}}"`
    {{end}}
    {{end}}
}
```

## 操作说明

### 1. 创建新模板
1. 在应用程序中点击"新建模板"按钮
2. 输入模板名称（如：`user`）
3. 系统会自动创建 `templates/imported/user.tpl` 文件
4. 编辑模板内容并保存

### 2. 编辑现有模板
1. 在模板列表中选择要编辑的模板
2. 在编辑器中修改模板内容
3. 点击"保存模板"按钮
4. 系统会同时更新存储和文件系统中的模板

### 3. 删除模板
1. 在模板列表中选择要删除的模板（默认模板不可删除）
2. 点击"删除模板"按钮
3. 确认删除操作
4. 系统会同时从存储和文件系统中删除模板

### 4. 手动管理模板文件
您也可以直接操作文件系统中的模板文件：

#### 添加新模板
```bash
# 创建新的模板文件
echo 'package {{.PackageName}}

type {{.StructName}} struct {
    {{range .Fields}}
    {{.Name}} {{.Type}} `json:"{{.JsonName}}"`{{end}}
}' > templates/imported/my_template.tpl
```

#### 修改现有模板
```bash
# 编辑模板文件
vim templates/imported/default.tpl
```

#### 删除模板文件
```bash
# 删除模板文件（应用程序重启后会从存储中同步删除）
rm templates/imported/unwanted.tpl
```

## 默认内置模板
系统提供以下默认模板，这些模板在首次运行时自动创建：

- **default.tpl** - 基础结构体模板
- **with_methods.tpl** - 包含方法的模板
- **gorm_model.tpl** - GORM 模型模板
- **json_api.tpl** - JSON API 响应模板

## 注意事项

1. **文件系统同步**：应用程序会自动监控 `templates/imported/` 目录的变化
2. **默认模板保护**：名为 "default" 的模板受到保护，无法删除
3. **文件编码**：模板文件使用 UTF-8 编码
4. **重启生效**：手动修改文件后需要重启应用程序才能生效
5. **备份建议**：定期备份重要的模板文件

## 故障排除

### 模板文件无法保存
- 检查 `templates/imported/` 目录的写入权限
- 确保磁盘空间充足

### 模板不显示
- 检查文件扩展名是否为 `.tpl`
- 重启应用程序重新加载模板

### 模板语法错误
- 使用有效的 Go 模板语法
- 确保所有模板变量都有正确的定义