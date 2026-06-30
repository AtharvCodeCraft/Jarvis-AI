from typing import Optional
from datetime import datetime
from pydantic import BaseModel, EmailStr

class UserBase(BaseModel):
    username: str
    email: EmailStr
    full_name: Optional[str] = None
    profile_image: Optional[str] = None
    preferred_model: Optional[str] = "default"
    preferred_voice: Optional[str] = "default"
    theme: Optional[str] = "dark"

class UserCreate(UserBase):
    password: str

class UserUpdate(BaseModel):
    full_name: Optional[str] = None
    username: Optional[str] = None
    profile_image: Optional[str] = None
    preferred_model: Optional[str] = None
    preferred_voice: Optional[str] = None
    theme: Optional[str] = None

class UserInDB(UserBase):
    id: str
    is_verified: bool
    role: str
    created_at: datetime
    updated_at: datetime
    last_login: Optional[datetime] = None

    class Config:
        from_attributes = True

class Token(BaseModel):
    access_token: str
    token_type: str
    user: UserInDB
