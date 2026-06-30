from sqlalchemy import Column, Integer, String, Text, ForeignKey, DateTime, Boolean
from sqlalchemy.sql import func
from app.core.database import Base

class User(Base):
    __tablename__ = "users"

    id = Column(String, primary_key=True, index=True)
    full_name = Column(String, nullable=True)
    username = Column(String, unique=True, index=True, nullable=False)
    email = Column(String, unique=True, index=True, nullable=False)
    password_hash = Column(String, nullable=False)
    profile_image = Column(String, nullable=True)
    created_at = Column(DateTime, server_default=func.now())
    updated_at = Column(DateTime, server_default=func.now(), onupdate=func.now())
    last_login = Column(DateTime, nullable=True)
    is_verified = Column(Boolean, default=False)
    role = Column(String, default="user")
    preferred_model = Column(String, default="default")
    preferred_voice = Column(String, default="default")
    theme = Column(String, default="dark")

class Memory(Base):
    __tablename__ = "memory"

    id = Column(Integer, primary_key=True, index=True, autoincrement=True)
    key = Column(String, unique=True, nullable=False)
    value = Column(Text, nullable=False)

class Conversation(Base):
    __tablename__ = "conversations"

    id = Column(String, primary_key=True, index=True)
    user_id = Column(String, ForeignKey("users.id", ondelete="CASCADE"), nullable=False)
    title = Column(String, nullable=False)
    created_at = Column(DateTime, server_default=func.now())
    updated_at = Column(DateTime, server_default=func.now(), onupdate=func.now())

class Message(Base):
    __tablename__ = "messages"

    id = Column(String, primary_key=True, index=True)
    conversation_id = Column(String, ForeignKey("conversations.id", ondelete="CASCADE"))
    sender = Column(String) # 'user', 'assistant', 'system'
    content = Column(Text, nullable=False)
    agent_assigned = Column(String, nullable=True)
    tokens_used = Column(Integer, default=0)
    created_at = Column(DateTime, server_default=func.now())

class MemoryNode(Base):
    __tablename__ = "memory_nodes"

    id = Column(String, primary_key=True, index=True)
    user_id = Column(String, nullable=False)
    memory_type = Column(String) # 'semantic', 'preference', 'episodic'
    content = Column(Text, nullable=False)
    last_accessed = Column(DateTime, server_default=func.now())
    created_at = Column(DateTime, server_default=func.now())

class Document(Base):
    __tablename__ = "documents"

    id = Column(String, primary_key=True, index=True)
    filename = Column(String, nullable=False)
    file_type = Column(String, nullable=False)
    file_path = Column(String, nullable=False)
    checksum = Column(String, unique=True, nullable=False)
    chunk_count = Column(Integer, nullable=False)
    created_at = Column(DateTime, server_default=func.now())

class ToolLog(Base):
    __tablename__ = "tool_logs"

    id = Column(String, primary_key=True, index=True)
    agent_name = Column(String, nullable=False)
    tool_name = Column(String, nullable=False)
    arguments = Column(Text, nullable=False)
    execution_status = Column(String) # 'pending', 'approved', 'rejected', 'success', 'failed'
    error_message = Column(Text, nullable=True)
    duration_ms = Column(Integer, nullable=True)
    created_at = Column(DateTime, server_default=func.now())
