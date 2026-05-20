@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

title Go-Zero ERP 前端启动器

echo ==============================================
echo          Go-Zero ERP 前端服务启动脚本
echo ==============================================

REM 切换到 frontend 目录
cd /d "%~dp0frontend" 2>nul
if %errorlevel% neq 0 (
    echo [错误] 找不到 frontend 目录，请确认脚本位于项目根目录。
    echo 按任意键退出...
    pause >nul
    exit /b 1
)

REM 检查 node_modules 是否存在（若不存在则自动安装）
if not exist "node_modules" (
    echo.
    echo [步骤 1/2] 安装依赖中，请稍候...
    call npm install
    if !errorlevel! neq 0 (
        echo [错误] npm install 失败，请检查网络或 package.json
        echo 按任意键退出...
        pause >nul
        exit /b 1
    )
    echo [成功] 依赖安装完成
) else (
    echo.
    echo [跳过] 依赖已存在，无需安装
)

echo.
echo [步骤 2/2] 启动开发服务器...
echo 服务地址: http://localhost:3000
echo 代理目标: http://localhost:8001
echo.
echo 提示: 按 Ctrl+C 可停止服务，关闭此窗口也会终止服务
echo ==============================================
echo.

REM 关键点：使用 call 执行，即使失败也不会导致窗口立即关闭
call npm run dev

REM 如果上面命令执行完毕（说明 dev 服务器已退出），则给出提示并等待
echo.
echo [警告] 开发服务器已停止运行，可能是启动失败或人为关闭。
echo 若启动失败，请检查端口 3000 是否被占用，或手动执行 "npm run dev" 查看详细错误。
echo 按任意键退出...
pause >nul
exit /b 0