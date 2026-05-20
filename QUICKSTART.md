# JARVIS Quick Start Guide

## Prerequisites

1. **Install Go** (already installed - Go 1.26.1)
2. **Install Node.js** (for frontend)
3. **Install Ollama** (for AI responses)
   ```bash
   # Download from: https://ollama.ai
   # After installation, run:
   ollama pull qwen2.5:1.5b
   ```

## Starting JARVIS

### Option 1: Using Air (Hot Reload - Recommended for Development)

1. **Start Backend with Air:**
   ```bash
   cd backend
   air
   ```

2. **Start Frontend:**
   ```bash
   cd frontend
   npm install  # First time only
   npm run dev
   ```

3. **Open Browser:**
   - Navigate to `http://localhost:5173`

### Option 2: Manual Start

1. **Start Ollama:**
   ```bash
   ollama serve
   ```

2. **Start Backend:**
   ```bash
   cd backend
   go run main.go router.go
   ```

3. **Start Frontend:**
   ```bash
   cd frontend
   npm run dev
   ```

## Testing Commands

### Test via Web Interface:
1. Open `http://localhost:5173`
2. Click the microphone icon
3. Say: "Jarvis open calculator"
4. Jarvis will execute and respond

### Test via API (without voice):
```bash
curl -X POST http://localhost:8080/command \
  -H "Content-Type: application/json" \
  -d '{"text":"open calculator"}'
```

## Example Commands to Try

1. **Open Apps:**
   - "Jarvis open Chrome"
   - "Jarvis open Calculator"
   - "Jarvis open Notepad"

2. **Web Actions:**
   - "Jarvis search Python tutorial on Google"
   - "Jarvis open YouTube"
   - "Jarvis open WhatsApp"

3. **System Control:**
   - "Jarvis volume up"
   - "Jarvis lock the computer"
   - "Jarvis open downloads folder"

4. **Media Control:**
   - "Jarvis play music"
   - "Jarvis next song"
   - "Jarvis pause"

## Troubleshooting

### WhatsApp Not Opening?
- WhatsApp is opened via web browser (WhatsApp Web)
- Make sure you have a browser installed (Chrome, Edge, Firefox)
- The command will open `https://web.whatsapp.com`

### Voice Not Working?
- Check microphone permissions in browser
- Ensure you're using HTTPS or localhost
- Check browser console for errors

### AI Not Responding?
- Make sure Ollama is running: `ollama serve`
- Check if model is installed: `ollama list`
- Install model if missing: `ollama pull qwen2.5:1.5b`

### Volume/Media Keys Not Working?
- These use Windows keyboard shortcuts
- Make sure no other app is blocking media keys
- Try testing with Windows Media Player or Spotify open

## Architecture

```
┌─────────────────┐
│   Frontend      │
│  (React + Vite) │
│   Port: 5173    │
└────────┬────────┘
         │ WebSocket
         ↓
┌─────────────────┐
│   Backend       │
│   (Go + Gin)    │
│   Port: 8080    │
└────────┬────────┘
         │
    ┌────┴────┬──────────┬──────────┐
    ↓         ↓          ↓          ↓
┌────────┐ ┌──────┐ ┌────────┐ ┌────────┐
│ Ollama │ │ TTS  │ │Whisper │ │Windows │
│  AI    │ │Voice │ │ STT    │ │  APIs  │
└────────┘ └──────┘ └────────┘ └────────┘
```

## Next Steps

1. **Customize Commands:** Edit `backend/modules/commands/command_parser.go`
2. **Add New Apps:** Update `backend/modules/system/apps.go`
3. **Improve UI:** Modify `frontend/src/components/`
4. **Add Wake Word:** Implement in `backend/modules/voice/wakeword.go`

## Full Command List

See [COMMANDS.md](./COMMANDS.md) for complete list of available commands.
