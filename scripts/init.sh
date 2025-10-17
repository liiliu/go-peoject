#!/bin/bash

# 项目初始化脚本
# 使用方法: ./scripts/init.sh your_project_name

if [ -z "$1" ]; then
    echo "使用方法: ./scripts/init.sh your_project_name"
    echo "示例: ./scripts/init.sh my_awesome_project"
    exit 1
fi

PROJECT_NAME=$1

echo "🚀 正在初始化项目: $PROJECT_NAME"
echo ""

# 1. 替换模块名称
echo "📝 步骤 1/5: 更新 go.mod 模块名称..."
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    sed -i '' "s/your_project/$PROJECT_NAME/g" go.mod
else
    # Linux
    sed -i "s/your_project/$PROJECT_NAME/g" go.mod
fi

# 2. 替换所有源代码中的 import 路径
echo "📝 步骤 2/5: 更新源代码中的 import 路径..."
find . -type f -name "*.go" -exec sed -i.bak "s/your_project/$PROJECT_NAME/g" {} \;
find . -type f -name "*.go.bak" -delete

# 3. 更新配置文件
echo "📝 步骤 3/5: 更新配置文件..."
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' "s/your_project/$PROJECT_NAME/g" config/config.toml
    sed -i '' "s/YOUR_PROJECT/${PROJECT_NAME^^}/g" config/config.toml
else
    sed -i "s/your_project/$PROJECT_NAME/g" config/config.toml
    sed -i "s/YOUR_PROJECT/${PROJECT_NAME^^}/g" config/config.toml
fi

# 4. 更新 Makefile
echo "📝 步骤 4/5: 更新 Makefile..."
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' "s/your_project/$PROJECT_NAME/g" Makefile
else
    sed -i "s/your_project/$PROJECT_NAME/g" Makefile
fi

# 5. 更新 docker-compose.yml
echo "📝 步骤 5/5: 更新 Docker 配置..."
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' "s/your_project/$PROJECT_NAME/g" docker-compose.yml
else
    sed -i "s/your_project/$PROJECT_NAME/g" docker-compose.yml
fi

# 下载依赖
echo ""
echo "📦 下载依赖包..."
go mod tidy
go mod download

echo ""
echo "✅ 项目初始化完成！"
echo ""
echo "接下来的步骤:"
echo "  1. 编辑 config/config.toml 配置数据库和 Redis"
echo "  2. 创建数据库表"
echo "  3. 运行项目: make run-api"
echo ""
echo "详细说明请查看 QUICK_START.md"
