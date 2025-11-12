@echo off
setlocal enabledelayedexpansion

REM ==== НАСТРОЙКА ПУТЕЙ (проверь и правь при необходимости) ====
set "MINGW=C:\mingw64"
set "GLFW=C:\glfw-3.4.bin.WIN64"

REM ==== ПАПКИ ПРОЕКТА ====
set "PROJ=%cd%"
set "OUT=%PROJ%\Release"
set "BIN=%OUT%\bin"

echo === Prepare output folders ===
if not exist "%BIN%" mkdir "%BIN%"
if exist "%OUT%\assets" rmdir /S /Q "%OUT%\assets"

REM ==== PATH/КОМПИЛЯТОР ====
set "PATH=%MINGW%\bin;%PATH%"
set "CC=%MINGW%\bin\x86_64-w64-mingw32-gcc.exe"

REM ==== CGO ФЛАГИ ДЛЯ GLFW (динамическая линковка) ====
REM Include с GLFW
set "CGO_CFLAGS=-I%GLFW%\include"
REM LinkDir GLFW + импортная либпа (dll) и системные Win32-библиотеки
REM В дистрибутиве GLFW обычно есть: lib-mingw-w64\glfw3.dll и libglfw3dll.a
set "CGO_LDFLAGS=-L%GLFW%\lib-mingw-w64 -lglfw3dll -lgdi32 -luser32 -lshell32 -lkernel32 -lopengl32"

echo === Build MyGame.exe ===
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=1
go build -ldflags "-s -w" -o "%BIN%\MyGame.exe" ".\cmd\game"
if errorlevel 1 (
  echo [ERROR] go build failed
  exit /b 1
)

echo === Copy assets ===
xcopy /E /I /Y "%PROJ%\assets" "%OUT%\assets" >nul

echo === Copy DLLs (GLFW + MinGW runtimes) ===
REM GLFW dll
if exist "%GLFW%\lib-mingw-w64\glfw3.dll" (
  copy "%GLFW%\lib-mingw-w64\glfw3.dll" "%BIN%\" >nul
) else (
  echo [WARN] %GLFW%\lib-mingw-w64\glfw3.dll not found. Check your GLFW package paths.
)

REM MinGW рантаймы (обычно нужны при линковке с dll)
for %%F in (libstdc++-6.dll libgcc_s_seh-1.dll libwinpthread-1.dll) do (
  if exist "%MINGW%\bin\%%F" (
    copy "%MINGW%\bin\%%F" "%BIN%\" >nul
  ) else (
    echo [WARN] %MINGW%\bin\%%F not found.
  )
)

echo === Make ZIP (optional) ===
powershell -Command "Compress-Archive -Path '%OUT%\*' -DestinationPath 'MyGame_Windows.zip' -Force" 2>nul

echo.
echo === DONE ===
echo Папка portable: %OUT%
echo Запускать: %BIN%\MyGame.exe
endlocal
