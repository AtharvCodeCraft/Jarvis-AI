# JARVIS System Architecture Diagram

## High-Level Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         USER INTERFACE                          │
│                                                                 │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐    │
│  │   Browser    │    │  Microphone  │    │   Speaker    │    │
│  │ (localhost:  │    │   (Voice     │    │   (Voice     │    │
│  │    5173)     │    │    Input)    │    │   Output)    │    │
│  └──────┬───────┘    └──────┬───────┘    └──────▲───────┘    │
│         │                   │                    │             │
└─────────┼───────────────────┼────────────────────┼─────────────┘
          │                   │                    │
          │ HTTP/WebSocket    │ Speech-to-Text     │ Text-to-Speech
          │                   │                    │
┌─────────▼───────────────────▼────────────────────┼─────────────┐
│                      FRONTEND (React)             │             │
│                                                   │             │
│  ┌────────────────┐  ┌────────────────┐  ┌───────┴──────┐    │
│  │  JarvisFace    │  │ MicIndicator   │  │ CommandConsole│   │
│  │  (Three.js)    │  │                │  │              │    │
│  └────────────────┘  └────────────────┘  └──────────────┘    │
│                                                                 │
│  ┌────────────────────────────────────────────────────────┐   │
│  │           useJarvisSocket (WebSocket Hook)             │   │
│  └────────────────────────┬───────────────────────────────┘   │
└───────────────────────────┼─────────────────────────────────────┘
                            │ WebSocket (ws://localhost:8080/ws)
                            │
┌───────────────────────────▼─────────────────────────────────────┐
│                    BACKEND (Go + Gin)                           │
│                                                                 │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │                    router.go                             │ │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────────┐    │ │
│  │  │ /command   │  │ /ws        │  │ /system/status │    │ │
│  │  │ (POST)     │  │ (WebSocket)│  │ (GET)          │    │ │
│  │  └─────┬──────┘  └─────┬──────┘  └────────────────┘    │ │
│  └────────┼───────────────┼─────────────────────────────────┘ │
│           │               │                                    │
│  ┌────────▼───────────────▼─────────────────────────────────┐ │
│  │         modules/commands/command_parser.go               │ │
│  │                                                          │ │
│  │  ParseAndExecute(input string) → string                 │ │
│  │                                                          │ │
│  │  • Remove "jarvis" prefix                               │ │
│  │  • Detect intent (app, web, media, power, etc.)        │ │
│  │  • Extract parameters (names, queries, folders)        │ │
│  │  • Route to appropriate handler                        │ │
│  └────────┬─────────────────────────────────────────────────┘ │
│           │                                                    │
│  ┌────────▼─────────────────────────────────────────────────┐ │
│  │              modules/system/ (Handlers)                  │ │
│  │                                                          │ │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │ │
│  │  │   apps.go    │  │web_actions.go│  │  folders.go  │ │ │
│  │  │              │  │              │  │              │ │ │
│  │  │ OpenApp()    │  │GoogleSearch()│  │OpenFolder()  │ │ │
│  │  │              │  │OpenWhatsApp()│  │              │ │ │
│  │  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘ │ │
│  │         │                 │                  │         │ │
│  │  ┌──────▼───────┐  ┌──────▼───────┐  ┌──────▼───────┐ │ │
│  │  │media_control │  │system_control│  │              │ │ │
│  │  │              │  │              │  │              │ │ │
│  │  │MediaControl()│  │PowerControl()│  │              │ │ │
│  │  │              │  │ControlVolume│  │              │ │ │
│  │  └──────┬───────┘  └──────┬───────┘  └──────────────┘ │ │
│  └─────────┼──────────────────┼────────────────────────────┘ │
│            │                  │                               │
│  ┌─────────▼──────────────────▼────────────────────────────┐ │
│  │              os/exec (Windows Commands)                 │ │
│  │                                                          │ │
│  │  cmd /c start [app]         → Launch applications       │ │
│  │  powershell SendKeys        → Media/Volume keys         │ │
│  │  shutdown /s /r             → Power commands            │ │
│  │  rundll32.exe               → Lock/Sleep                │ │
│  └────────┬─────────────────────────────────────────────────┘ │
└───────────┼───────────────────────────────────────────────────┘
            │
┌───────────▼───────────────────────────────────────────────────┐
│                    WINDOWS OPERATING SYSTEM                   │
│                                                               │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │  Apps    │  │ Browser  │  │  Media   │  │  Power   │   │
│  │          │  │          │  │  Player  │  │  Mgmt    │   │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘   │
└───────────────────────────────────────────────────────────────┘


┌───────────────────────────────────────────────────────────────┐
│                    OPTIONAL: AI MODULE                        │
│                                                               │
│  ┌────────────────────────────────────────────────────────┐  │
│  │              modules/ai/ollama.go                      │  │
│  │                                                        │  │
│  │  QueryOllama(prompt) → AI Response                    │  │
│  │                                                        │  │
│  │  • Fallback for unknown commands                      │  │
│  │  • Conversational queries                             │  │
│  │  • Context-aware responses                            │  │
│  └────────────────────┬───────────────────────────────────┘  │
│                       │ HTTP (localhost:11434)               │
│  ┌────────────────────▼───────────────────────────────────┐  │
│  │              Ollama Server                             │  │
│  │         (qwen2.5:1.5b model)                          │  │
│  └────────────────────────────────────────────────────────┘  │
└───────────────────────────────────────────────────────────────┘
```

## Data Flow Example: "Jarvis open Chrome"

```
1. User speaks: "Jarvis open Chrome"
   │
2. Browser captures audio → Speech-to-Text
   │
3. Frontend sends via WebSocket: {"type": "voice_command", "text": "jarvis open chrome"}
   │
4. Backend router.go receives WebSocket message
   │
5. Calls: commands.ParseAndExecute("jarvis open chrome")
   │
6. command_parser.go:
   - Removes "jarvis" prefix → "open chrome"
   - Detects "open" command
   - Extracts app name → "chrome"
   - Routes to: system.OpenApp("chrome")
   │
7. apps.go:
   - Matches "chrome" case
   - Executes: cmd /c start chrome
   - Returns: "Opening chrome"
   │
8. Response sent back via WebSocket
   │
9. Frontend displays: "Opening chrome"
   │
10. TTS speaks: "Opening chrome"
    │
11. Chrome launches on Windows
```

## Command Routing Logic

```
Input: "jarvis [command]"
         │
         ├─ Contains "whatsapp" → web_actions.go
         │                         ├─ Contains "call" → WhatsAppCall()
         │                         └─ Else → OpenWhatsApp()
         │
         ├─ Contains "search" + "google" → GoogleSearch()
         │
         ├─ Contains "youtube" → YouTubeSearch() or OpenYouTube()
         │
         ├─ Contains "folder" or folder names → folders.go → OpenFolder()
         │
         ├─ Contains "music/play/pause/next" → media_control.go → MediaControl()
         │
         ├─ Contains "volume/mute" → system_control.go → ControlVolume()
         │
         ├─ Contains "shutdown/restart/lock/sleep" → PowerControl()
         │
         ├─ Starts with "open" → apps.go → OpenApp()
         │
         └─ Else → ai/ollama.go → QueryOllama() (AI fallback)
```

## Module Dependencies

```
main.go
  │
  ├─ router.go
  │   ├─ commands/command_parser.go
  │   │   ├─ system/apps.go
  │   │   ├─ system/web_actions.go
  │   │   ├─ system/folders.go
  │   │   ├─ system/media_control.go
  │   │   ├─ system/system_control.go
  │   │   └─ ai/ollama.go
  │   │       └─ ai/memory.go
  │   └─ system/system_control.go (GetStatus)
  │
  ├─ database/sqlite.go
  └─ plugins/plugin_loader.go
```

## Technology Stack

```
┌─────────────────────────────────────────────────────────┐
│                    FRONTEND STACK                       │
├─────────────────────────────────────────────────────────┤
│ React 18          │ UI Framework                        │
│ Vite              │ Build Tool & Dev Server             │
│ Three.js          │ 3D Jarvis Face Animation            │
│ Web Speech API    │ Browser Speech Recognition          │
│ WebSocket         │ Real-time Communication             │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│                    BACKEND STACK                        │
├─────────────────────────────────────────────────────────┤
│ Go 1.26           │ Programming Language                │
│ Gin               │ HTTP Framework                      │
│ Gorilla WebSocket │ WebSocket Support                   │
│ SQLite            │ Database (modernc.org/sqlite)       │
│ os/exec           │ System Command Execution            │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│                      AI STACK                           │
├─────────────────────────────────────────────────────────┤
│ Ollama            │ Local LLM Server                    │
│ Qwen 2.5 1.5B     │ Lightweight AI Model                │
│ Whisper (planned) │ Speech-to-Text                      │
│ Windows SAPI      │ Text-to-Speech                      │
└─────────────────────────────────────────────────────────┘
```

## Port Allocation

| Service | Port | Protocol | Purpose |
|---------|------|----------|---------|
| Frontend | 5173 | HTTP | React Dev Server |
| Backend | 8080 | HTTP/WS | API & WebSocket |
| Ollama | 11434 | HTTP | AI Model API |

## File System Layout

```
jarvis-ai/
│
├── backend/
│   ├── main.go                    # Entry point
│   ├── router.go                  # HTTP/WS routes
│   ├── jarvis.db                  # SQLite database
│   │
│   └── modules/
│       ├── ai/
│       │   ├── ollama.go          # AI integration
│       │   └── memory.go          # Conversation memory
│       │
│       ├── commands/
│       │   ├── command_parser.go  # ⭐ Main command router
│       │   └── command_registry.go
│       │
│       ├── system/
│       │   ├── apps.go            # ⭐ App launcher
│       │   ├── web_actions.go     # ⭐ Web automation
│       │   ├── folders.go         # ⭐ Folder access
│       │   ├── media_control.go   # ⭐ Media keys
│       │   ├── system_control.go  # ⭐ Power/volume
│       │   └── file_search.go
│       │
│       ├── voice/
│       │   ├── tts.go             # Text-to-speech
│       │   ├── whisper.go         # Speech-to-text
│       │   └── wakeword.go        # Wake word detection
│       │
│       ├── database/
│       │   └── sqlite.go          # Database operations
│       │
│       └── plugins/
│           └── plugin_loader.go   # Plugin system
│
└── frontend/
    └── src/
        ├── components/
        │   ├── JarvisFace.jsx     # 3D animated face
        │   ├── MicIndicator.jsx   # Voice input UI
        │   └── CommandConsole.jsx # Command history
        │
        ├── hooks/
        │   ├── useJarvisSocket.js # WebSocket hook
        │   └── useVoice.js        # Voice input hook
        │
        └── pages/
            └── Dashboard.jsx      # Main UI
```

---

**Legend:**
- ⭐ = Core command execution modules
- → = Data flow direction
- │ = Dependency relationship
