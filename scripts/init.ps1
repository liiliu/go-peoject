# 项目初始化脚本 (PowerShell)
# 使用方法: .\scripts\init.ps1 your_project_name

param(
    [Parameter(Mandatory=$true)]
    [string]$ProjectName
)

Write-Host "🚀 正在初始化项目: $ProjectName" -ForegroundColor Green
Write-Host ""

# 1. 替换模块名称
Write-Host "📝 步骤 1/5: 更新 go.mod 模块名称..." -ForegroundColor Cyan
(Get-Content go.mod) -replace 'your_project', $ProjectName | Set-Content go.mod

# 2. 替换所有源代码中的 import 路径
Write-Host "📝 步骤 2/5: 更新源代码中的 import 路径..." -ForegroundColor Cyan
Get-ChildItem -Path . -Filter *.go -Recurse | ForEach-Object {
    (Get-Content $_.FullName) -replace 'your_project', $ProjectName | Set-Content $_.FullName
}

# 3. 更新配置文件
Write-Host "📝 步骤 3/5: 更新配置文件..." -ForegroundColor Cyan
$ProjectNameUpper = $ProjectName.ToUpper()
(Get-Content config\config.toml) -replace 'your_project', $ProjectName -replace 'YOUR_PROJECT', $ProjectNameUpper | Set-Content config\config.toml

# 4. 更新 Makefile
Write-Host "📝 步骤 4/5: 更新 Makefile..." -ForegroundColor Cyan
(Get-Content Makefile) -replace 'your_project', $ProjectName | Set-Content Makefile

# 5. 更新 docker-compose.yml
Write-Host "📝 步骤 5/5: 更新 Docker 配置..." -ForegroundColor Cyan
(Get-Content docker-compose.yml) -replace 'your_project', $ProjectName | Set-Content docker-compose.yml

# 下载依赖
Write-Host ""
Write-Host "📦 下载依赖包..." -ForegroundColor Cyan
go mod tidy
go mod download

Write-Host ""
Write-Host "✅ 项目初始化完成！" -ForegroundColor Green
Write-Host ""
Write-Host "接下来的步骤:" -ForegroundColor Yellow
Write-Host "  1. 编辑 config/config.toml 配置数据库和 Redis"
Write-Host "  2. 创建数据库表"
Write-Host "  3. 运行项目: make run-api 或 go run main.go api"
Write-Host ""
Write-Host "详细说明请查看 QUICK_START.md"
