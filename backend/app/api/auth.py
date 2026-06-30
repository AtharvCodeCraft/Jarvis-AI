import uuid
from datetime import datetime
from fastapi import APIRouter, Depends, HTTPException, status
from fastapi.security import OAuth2PasswordRequestForm
from sqlalchemy.orm import Session
from app.core.database import get_db
from app.models.models import User
from app.schemas.user_schemas import UserCreate, Token, UserInDB
from app.auth.utils import get_password_hash, verify_password, create_access_token
from pydantic import BaseModel

router = APIRouter()

class RegisterResponse(BaseModel):
    message: str
    user: UserInDB

class ForgotPasswordReq(BaseModel):
    email: str

@router.post("/register", response_model=RegisterResponse)
def register(user_in: UserCreate, db: Session = Depends(get_db)):
    # Check if email exists
    user = db.query(User).filter(User.email == user_in.email).first()
    if user:
        raise HTTPException(
            status_code=400,
            detail="A user with this email already exists."
        )
    
    # Check if username exists
    user_by_username = db.query(User).filter(User.username == user_in.username).first()
    if user_by_username:
        raise HTTPException(
            status_code=400,
            detail="A user with this username already exists."
        )

    # Password validation based on requirements
    import re
    if len(user_in.password) < 8 or not re.search(r"[A-Z]", user_in.password) or \
       not re.search(r"[a-z]", user_in.password) or not re.search(r"[0-9]", user_in.password) or \
       not re.search(r"[!@#$%^&*(),.?\":{}|<>]", user_in.password):
        raise HTTPException(
            status_code=400,
            detail="Password must be at least 8 characters, and contain an uppercase letter, lowercase letter, number, and special character."
        )

    user_id = str(uuid.uuid4())
    hashed_password = get_password_hash(user_in.password)
    
    new_user = User(
        id=user_id,
        email=user_in.email,
        username=user_in.username,
        full_name=user_in.full_name,
        password_hash=hashed_password,
        role="user",
        is_verified=False
    )
    db.add(new_user)
    db.commit()
    db.refresh(new_user)
    
    return {"message": "User registered successfully", "user": new_user}

@router.post("/login", response_model=Token)
def login(form_data: OAuth2PasswordRequestForm = Depends(), db: Session = Depends(get_db)):
    # Authenticate
    user = db.query(User).filter((User.email == form_data.username) | (User.username == form_data.username)).first()
    if not user or not verify_password(form_data.password, user.password_hash):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Incorrect email/username or password",
            headers={"WWW-Authenticate": "Bearer"},
        )
    
    user.last_login = datetime.utcnow()
    db.commit()
    db.refresh(user)

    access_token = create_access_token(subject=user.id)
    return {"access_token": access_token, "token_type": "bearer", "user": user}

@router.post("/logout")
def logout():
    # In a real JWT stateless setup, logout is handled client side by removing the token.
    # Optionally, we could blacklist tokens in a DB/Redis.
    return {"message": "Successfully logged out"}

@router.post("/refresh")
def refresh():
    # Placeholder for refresh token logic. Usually handled similarly to login, but using a refresh_token
    pass

@router.post("/forgot-password")
def forgot_password(req: ForgotPasswordReq, db: Session = Depends(get_db)):
    # Placeholder for sending reset link
    return {"message": "If this email is registered, a password reset link will be sent."}

@router.post("/reset-password")
def reset_password():
    # Placeholder
    pass
