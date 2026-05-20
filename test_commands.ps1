# JARVIS Command Testing Script (PowerShell)
# Tests all commands via REST API

$API_URL = "http://localhost:8080/command"

Write-Host "🤖 JARVIS Command Testing Suite" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""

# Function to test a command
function Test-Command {
    param($cmd)
    Write-Host "Testing: $cmd" -ForegroundColor Yellow
    
    $body = @{
        text = $cmd
    } | ConvertTo-Json
    
    try {
        $response = Invoke-RestMethod -Uri $API_URL -Method Post -Body $body -ContentType "application/json"
        Write-Host "Response: $($response.result)" -ForegroundColor Green
    } catch {
        Write-Host "Error: $_" -ForegroundColor Red
    }
    
    Write-Host ""
    Start-Sleep -Seconds 2
}

Write-Host "📱 Application Control Tests" -ForegroundColor Magenta
Write-Host "----------------------------"
Test-Command "open calculator"
Test-Command "open notepad"

Write-Host "🔍 Web Search Tests" -ForegroundColor Magenta
Write-Host "-------------------"
Test-Command "search Go programming on Google"

Write-Host "💬 WhatsApp Tests" -ForegroundColor Magenta
Write-Host "-----------------"
Test-Command "open whatsapp"

Write-Host "🌐 Web Navigation Tests" -ForegroundColor Magenta
Write-Host "-----------------------"
Test-Command "open youtube"
Test-Command "open gmail"

Write-Host "📁 Folder Access Tests" -ForegroundColor Magenta
Write-Host "----------------------"
Test-Command "open downloads folder"
Test-Command "open documents"

Write-Host "🎵 Media Control Tests" -ForegroundColor Magenta
Write-Host "----------------------"
Test-Command "play music"
Test-Command "pause music"

Write-Host "🔊 Volume Control Tests" -ForegroundColor Magenta
Write-Host "-----------------------"
Test-Command "volume up"
Test-Command "volume down"

Write-Host "✅ Testing Complete!" -ForegroundColor Green
Write-Host ""
Write-Host "Note: Some commands (like shutdown/restart) are not tested automatically for safety." -ForegroundColor Yellow
Write-Host "Test them manually: 'Jarvis lock the computer' is safe to test." -ForegroundColor Yellow
