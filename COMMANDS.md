# JARVIS Voice Commands Guide

## 1. Application Control
Open any installed Windows application:
- `Jarvis open Chrome`
- `Jarvis open VS Code`
- `Jarvis open Notepad`
- `Jarvis open Calculator`
- `Jarvis open File Explorer`
- `Jarvis open Paint`
- `Jarvis open Camera`
- `Jarvis open Settings`
- `Jarvis open Task Manager`

## 2. Google Search
Perform Google searches automatically:
- `Jarvis search Go programming tutorial on Google`
- `Jarvis search weather in Mumbai`
- `Jarvis search best restaurants near me`

## 3. WhatsApp Actions
Open WhatsApp Web and interact:
- `Jarvis open WhatsApp`
- `Jarvis call Rahul on WhatsApp` (opens WhatsApp, you search for contact)
- `Jarvis message John on WhatsApp`

## 4. System Power Commands
Control system power state:
- `Jarvis shutdown the laptop`
- `Jarvis restart the system`
- `Jarvis lock the computer`
- `Jarvis put the system to sleep`

## 5. File and Folder Access
Open common Windows folders:
- `Jarvis open downloads folder`
- `Jarvis open documents`
- `Jarvis open pictures`
- `Jarvis open desktop`
- `Jarvis open music folder`
- `Jarvis open videos`

## 6. Media Control
Control music playback:
- `Jarvis play music`
- `Jarvis pause music`
- `Jarvis next song`
- `Jarvis previous song`
- `Jarvis open music player`

## 7. Volume Control
Adjust system volume:
- `Jarvis volume up`
- `Jarvis volume down`
- `Jarvis mute`
- `Jarvis increase volume`
- `Jarvis decrease volume`

## 8. Camera and Photos
Access camera and pictures:
- `Jarvis open camera`
- `Jarvis show my pictures`
- `Jarvis open photos`

## 9. Web Navigation
Open popular websites:
- `Jarvis open YouTube`
- `Jarvis open Gmail`
- `Jarvis open Facebook.com`
- `Jarvis search Python tutorial on YouTube`

## 10. Smart AI Responses
Ask general questions and Jarvis will respond using local AI:
- `Jarvis what's the weather like?`
- `Jarvis tell me a joke`
- `Jarvis what can you do?`

---

## How to Use

1. Start the backend server: `cd backend && go run main.go router.go`
2. Start the frontend: `cd frontend && npm run dev`
3. Open the web interface and click the microphone to speak
4. Say "Jarvis" followed by your command
5. Jarvis will execute the command and respond with voice feedback

## Technical Architecture

```
Voice Input → Speech-to-Text (Whisper) → Command Parser → 
Intent Detection → Action Executor → Text-to-Speech → Voice Response
```

## Requirements

- Windows OS
- Go 1.26+
- Node.js for frontend
- Ollama (for AI responses)
- Modern web browser with microphone access
