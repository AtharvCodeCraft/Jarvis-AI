import logging
import psutil
from fastapi import APIRouter, WebSocket, WebSocketDisconnect, Depends
from app.schemas.schemas import CommandRequest, SystemStatusResponse, CommandResponse
from app.services.command_parser import parse_and_execute
import json
from app.core.context import current_user_id
from app.auth.dependencies import get_current_user
from app.models.models import User
from jose import jwt, JWTError
from app.core.config import settings

router = APIRouter()
logger = logging.getLogger(__name__)

@router.post("/command", response_model=CommandResponse)
async def handle_command(req: CommandRequest, current_user: User = Depends(get_current_user)):
    current_user_id.set(current_user.id)
    result = await parse_and_execute(req.text)
    return CommandResponse(result=result)

@router.get("/system/status", response_model=SystemStatusResponse)
async def handle_system_status():
    # Simulated system stats, or real using psutil
    cpu = psutil.cpu_percent(interval=None)
    ram = psutil.virtual_memory().percent
    return SystemStatusResponse(cpu=cpu, ram=ram, status="Online")

@router.websocket("/ws")
async def websocket_endpoint(websocket: WebSocket, token: str = None):
    await websocket.accept()
    logger.info("Client connected to WS")
    
    user_id = "default_user"
    if token:
        try:
            payload = jwt.decode(token, settings.SECRET_KEY, algorithms=[settings.ALGORITHM])
            uid = payload.get("sub")
            if uid:
                user_id = uid
        except JWTError:
            pass
            
    current_user_id.set(user_id)
    
    try:
        while True:
            data = await websocket.receive_text()
            try:
                msg = json.loads(data)
                event_type = msg.get("type")
                if event_type in ("voice_command", "manual_command"):
                    text = msg.get("text", "")
                    language = msg.get("language", "en-IN")
                    logger.info(f"Received command via WS [{language}]: {text}")
                    
                    # Process command
                    response = await parse_and_execute(text)
                    
                    await websocket.send_json({
                        "type": "command_result",
                        "text": response,
                        "language": language
                    })
            except Exception as e:
                logger.error(f"WS Message Processing Error: {e}")
    except WebSocketDisconnect:
        logger.info("Client disconnected from WS")
