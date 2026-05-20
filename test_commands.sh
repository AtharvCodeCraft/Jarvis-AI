#!/bin/bash

# JARVIS Command Testing Script
# Tests all commands via REST API

API_URL="http://localhost:8080/command"

echo "🤖 JARVIS Command Testing Suite"
echo "================================"
echo ""

# Function to test a command
test_command() {
    local cmd="$1"
    echo "Testing: $cmd"
    response=$(curl -s -X POST $API_URL \
        -H "Content-Type: application/json" \
        -d "{\"text\":\"$cmd\"}")
    echo "Response: $response"
    echo ""
    sleep 2
}

echo "📱 Application Control Tests"
echo "----------------------------"
test_command "open calculator"
test_command "open notepad"

echo "🔍 Web Search Tests"
echo "-------------------"
test_command "search Go programming on Google"

echo "💬 WhatsApp Tests"
echo "-----------------"
test_command "open whatsapp"

echo "🌐 Web Navigation Tests"
echo "-----------------------"
test_command "open youtube"
test_command "open gmail"

echo "📁 Folder Access Tests"
echo "----------------------"
test_command "open downloads folder"
test_command "open documents"

echo "🎵 Media Control Tests"
echo "----------------------"
test_command "play music"
test_command "pause music"

echo "🔊 Volume Control Tests"
echo "-----------------------"
test_command "volume up"
test_command "volume down"

echo "✅ Testing Complete!"
echo ""
echo "Note: Some commands (like shutdown/restart) are not tested automatically for safety."
echo "Test them manually: 'Jarvis lock the computer' is safe to test."
