# 🤖 JARVIS - Voice Controlled Desktop Assistant

A powerful voice-controlled desktop assistant built with Go and React that can execute Windows system commands, web actions, and provide AI-powered responses.

## ✨ Features

### 1. Application Control
Open any Windows application with voice commands:
- Chrome, VS Code, Notepad, Calculator, File Explorer, Paint, Camera, and more

### 2. Web Actions
- **Google Search**: Automatically search Google from voice commands
- **WhatsApp**: Open WhatsApp Web and initiate calls/chats
- **YouTube**: Open and search YouTube
- **Gmail**: Quick access to Gmail
- **Any Website**: Open any website by name

### 3. System Control
- **Power Management**: Shutdown, restart, lock, sleep
- **Volume Control**: Increase, decrease, mute
- **Folder Access**: Open Downloads, Documents, Pictures, Desktop, etc.

### 4. Media Control
- Play/pause music
- Next/previous track
- Open music player
- Full media key support

### 5. AI Integration
- Local AI responses using Ollama (Qwen 2.5 1.5B model)
- Context-aware conversations
- Memory of past interactions

### 6. Voice Interface
- Speech-to-text using Whisper
- Text-to-speech responses
- Wake word detection (Jarvis)
- Real-time voice feedback

## 🚀 Quick Start

### Prerequisites
- Windows OS
- Go 1.26+ (installed ✓)
- Node.js 16+
- Ollama (for AI features)

### Installation

1. **Install Ollama and AI model:**
   ```bash
   # Download from: https://ollama.ai
   ollama pull qwen2.5:1.5b
   ```

2. **Start Backend:**
   ```bash
   cd backend
   air  # or: go run main.go router.go
   ```

3. **Start Frontend:**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

4. **Open Browser:**
   - Navigate to `http://localhost:5173`
   - Allow microphone access
   - Click mic icon and say "Jarvis open calculator"

## 📖 Documentation

- [QUICKSTART.md](./QUICKSTART.md) - Detailed setup guide
- [COMMANDS.md](./COMMANDS.md) - Complete command reference
- Test scripts: `test_commands.ps1` or `test_commands.sh`

## 🎯 Example Commands

```
"Jarvis open Chrome"
"Jarvis search Python tutorial on Google"
"Jarvis open WhatsApp"
"Jarvis call Rahul on WhatsApp"
"Jarvis shutdown the laptop"
"Jarvis open downloads folder"
"Jarvis play music"
"Jarvis volume up"
"Jarvis open camera"
"Jarvis open YouTube"
```

## 🏗️ Architecture

```
Voice Input → Whisper (STT) → Command Parser → Intent Detection
                                      ↓
                              Action Executor
                                      ↓
                    ┌─────────────────┼─────────────────┐
                    ↓                 ↓                 ↓
              System APIs        Web Actions      Ollama AI
                    ↓                 ↓                 ↓
                    └─────────────────┼─────────────────┘
                                      ↓
                              TTS Response → Voice Output
```

## 🛠️ Technology Stack

### Backend (Go)
- **Framework**: Gin (HTTP/WebSocket)
- **Database**: SQLite (conversation memory)
- **AI**: Ollama API integration
- **Voice**: Whisper.cpp, Windows SAPI TTS

### Frontend (React)
- **Framework**: React + Vite
- **UI**: Custom components with Three.js face
- **Communication**: WebSocket for real-time updates
- **Voice**: Web Speech API

## 📁 Project Structure

```
jarvis-ai/
├── backend/
│   ├── main.go                 # Entry point
│   ├── router.go               # HTTP/WebSocket routes
│   ├── modules/
│   │   ├── ai/                 # Ollama integration
│   │   ├── commands/           # Command parser
│   │   ├── system/             # System control
│   │   │   ├── apps.go         # App launcher
│   │   │   ├── web_actions.go  # Web automation
│   │   │   ├── media_control.go # Media keys
│   │   │   ├── folders.go      # Folder access
│   │   │   └── system_control.go # Power/volume
│   │   ├── voice/              # TTS/STT/Wake word
│   │   └── database/           # SQLite
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── components/         # React components
│   │   ├── hooks/              # Custom hooks
│   │   └── pages/              # Dashboard
│   └── package.json
├── COMMANDS.md                 # Command reference
├── QUICKSTART.md              # Setup guide
└── README.md                  # This file
```

## 🧪 Testing

### Test via API:
```bash
# PowerShell
.\test_commands.ps1

# Bash
bash test_commands.sh

# Manual curl
curl -X POST http://localhost:8080/command \
  -H "Content-Type: application/json" \
  -d '{"text":"open calculator"}'
```

### Test via Web Interface:
1. Open `http://localhost:5173`
2. Click microphone icon
3. Speak commands
4. Watch Jarvis respond

## 🔧 Customization

### Add New Commands:
Edit `backend/modules/commands/command_parser.go`

### Add New Apps:
Edit `backend/modules/system/apps.go`

### Customize UI:
Modify `frontend/src/components/`

### Change AI Model:
Edit `backend/modules/ai/ollama.go` (change model name)

## 🐛 Troubleshooting

### WhatsApp Opens Browser Instead of App?
✓ This is correct! WhatsApp Web is used for better automation

### Voice Not Working?
- Check microphone permissions
- Use HTTPS or localhost
- Check browser console

### AI Not Responding?
- Ensure Ollama is running: `ollama serve`
- Check model: `ollama list`
- Install if missing: `ollama pull qwen2.5:1.5b`

### Commands Not Executing?
- Check backend logs
- Verify Windows permissions
- Test with API first (curl)

## 🤝 Contributing

1. Fork the repository
2. Create feature branch
3. Make changes
4. Test thoroughly
5. Submit pull request

## 📝 License

MIT License - feel free to use and modify

## 🙏 Acknowledgments

- Ollama for local AI
- Whisper for speech recognition
- Go Gin framework
- React and Vite

---

**Made with ❤️ for voice-controlled productivity**
"# Jarvis-AI" 
