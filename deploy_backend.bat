@echo off
chcp 65001 >nul
setlocal

echo ==============================================
echo          Go-Zero ERP 后端服务启动脚本
echo ==============================================

cd /d "%~dp0admin"

echo.
echo [1/3] 检查并安装依赖...
go mod tidy

if %errorlevel% neq 0 (
    echo ❌ 依赖安装失败！
    pause
    exit /b 1
)

echo ✅ 依赖安装成功

echo.
echo [2/3] 编译后端服务...
go build -o admin.exe

if %errorlevel% neq 0 (
    echo ❌ 编译失败！
    pause
    exit /b 1
)

echo ✅ 编译成功

echo.
echo [3/3] 启动后端服务...
echo 服务端口: 8001
echo 配置文件: etc/admin-api.yaml
echo.
start cmd /k "admin.exe -f etc/admin-api.yaml"

echo ✅ 后端服务已启动
echo.
echo 按任意键退出...
pause >nul