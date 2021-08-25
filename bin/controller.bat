@echo off
set CURRENT_DIR=%~dp0
cd %CURRENT_DIR%
%CURRENT_DIR:~0,2%

:AGAIN
echo * 0 cancel execute
echo * 1 install master
echo * 2 install slave
echo * 3 uninstall master
echo * 4 uninstall slave

set choice=-1
set /p choice=choice commnad: 
if "%choice%"=="0" (
    goto END
)
if "%choice%"=="1" (
    webpc-master install
    pause
    goto END
)
if "%choice%"=="2" (
    webpc-slave install
    pause
    goto END
)
if "%choice%"=="3" (
    webpc-master uninstall
    pause
    goto END
)
if "%choice%"=="4" (
    webpc-slave uninstall
    pause
    goto END
)
goto AGAIN

:END
