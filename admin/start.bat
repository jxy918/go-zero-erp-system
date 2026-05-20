@echo off
title Admin API Server
cd /d d:\work\go\src\go-zero-erp\admin
echo Starting Admin API Server...
.\admin.exe -f etc/admin-api.yaml
pause