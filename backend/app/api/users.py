from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session
from app.core.database import get_db
from app.models.models import User
from app.schemas.user_schemas import UserInDB, UserUpdate
from app.auth.dependencies import get_current_user

router = APIRouter()

@router.get("/me", response_model=UserInDB)
def read_users_me(current_user: User = Depends(get_current_user)):
    return current_user

@router.put("/update", response_model=UserInDB)
def update_user(
    user_in: UserUpdate, 
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    if user_in.full_name is not None:
        current_user.full_name = user_in.full_name
    if user_in.username is not None:
        # Check uniqueness if username changed
        if user_in.username != current_user.username:
            existing = db.query(User).filter(User.username == user_in.username).first()
            if existing:
                raise HTTPException(status_code=400, detail="Username already taken")
        current_user.username = user_in.username
    if user_in.profile_image is not None:
        current_user.profile_image = user_in.profile_image
    if user_in.preferred_model is not None:
        current_user.preferred_model = user_in.preferred_model
    if user_in.preferred_voice is not None:
        current_user.preferred_voice = user_in.preferred_voice
    if user_in.theme is not None:
        current_user.theme = user_in.theme
        
    db.commit()
    db.refresh(current_user)
    return current_user

@router.delete("/delete")
def delete_user(current_user: User = Depends(get_current_user), db: Session = Depends(get_db)):
    db.delete(current_user)
    db.commit()
    return {"message": "User deleted successfully"}
