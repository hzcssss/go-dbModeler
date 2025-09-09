#!/bin/bash

# 打包前清理脚本
# 确保不包含敏感数据

echo "开始清理打包环境..."

# 确保不存在临时配置文件
if [ -f "./config.json" ]; then
    echo "删除临时配置文件..."
    rm ./config.json
fi

# 确保不存在测试数据库连接
if [ -d "./testdata" ]; then
    echo "删除测试数据..."
    rm -rf ./testdata
fi

# 确保不存在日志文件
if [ -f "./godbmodeler.log" ]; then
    echo "删除日志文件..."
    rm ./godbmodeler.log
fi

echo "清理完成，可以安全打包"

# 打包命令提示
echo ""
echo "使用以下命令打包应用程序："
echo "1. 编译应用："
echo "   go build -o bin/go-DBmodeler ./cmd/app/"
echo ""
echo "2. 打包为 .app："
echo "   fyne package --os darwin --executable bin/go-DBmodeler --app-id com.godbmodeler.app --name DBmodeler"
echo ""
echo "3. 创建 DMG："
echo "   hdiutil create -volname \"DBmodeler\" -srcfolder DBmodeler.app -ov -format UDZO DBmodeler.dmg"