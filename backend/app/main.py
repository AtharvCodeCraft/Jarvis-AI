import logging
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from app.api.routes import router
from app.api.auth import router as auth_router
from app.api.users import router as users_router
from app.core.database import engine, Base
from app.services.plugin_service import load_all_plugins

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Initialize database
Base.metadata.create_all(bind=engine)

app = FastAPI(title="JARVIS Core Engine")

# CORS for frontend
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Include API router
app.include_router(router)
app.include_router(auth_router, prefix="/api/auth", tags=["auth"])
app.include_router(users_router, prefix="/api/users", tags=["users"])

@app.on_event("startup")
async def startup_event():
    logger.info("Initializing JARVIS Core Engine...")
    # Load Plugins
    load_all_plugins()
    logger.info("JARVIS system online. Listening on port 8080...")

if __name__ == "__main__":
    import uvicorn
    uvicorn.run("app.main:app", host="0.0.0.0", port=8080, reload=True)
