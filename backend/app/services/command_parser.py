import logging
from app.services.agent_orchestrator import execute_agent

logger = logging.getLogger(__name__)

def normalize_input(original: str) -> str:
    original = original.lower().strip()
    wake_words = [
        "jarvis ", "hey jarvis ", "ok jarvis ", "अरे जार्विस ", "जार्विस ", "hey rudra ", "rudra "
    ]
    for ww in wake_words:
        if original.startswith(ww):
            original = original[len(ww):]
    return original.strip()

async def parse_and_execute(input_text: str) -> str:
    normalized = normalize_input(input_text)
    logger.info(f"[CommandRouter] Normalized request: {normalized}")
    
    try:
        response = await execute_agent(normalized)
        return response
    except Exception as e:
        logger.error(f"[CommandRouter] Error invoking Agent Orchestrator: {e}")
        return "I encountered an error executing that request. Please verify that Ollama is running locally."
