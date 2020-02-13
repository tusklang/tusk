@echo off

node "%~dp0src\cli.js" %cd% %* --max-old-space-size=4096
