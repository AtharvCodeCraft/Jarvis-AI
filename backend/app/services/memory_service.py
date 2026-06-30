import logging
import uuid
import uuid as uuid_pkg
from typing import List, Dict
from sqlalchemy.orm import Session
from sqlalchemy import or_, desc
from app.core.database import SessionLocal
from app.models.models import Conversation, Message, MemoryNode
from app.core.context import current_user_id

logger = logging.getLogger(__name__)

def get_db_session() -> Session:
    return SessionLocal()

def get_or_create_conversation() -> str:
    user_id = current_user_id.get()
    with get_db_session() as db:
        conv = db.query(Conversation).filter(Conversation.user_id == user_id).order_by(desc(Conversation.updated_at)).first()
        if not conv:
            conv_id = str(uuid_pkg.uuid4())
            conv = Conversation(id=conv_id, user_id=user_id, title="Jarvis Interaction Session")
            db.add(conv)
            db.commit()
        else:
            conv_id = conv.id
        return conv_id

def save_message(conv_id: str, sender: str, content: str, agent_assigned: str = None) -> None:
    with get_db_session() as db:
        msg = Message(
            id=str(uuid_pkg.uuid4()),
            conversation_id=conv_id,
            sender=sender,
            content=content,
            agent_assigned=agent_assigned
        )
        db.add(msg)
        db.commit()

def get_conversation_history(conv_id: str, limit: int = 8) -> List[Dict[str, str]]:
    with get_db_session() as db:
        msgs = db.query(Message).filter(Message.conversation_id == conv_id).order_by(Message.created_at).limit(limit).all()
        history = []
        for msg in msgs:
            role = msg.sender if msg.sender in ["user", "assistant", "system"] else ("assistant" if msg.sender != "user" else "user")
            history.append({"role": role, "content": msg.content})
        return history

def save_semantic_memory(content: str, memory_type: str) -> None:
    user_id = current_user_id.get()
    with get_db_session() as db:
        node = MemoryNode(
            id=str(uuid_pkg.uuid4()),
            user_id=user_id,
            memory_type=memory_type,
            content=content
        )
        db.add(node)
        db.commit()

def get_relevant_memories(prompt: str) -> str:
    words = [w.strip() for w in prompt.lower().split() if w.strip()]
    if not words:
        return ""
    
    user_id = current_user_id.get()
    with get_db_session() as db:
        filters = [MemoryNode.content.ilike(f"%{w}%") for w in words]
        nodes = db.query(MemoryNode).filter(MemoryNode.user_id == user_id).filter(or_(*filters)).limit(5).all()
        
        memories = [node.content for node in nodes]
        if memories:
            return "\n[Retrieved User Memories]:\n- " + "\n- ".join(memories)
        return ""

def save_interaction(user_msg: str, ai_msg: str) -> None:
    conv_id = get_or_create_conversation()
    save_message(conv_id, "user", user_msg)
    save_message(conv_id, "assistant", ai_msg, "Orchestrator")
    
    lower_msg = user_msg.lower()
    if "i prefer" in lower_msg or "my name is" in lower_msg or "remember that" in lower_msg:
        save_semantic_memory(user_msg, "preference")
        logger.info("New preference saved to semantic memory node.")

def get_context(prompt: str) -> str:
    mem_context = get_relevant_memories(prompt)
    base_context = """You are Jarvis (also called Rudra), a highly advanced local AI assistant managing this system. Keep your responses concise.

IMPORTANT — Language Rules:
- Detect the language of the user's message automatically.
- If the user writes or speaks in Hindi (हिंदी), you MUST reply in Hindi.
- If the user writes or speaks in Marathi (मराठी), you MUST reply in Marathi.
- If the user writes or speaks in English, reply in English.
- Never switch languages on your own. Always match the user's language exactly.
- You are fluent in English, Hindi, and Marathi."""
    if mem_context:
        base_context += "\n" + mem_context
    return base_context
