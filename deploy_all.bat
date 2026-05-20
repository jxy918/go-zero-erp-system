@echo off
chcp 65001 >nul
setlocal

echo ==============================================
echo          Go-Zero ERP 一键启动脚本
echo ==============================================

echo.
echo 正在启动后端服务...
start "" "deploy_backend.bat"

echo.
echo 等待后端启动...
timeout /t 3 /nobreak >nul

echo.
echo 正在启动前端服务...
start "" "deploy_frontend.bat"

echo.
echo ✅ 服务启动完成！
echo.
echo 后端服务: http://localhost:8001
echo 前端服务: http://localhost:3000
echo 默认账号: admin / admin123
echo.
echo 按任意键退出...
pause >nul