# go-DBmodeler

go-DBmodeler 是一个使用 Go 语言开发的轻量级数据库建模工具。

## 功能概述

GoDBModeler 是一个简化版的数据库建模工具，支持从多种数据库生成 TypeScript 模型代码，包含以下核心功能：

- **连接管理**：支持 MySQL、PostgreSQL、SQLite 数据库连接配置
- **TS模型生成**：使用默认模板和自定义脚本生成 TypeScript 代码

## 技术栈

- Go 1.21+
- Fyne UI 框架
- MySQL/PostgreSQL/SQLite 数据库驱动

## 快速开始

### 克隆项目

```bash
git clone <your-repo-url>
cd go-DBmodeler
```

### 安装依赖

```bash
go mod download
```

### 运行应用

```bash
go run cmd/app/main.go
```

### 构建应用

```bash
go build -o go-DBmodeler cmd/app/main.go
```

## 打包指南

### 打包前清理

为确保打包的应用程序不包含任何敏感数据（如数据库连接信息），请先运行清理脚本：

```bash
# 运行打包前清理脚本
./scripts/package_clean.sh
```

### macOS 打包

```bash
# 编译 macOS 版本
GOOS=darwin GOARCH=arm64 go build -o bin/go-DBmodeler ./cmd/app/

# 使用 fyne 工具打包为 .app（确保指定正确的可执行文件）
fyne package --os darwin --executable bin/go-DBmodeler --app-id com.godbmodeler.app --name DBmodeler

# 创建 DMG 安装包
hdiutil create -volname "DBmodeler" -srcfolder DBmodeler.app -ov -format UDZO DBmodeler.dmg
```

## 项目结构

```
go-DBmodeler/
├── cmd/app/           # 主应用程序入口
├── internal/          # 内部包
│   ├── app/          # 应用核心
│   ├── config/       # 配置管理
│   ├── db/           # 数据库连接和元数据
│   ├── generator/    # 代码生成器
│   └── ui/           # 用户界面
├── resources/         # 资源文件
├── scripts/           # 脚本文件
└── templates/         # 代码模板
```

## 开发指南

### 环境要求

- Go 1.21 或更高版本
- Fyne 依赖项

### 依赖管理

使用 Go Modules 管理依赖：

```bash
# 添加新依赖
go get <package>

# 清理未使用的依赖
go mod tidy
```

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情