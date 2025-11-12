@echo off
setlocal enabledelayedexpansion

set "SCRIPT_DIR=%~dp0"
if "%SCRIPT_DIR:~-1%"=="\" set "SCRIPT_DIR=%SCRIPT_DIR:~0,-1%"
set "PROJ=%SCRIPT_DIR%"
set "OUT_BASE=%PROJ%\Release"
set "DEPS=%PROJ%\deps"
set "GLFW_VERSION=3.4"
set "GLFW_PACKAGE=glfw-%GLFW_VERSION%.bin.WIN64"
set "GLFW_URL=https://github.com/glfw/glfw/releases/download/%GLFW_VERSION%/%GLFW_PACKAGE%.zip"
set "GLFW_ARCHIVE=%DEPS%\%GLFW_PACKAGE%.zip"
set "GLFW_DIR=%DEPS%\%GLFW_PACKAGE%"
set "VAR_PATH=example.com/go-ogl-engine/cmd/game"
set "AVAILABLE_PROJECTS=daylight fallline"

set "PROJECT_NAME=%~1"
if "%PROJECT_NAME%"=="" (
    echo Usage: %~nx0 ^<project-name^> [output-name]
    echo.
    echo Available projects:
    for %%P in (%AVAILABLE_PROJECTS%) do echo   %%P
    exit /b 1
)

set "KNOWN=0"
for %%P in (%AVAILABLE_PROJECTS%) do (
    if /I "%PROJECT_NAME%"=="%%P" (
        set "PROJECT_NAME=%%P"
        set "KNOWN=1"
    )
)
if "%KNOWN%"=="0" (
    echo [ERROR] Unknown project %PROJECT_NAME%.
    echo Available projects:
    for %%P in (%AVAILABLE_PROJECTS%) do echo   %%P
    exit /b 1
)

set "OUTPUT_NAME=%~2"
if "%OUTPUT_NAME%"=="" set "OUTPUT_NAME=%PROJECT_NAME%"

set "TARGET_DIR=%OUT_BASE%\%OUTPUT_NAME%"
set "BIN=%TARGET_DIR%"

if not exist "%OUT_BASE%" mkdir "%OUT_BASE%"
if exist "%TARGET_DIR%" rmdir /S /Q "%TARGET_DIR%"
mkdir "%TARGET_DIR%"

if not exist "%DEPS%" mkdir "%DEPS%"

echo === Checking Go toolchain ===
where go >nul 2>nul
if errorlevel 1 (
    echo [ERROR] Go compiler not found in PATH.
    exit /b 1
)

echo === Checking C compiler (gcc/MinGW-w64) ===
set "MINGW_BIN="
for /f "delims=" %%I in ('where gcc 2^>nul') do (
    set "MINGW_BIN=%%~dpI"
    goto :foundGCC
)

echo [ERROR] gcc compiler not found. Install MinGW-w64 (MSYS2) or LLVM with gcc shim and add it to PATH.
exit /b 1

:foundGCC
if "%MINGW_BIN:~-1%"=="\" set "MINGW_BIN=%MINGW_BIN:~0,-1%"
set "PATH=%MINGW_BIN%;%PATH%"
set "CC=%MINGW_BIN%\gcc.exe"

set "CGO_CFLAGS=-I%GLFW_DIR%\include"
set "CGO_LDFLAGS=-L%GLFW_DIR%\lib-mingw-w64 -lglfw3dll -lopengl32 -lgdi32 -luser32 -lshell32 -lkernel32"

if not exist "%GLFW_DIR%" (
    echo === Ensure GLFW %GLFW_VERSION% binaries ===
    echo Downloading GLFW prebuilt package...
    powershell -NoLogo -NoProfile -Command "try { Invoke-WebRequest -Uri '%GLFW_URL%' -OutFile '%GLFW_ARCHIVE%' -ErrorAction Stop } catch { exit 1 }"
    if errorlevel 1 (
        echo [ERROR] Failed to download GLFW archive.
        exit /b 1
    )
    powershell -NoLogo -NoProfile -Command "try { Expand-Archive -LiteralPath '%GLFW_ARCHIVE%' -DestinationPath '%DEPS%' -Force -ErrorAction Stop } catch { exit 1 }"
    if errorlevel 1 (
        echo [ERROR] Failed to extract GLFW archive.
        exit /b 1
    )
    if exist "%GLFW_ARCHIVE%" del "%GLFW_ARCHIVE%"
)

if not exist "%GLFW_DIR%\include" (
    echo [ERROR] GLFW include directory not found: %GLFW_DIR%\include
    exit /b 1
)

echo === Fetch Go module dependencies ===
go mod download
if errorlevel 1 (
    echo [ERROR] go mod download failed.
    exit /b 1
)

echo === Building project %PROJECT_NAME% (Windows amd64) ===
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=1
set "LDFLAGS=-s -w -X %VAR_PATH%.defaultProject=%PROJECT_NAME%"

go build -ldflags "%LDFLAGS%" -o "%BIN%\%OUTPUT_NAME%.exe" ".\cmd\game"
if errorlevel 1 (
    echo [ERROR] go build failed.
    exit /b 1
)

echo === Copy assets ===
xcopy /E /I /Y "%PROJ%\assets" "%BIN%\assets" >nul
if not exist "%BIN%\assets" (
    echo [ERROR] Assets copy failed.
    exit /b 1
)

echo === Copy runtime DLLs ===
if exist "%GLFW_DIR%\lib-mingw-w64\glfw3.dll" (
    copy "%GLFW_DIR%\lib-mingw-w64\glfw3.dll" "%BIN%\" >nul
) else (
    echo [WARN] %GLFW_DIR%\lib-mingw-w64\glfw3.dll not found. Please copy it manually.
)

for %%F in (libstdc++-6.dll libgcc_s_seh-1.dll libwinpthread-1.dll) do (
    if exist "%MINGW_BIN%\%%F" (
        copy "%MINGW_BIN%\%%F" "%BIN%\" >nul
    ) else (
        echo [WARN] %%F not found in %MINGW_BIN%.
    )
)

echo.
echo === Build complete ===
echo Binary: %BIN%\%OUTPUT_NAME%.exe
echo Assets: %BIN%\assets
echo Default project: %PROJECT_NAME%

echo Done.
endlocal
