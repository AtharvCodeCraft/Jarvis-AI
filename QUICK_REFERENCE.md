# 🎯 JARVIS Quick Reference Card

## 🚀 Start Commands

```bash
# Terminal 1: Start Backend
cd backend && air

# Terminal 2: Start Frontend  
cd frontend && npm run dev

# Terminal 3: Start Ollama (optional, for AI)
ollama serve
```

## 💬 Voice Commands Cheat Sheet

### 📱 Apps
| Say This | Opens |
|----------|-------|
| "Jarvis open Chrome" | Google Chrome |
| "Jarvis open VS Code" | Visual Studio Code |
| "Jarvis open Calculator" | Calculator |
| "Jarvis open Notepad" | Notepad |
| "Jarvis open Camera" | Windows Camera |
| "Jarvis open File Explorer" | File Explorer |

### 🌐 Web
| Say This | Action |
|----------|--------|
| "Jarvis search [query] on Google" | Google Search |
| "Jarvis open WhatsApp" | WhatsApp Web |
| "Jarvis call [name] on WhatsApp" | WhatsApp (manual call) |
| "Jarvis open YouTube" | YouTube |
| "Jarvis search [query] on YouTube" | YouTube Search |
| "Jarvis open Gmail" | Gmail |

### 📁 Folders
| Say This | Opens |
|----------|-------|
| "Jarvis open downloads" | Downloads folder |
| "Jarvis open documents" | Documents folder |
| "Jarvis open pictures" | Pictures folder |
| "Jarvis open desktop" | Desktop folder |

### 🎵 Media
| Say This | Action |
|----------|--------|
| "Jarvis play music" | Play/Resume |
| "Jarvis pause music" | Pause |
| "Jarvis next song" | Next track |
| "Jarvis previous song" | Previous track |

### 🔊 Volume
| Say This | Action |
|----------|--------|
| "Jarvis volume up" | Increase volume |
| "Jarvis volume down" | Decrease volume |
| "Jarvis mute" | Toggle mute |

### ⚡ Power
| Say This | Action |
|----------|--------|
| "Jarvis lock the computer" | Lock screen |
| "Jarvis put system to sleep" | Sleep mode |
| "Jarvis shutdown the laptop" | Shutdown |
| "Jarvis restart the system" | Restart |

## 🧪 Test Without Voice

```bash
# PowerShell
curl -X POST http://localhost:8080/command `
  -H "Content-Type: application/json" `
  -d '{"text":"open calculator"}'

# Bash
curl -X POST http://localhost:8080/command \
  -H "Content-Type: application/json" \
  -d '{"text":"open calculator"}'
```

## 🔧 Quick Troubleshooting

| Problem | Solution |
|---------|----------|
| WhatsApp not opening | ✓ Opens in browser (WhatsApp Web) - this is correct! |
| Voice not working | Check mic permissions in browser |
| AI not responding | Run `ollama serve` in terminal |
| Command not executing | Check backend logs for errors |
| Volume keys not working | Ensure no other app is blocking media keys |

## 📊 System Status

Check if everything is running:

```bash
# Check backend
curl http://localhost:8080/system/status

# Check frontend
# Open: http://localhost:5173

# Check Ollama
curl http://localhost:11434/api/tags
```

## 🎨 URLs

- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:8080
- **WebSocket**: ws://localhost:8080/ws
- **Ollama**: http://localhost:11434

## 📝 File Locations

```
backend/modules/commands/command_parser.go  → Add new commands
backend/modules/system/apps.go              → Add new apps
backend/modules/system/web_actions.go       → Add web actions
frontend/src/components/                    → Customize UI
```

## 🎯 Most Common Commands

1. `"Jarvis open Chrome"`
2. `"Jarvis search [topic] on Google"`
3. `"Jarvis open downloads"`
4. `"Jarvis volume up"`
5. `"Jarvis play music"`
6. `"Jarvis lock the computer"`
7. `"Jarvis open WhatsApp"`
8. `"Jarvis open YouTube"`

---

**Print this and keep it handy! 📄**
