@echo off

:loop
echo Welcome %ShellUser% shell as %UserName%

E:\opt\msys64\usr\bin\bash -l

:read
set quit=nil

set /p quit="Are you sure you want to close the shell <yes/no> : "

if "%quit%"=="y" goto end
if "%quit%"=="yes" goto end

if "%quit%"=="n" goto loop
if "%quit%"=="no" goto loop

goto read

:end
echo bye %ShellUser% of %UserName%

timeout /t 1