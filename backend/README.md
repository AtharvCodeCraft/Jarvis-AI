# Jarvis backend migration to Python

## Overview
This repository contains the Python migration of the Jarvis AI backend, originally written in Go (Gin framework). The backend powers the entire AI assistant ecosystem, providing REST APIs, WebSocket support for voice interaction, OS-level integration tools, and an advanced Agent Orchestrator to route user intents to specialized sub-agents.

## Technology Stack
*   **Framework**: FastAPI
*   **Server**: Uvicorn
*   **Database**: SQLite with SQLAlchemy ORM
*   **LLM Inference**: Local Ollama (e.g. qwen2.5:1.5b)
*   **System Automation**: pywin32, pyautogui, subprocess
*   **Typing/Validation**: Pydantic

## Setup Instructions

### 1. Prerequisites
- Python 3.10+
- Install Ollama locally and pull your preferred model (e.g. `ollama run qwen2.5:1.5b`).

### 2. Installation
```bash
# Clone or navigate to the project directory
cd jarvis-backend

# Create a virtual environment
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate

# Install dependencies
pip install -r requirements.txt
```

### 3. Environment Configuration
Copy the `.env.example` file to `.env`:
```bash
copy .env.example .env
```
Ensure `OLLAMA_API_URL` points to your local instance (default `http://127.0.0.1:11434/api`).

### 4. Running the Server
```bash
# Start the FastAPI server using Uvicorn
uvicorn app.main:app --host 0.0.0.0 --port 8080 --reload
```
The API documentation will be available at: http://localhost:8080/docs

## API Documentation

- `POST /command`: Takes JSON `{"text": "your command"}` and returns `{"result": "response text"}`.
- `GET /system/status`: Returns JSON with `cpu`, `ram`, and `status`.
- `WebSocket /ws`: Connects to stream bidirectional events (e.g. `voice_command`, `manual_command`). Sends back `{"type": "command_result", "text": "response", "language": "en-IN"}`.

## Migration Notes
- **Goroutines to Asyncio**: The system now extensively leverages FastAPI's async capabilities with `async/await` syntax instead of goroutines.
- **Gin to FastAPI**: Replaced Gin router and JSON bindings with FastAPI routing and Pydantic schema validations.
- **Go HTTP client to httpx**: Inter-process communication to Ollama is handled through `httpx.AsyncClient` for unblocked, async network requests.
- **Agent Orchestrator**: The semantic agent router logic was ported directly to `app/services/agent_orchestrator.py`, maintaining compatibility with `ToolRegistry`.
- **System Interactions**: OS interactions previously using Go's `os/exec` to invoke `powershell` or `rundll32` have been replicated accurately using Python's `subprocess` and `webbrowser` modules.
- **Database**: Ported from raw `modernc.org/sqlite` to `SQLAlchemy` using `declarative_base`. Schemas precisely match legacy SQLite tables.
