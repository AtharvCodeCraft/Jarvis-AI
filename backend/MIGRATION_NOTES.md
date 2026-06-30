# Migration Notes: Go to Python

## 1. `main.go` -> `app/main.py`
### Original Go
```go
func main() {
    log.Println("Initializing JARVIS Core Engine...")
    database.InitDB()
    plugins.LoadAll()
    r := SetupRouter()
    r.Run(":8080")
}
```

### Equivalent Python
```python
@app.on_event("startup")
async def startup_event():
    logger.info("Initializing JARVIS Core Engine...")
    load_all_plugins()
    logger.info("JARVIS system online. Listening on port 8080...")

if __name__ == "__main__":
    uvicorn.run("app.main:app", host="0.0.0.0", port=8080, reload=True)
```
**Explanation**: Replaced standard Go Gin bootstrapping with FastAPI's `startup` event and `uvicorn.run()`. Database schemas are provisioned automatically via SQLAlchemy metadata.

---

## 2. `router.go` -> `app/api/routes.py`
### Original Go
```go
func handleCommand(c *gin.Context) { ... }
func handleWebSocket(c *gin.Context) { ... }
func handleSystemStatus(c *gin.Context) { ... }
```

### Equivalent Python
```python
@router.post("/command", response_model=CommandResponse)
async def handle_command(req: CommandRequest): ...

@router.websocket("/ws")
async def websocket_endpoint(websocket: WebSocket): ...

@router.get("/system/status", response_model=SystemStatusResponse)
async def handle_system_status(): ...
```
**Explanation**: Gin's imperative `c.BindJSON` was replaced with FastAPI's automatic request validation using Pydantic schemas. WebSockets use FastAPI's `WebSocket` class instead of Gorilla WebSocket.

---

## 3. `database/sqlite.go` -> `app/models/models.py` & `app/core/database.py`
### Original Go
Raw SQL schemas embedded in `database.DB.Exec(...)`. Example:
```sql
CREATE TABLE IF NOT EXISTS messages (
    id TEXT PRIMARY KEY,
    conversation_id TEXT REFERENCES conversations(id) ON DELETE CASCADE, ...
);
```

### Equivalent Python
```python
class Message(Base):
    __tablename__ = "messages"
    id = Column(String, primary_key=True, index=True)
    conversation_id = Column(String, ForeignKey("conversations.id", ondelete="CASCADE"))
    ...
```
**Explanation**: Adopted SQLAlchemy ORM over raw string interpolations, enhancing maintainability and providing Python object mapping for the database entries. 

---

## 4. `ai/memory.go` -> `app/services/memory_service.py`
### Original Go
```go
func GetRelevantMemories(prompt string) string { ... }
func SaveInteraction(userMsg, aiMsg string) { ... }
```

### Equivalent Python
```python
def get_relevant_memories(prompt: str) -> str: ...
def save_interaction(user_msg: str, ai_msg: str) -> None: ...
```
**Explanation**: The direct SQL queries (`DB.QueryRow`) were replaced with SQLAlchemy `db.query()`. Logic mapping `role` and extracting semantics like "i prefer" was fully retained.

---

## 5. `ai/ollama.go` -> `app/services/ollama_service.py`
### Original Go
```go
func queryOllamaChat(body []byte) ([]byte, error) {
    resp, err := http.Post("http://127.0.0.1:11434/api/chat", ...)
}
```

### Equivalent Python
```python
async def query_ollama_chat(body: Dict[str, Any]) -> Dict[str, Any]:
    async with httpx.AsyncClient() as client:
        resp = await client.post("http://127.0.0.1:11434/api/chat", json=body)
```
**Explanation**: Used asynchronous `httpx` to replace Go's blocking `net/http` client.

---

## 6. `agent/orchestrator.go` -> `app/services/agent_orchestrator.py`
### Original Go
Registered system agents like `SystemAgent`, `BrowserAgent` via `GlobalRouter.Register()`. Formulated LLM payload and mapped function arguments directly.

### Equivalent Python
Created `Agent` and `Router` classes. Python's orchestrator formats standard `ollama_tools` dict lists out of Pydantic-based `ToolDefinition`.

---

## 7. `tool/system_tools.go` -> `app/services/plugin_service.py`
### Original Go
```go
GlobalRegistry.Register("open_application", ..., func(args map[string]interface{}) (string, error) {
    return system.OpenApp(args["app_name"].(string)), nil
})
```

### Equivalent Python
```python
global_registry.register(
    ToolDefinition(name="open_application", ...),
    lambda app_name: system_controls.open_application(app_name)
)
```
**Explanation**: Converted statically typed Go tool schemas to `ToolDefinition` objects leveraging Pydantic, ensuring better static validation. Callbacks use python lambda.

---

## 8. `system/system_control.go` & `apps.go` -> `app/services/system_controls.py`
### Original Go
```go
cmd = exec.Command("powershell", "-Command", "(New-Object -ComObject WScript.Shell).SendKeys([char]175)")
```

### Equivalent Python
```python
subprocess.run(["powershell", "-Command", "(New-Object -ComObject WScript.Shell).SendKeys([char]175)"], creationflags=subprocess.CREATE_NO_WINDOW)
```
**Explanation**: Python's `subprocess.run` directly replaces Go's `os/exec`. Added `creationflags=subprocess.CREATE_NO_WINDOW` to replicate Go's `HideWindow: true`. Web actions like WhatsApp opening were translated perfectly using Python's `webbrowser` library.
