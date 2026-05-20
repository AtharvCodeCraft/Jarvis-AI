# 🔧 JARVIS Troubleshooting Guide

## Common Issues and Solutions

### 1. WhatsApp Opens in Browser Instead of Desktop App

**Issue**: When I say "Jarvis open WhatsApp", it opens in my browser instead of the desktop app.

**Solution**: ✅ This is the correct behavior!

**Why?**
- WhatsApp desktop app doesn't support command-line launching on Windows
- WhatsApp Web provides better automation capabilities
- This approach works consistently across all Windows versions
- You can still use all WhatsApp features through the web interface

**Alternative**: If you prefer the desktop app, you can manually open it after Jarvis opens WhatsApp Web.

---

### 2. Voice Recognition Not Working

**Symptoms**:
- Microphone icon doesn't respond
- No text appears when speaking
- Browser doesn't ask for microphone permission

**Solutions**:

1. **Check Browser Permissions**:
   - Chrome: Settings → Privacy and Security → Site Settings → Microphone
   - Edge: Settings → Cookies and site permissions → Microphone
   - Ensure `localhost:5173` has microphone access

2. **Use HTTPS or Localhost**:
   - Web Speech API only works on HTTPS or localhost
   - If accessing from another device, use HTTPS

3. **Check Microphone Hardware**:
   ```powershell
   # Test microphone in Windows
   # Settings → System → Sound → Input → Test your microphone
   ```

4. **Browser Console Errors**:
   - Press F12 to open Developer Tools
   - Check Console tab for errors
   - Look for "NotAllowedError" or "NotFoundError"

---

### 3. Commands Not Executing

**Symptoms**:
- Jarvis responds but nothing happens
- Applications don't open
- System commands don't work

**Solutions**:

1. **Check Backend is Running**:
   ```bash
   curl http://localhost:8080/system/status
   ```
   Should return: `{"cpu":12,"ram":48,"status":"Online"}`

2. **Check Backend Logs**:
   ```bash
   cd backend
   # Look for error messages in the terminal
   ```

3. **Test with API Directly**:
   ```bash
   curl -X POST http://localhost:8080/command \
     -H "Content-Type: application/json" \
     -d '{"text":"open calculator"}'
   ```

4. **Windows Permissions**:
   - Run backend as Administrator if needed
   - Check Windows Defender isn't blocking commands

5. **Application Not Installed**:
   - Ensure the app you're trying to open is installed
   - Try: "Jarvis open calculator" (always available)

---

### 4. AI Not Responding

**Symptoms**:
- "I encountered an error connecting to my AI core"
- Long delays before responses
- No response to conversational queries

**Solutions**:

1. **Check Ollama is Running**:
   ```bash
   # Check if Ollama is running
   curl http://localhost:11434/api/tags
   
   # If not running, start it
   ollama serve
   ```

2. **Check Model is Installed**:
   ```bash
   ollama list
   
   # Should show: qwen2.5:1.5b
   # If not, install it:
   ollama pull qwen2.5:1.5b
   ```

3. **Check Ollama Logs**:
   - Look at the terminal where `ollama serve` is running
   - Check for memory errors or model loading issues

4. **Try Different Model**:
   Edit `backend/modules/ai/ollama.go`:
   ```go
   Model: "llama3.2:1b", // Smaller, faster model
   ```

---

### 5. Volume/Media Keys Not Working

**Symptoms**:
- "Jarvis volume up" doesn't change volume
- Media controls don't work
- No error message but no action

**Solutions**:

1. **Test Media Keys Manually**:
   - Press volume up/down keys on keyboard
   - If they don't work, Windows might be blocking them

2. **Close Conflicting Apps**:
   - Some apps (Spotify, Discord) intercept media keys
   - Try closing them and test again

3. **Test with Media Player Open**:
   ```
   "Jarvis open music player"
   "Jarvis play music"
   ```

4. **Check PowerShell Execution**:
   ```powershell
   # Test volume up manually
   powershell -Command "(New-Object -ComObject WScript.Shell).SendKeys([char]175)"
   ```

---

### 6. Frontend Won't Start

**Symptoms**:
- `npm run dev` fails
- Port 5173 already in use
- Module not found errors

**Solutions**:

1. **Install Dependencies**:
   ```bash
   cd frontend
   rm -rf node_modules package-lock.json
   npm install
   ```

2. **Port Already in Use**:
   ```bash
   # Kill process on port 5173
   netstat -ano | findstr :5173
   taskkill /PID <PID> /F
   
   # Or use different port
   npm run dev -- --port 3000
   ```

3. **Node Version**:
   ```bash
   node --version  # Should be 16+
   npm --version   # Should be 8+
   ```

---

### 7. Backend Won't Compile

**Symptoms**:
- `go build` fails
- Import errors
- Module not found

**Solutions**:

1. **Update Dependencies**:
   ```bash
   cd backend
   go mod tidy
   go mod download
   ```

2. **Check Go Version**:
   ```bash
   go version  # Should be 1.21+
   ```

3. **Clean Build**:
   ```bash
   go clean -cache
   go build -o jarvis.exe
   ```

---

### 8. WebSocket Connection Failed

**Symptoms**:
- "WebSocket connection failed" in console
- Commands don't reach backend
- No real-time updates

**Solutions**:

1. **Check Backend is Running**:
   ```bash
   curl http://localhost:8080/system/status
   ```

2. **Check WebSocket URL**:
   In `frontend/src/hooks/useJarvisSocket.js`:
   ```javascript
   const ws = new WebSocket('ws://localhost:8080/ws');
   ```

3. **Firewall Blocking**:
   - Check Windows Firewall
   - Allow Go application through firewall

4. **Test WebSocket Manually**:
   Use a WebSocket testing tool or browser extension

---

### 9. Specific Commands Not Working

#### "Jarvis open Chrome" - Chrome doesn't open

**Solutions**:
- Ensure Chrome is installed
- Try: "Jarvis open browser"
- Check if Chrome is in PATH

#### "Jarvis open downloads" - Wrong folder opens

**Solutions**:
- Check `USERPROFILE` environment variable:
  ```powershell
  echo $env:USERPROFILE
  ```
- Should be: `C:\Users\YourUsername`

#### "Jarvis shutdown" - Nothing happens

**Solutions**:
- Run backend as Administrator
- Test manually:
  ```powershell
  shutdown /s /t 0
  ```

---

### 10. Performance Issues

**Symptoms**:
- Slow responses
- High CPU usage
- Memory leaks

**Solutions**:

1. **Use Lighter AI Model**:
   ```bash
   ollama pull qwen2.5:0.5b  # Even smaller
   ```

2. **Disable AI Fallback**:
   Comment out AI fallback in `command_parser.go` for faster responses

3. **Check System Resources**:
   ```powershell
   # Open Task Manager
   taskmgr
   ```

4. **Restart Services**:
   ```bash
   # Restart backend
   Ctrl+C in backend terminal
   go run main.go router.go
   
   # Restart frontend
   Ctrl+C in frontend terminal
   npm run dev
   ```

---

## Diagnostic Commands

### Check All Services

```bash
# Backend
curl http://localhost:8080/system/status

# Frontend
curl http://localhost:5173

# Ollama
curl http://localhost:11434/api/tags

# WebSocket (requires wscat)
wscat -c ws://localhost:8080/ws
```

### Test Individual Modules

```bash
# Test command parser
curl -X POST http://localhost:8080/command \
  -H "Content-Type: application/json" \
  -d '{"text":"open calculator"}'

# Test different commands
curl -X POST http://localhost:8080/command \
  -H "Content-Type: application/json" \
  -d '{"text":"volume up"}'

curl -X POST http://localhost:8080/command \
  -H "Content-Type: application/json" \
  -d '{"text":"search test on google"}'
```

### Enable Debug Logging

Add to `backend/main.go`:
```go
log.SetFlags(log.LstdFlags | log.Lshortfile)
```

---

## Getting Help

### 1. Check Logs

**Backend Logs**:
- Look at terminal where backend is running
- Check for error messages and stack traces

**Frontend Logs**:
- Open browser Developer Tools (F12)
- Check Console tab for errors

**Ollama Logs**:
- Look at terminal where `ollama serve` is running

### 2. Test Systematically

1. Test backend API with curl
2. Test WebSocket connection
3. Test voice input in browser
4. Test specific commands one by one

### 3. Verify Installation

```bash
# Check all required tools
go version
node --version
npm --version
ollama --version

# Check services are running
netstat -ano | findstr :8080   # Backend
netstat -ano | findstr :5173   # Frontend
netstat -ano | findstr :11434  # Ollama
```

### 4. Common Error Messages

| Error Message | Cause | Solution |
|---------------|-------|----------|
| "Failed to open X" | App not installed | Install app or use different name |
| "AI Query Error" | Ollama not running | Run `ollama serve` |
| "WebSocket connection failed" | Backend not running | Start backend |
| "NotAllowedError" | Mic permission denied | Allow mic in browser settings |
| "Port already in use" | Service already running | Kill process or use different port |

---

## Still Having Issues?

1. **Run Test Scripts**:
   ```bash
   # PowerShell
   .\test_commands.ps1
   
   # Bash
   bash test_commands.sh
   ```

2. **Check Documentation**:
   - [QUICKSTART.md](./QUICKSTART.md)
   - [COMMANDS.md](./COMMANDS.md)
   - [SYSTEM_DIAGRAM.md](./SYSTEM_DIAGRAM.md)

3. **Reset Everything**:
   ```bash
   # Stop all services (Ctrl+C in each terminal)
   
   # Clean backend
   cd backend
   go clean -cache
   rm jarvis.db
   
   # Clean frontend
   cd frontend
   rm -rf node_modules
   npm install
   
   # Restart everything
   # Terminal 1: ollama serve
   # Terminal 2: cd backend && go run main.go router.go
   # Terminal 3: cd frontend && npm run dev
   ```

---

**Remember**: Most issues are due to services not running or permissions. Always check that backend, frontend, and Ollama are all running!
