@echo off
TITLE Jarvis AI Server
echo Starting Jarvis AI Servers...

echo [1/3] Starting backend server...
START "Jarvis Backend" /MIN cmd /c "cd /d "%~dp0\backend" && go run main.go router.go"

echo [2/3] Starting frontend server...
START "Jarvis Frontend" /MIN cmd /c "cd /d "%~dp0\frontend" && npm run dev"

echo [3/3] Waiting for servers to initialize...
timeout /t 4 /nobreak > nul

echo Opening Jarvis as a Desktop App...
:: You can change 'msedge' to 'chrome' if you prefer Google Chrome
start msedge --app=http://localhost:5173

echo.
echo ===================================================
echo Jarvis is running!
echo The UI has opened in a dedicated app window.
echo.
echo Note: Two minimized terminal windows are running the servers.
echo To completely stop Jarvis, close those terminal windows.
echo ===================================================
pause
