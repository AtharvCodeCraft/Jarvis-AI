import httpx
import logging
import json
from typing import Dict, Any, List

logger = logging.getLogger(__name__)

async def get_available_model() -> str:
    try:
        async with httpx.AsyncClient() as client:
            resp = await client.get("http://127.0.0.1:11434/api/tags", timeout=5.0)
            if resp.status_code == 200:
                data = resp.json()
                models = [m.get("name") for m in data.get("models", [])]
                if models:
                    order = ["qwen2.5:1.5b", "qwen2.5:0.5b", "qwen2.5-coder:7b", "llama3.2:1b"]
                    for preferred in order:
                        for m in models:
                            if m == preferred or m.startswith(preferred):
                                return m
                    return models[0]
    except Exception as e:
        logger.warning(f"Error fetching Ollama models: {e}")
    return "qwen2.5:1.5b"

async def query_ollama_chat(body: Dict[str, Any]) -> Dict[str, Any]:
    async with httpx.AsyncClient() as client:
        resp = await client.post(
            "http://127.0.0.1:11434/api/chat",
            json=body,
            timeout=30.0
        )
        if resp.status_code != 200:
            raise Exception(f"Ollama API status {resp.status_code}: {resp.text}")
        return resp.json()
