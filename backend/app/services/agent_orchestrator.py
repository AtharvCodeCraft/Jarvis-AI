import logging
import json
import time
import uuid as uuid_pkg
from typing import Dict, Any, List
from app.services.memory_service import get_relevant_memories, get_or_create_conversation, get_conversation_history, save_interaction, get_context
from app.services.tool_registry import global_registry
from app.services.ollama_service import get_available_model, query_ollama_chat
from app.core.database import SessionLocal
from app.models.models import ToolLog

logger = logging.getLogger(__name__)

class Agent:
    def __init__(self, name: str, description: str, system_prompt: str, allowed_tools: List[str]):
        self.name = name
        self.description = description
        self.system_prompt = system_prompt
        self.allowed_tools = allowed_tools

class Router:
    def __init__(self):
        self.agents: Dict[str, Agent] = {}

    def register(self, agent: Agent):
        self.agents[agent.name] = agent
        logger.info(f"Agent registered: {agent.name}")

    def route(self, prompt: str) -> Agent:
        prompt_lower = prompt.lower()
        
        coding_keywords = ["code", "program", "function", "debug", "write a script", "compile", "refactor", "bug", "scaffold", "golang", "javascript", "react", "html", "python"]
        for kw in coding_keywords:
            if kw in prompt_lower:
                return self.agents["CodingAgent"]
        
        browser_keywords = ["search", "google", "youtube", "website", "whatsapp", "call rahul", "chat with", "link", "browse", "internet", "open chrome"]
        for kw in browser_keywords:
            if kw in prompt_lower:
                return self.agents["BrowserAgent"]
        
        return self.agents["SystemAgent"]

global_router = Router()

# Register Agents
global_router.register(Agent(
    name="SystemAgent",
    description="Handles OS tasks, starting desktop applications, folder navigation, volume adjustments, and system power commands.",
    system_prompt="""You are the System Agent of the Jarvis platform. Your primary purpose is to help the user manage their computer system. 
You have direct access to tools for launching apps, navigating local user directories, adjusting sound volume, and managing power status.
Be concise and focus on executing the exact system request.""",
    allowed_tools=["open_application", "open_folder", "control_volume", "power_control", "system_status"]
))

global_router.register(Agent(
    name="BrowserAgent",
    description="Handles web navigation, opening websites, performing Google/YouTube searches, and WhatsApp automation.",
    system_prompt="""You are the BrowserAgent. Your role is to assist with web navigation, online search queries, YouTube video streaming lookups, and opening WhatsApp.
Formulate clear search intents. Use Google Search and Web navigation tools dynamically.""",
    allowed_tools=["google_search", "youtube_search", "open_website", "whatsapp_action"]
))

global_router.register(Agent(
    name="CodingAgent",
    description="Specialized in code generation, code review, debugging, project scaffolding, and architecture suggestions.",
    system_prompt="""You are the Coding Agent. Your task is to review, write, edit, and explain software source code.
Maintain high quality, use best practices, comment complex logic, and help scaffold modular structures.""",
    allowed_tools=[]
))

def log_tool_execution(agent_name: str, tool_name: str, args: dict, status: str, err: str, duration_ms: int):
    with SessionLocal() as db:
        log_entry = ToolLog(
            id=str(uuid_pkg.uuid4()),
            agent_name=agent_name,
            tool_name=tool_name,
            arguments=json.dumps(args),
            execution_status=status,
            error_message=err,
            duration_ms=duration_ms
        )
        db.add(log_entry)
        db.commit()

async def execute_agent(prompt: str) -> str:
    agent = global_router.route(prompt)
    logger.info(f"[AgentOrchestrator] Selected Specialized Agent: {agent.name}")

    mem_context = get_relevant_memories(prompt)
    full_system_prompt = agent.system_prompt
    if mem_context:
        full_system_prompt += "\n" + mem_context
    
    # Also inject base context which has language rules
    full_system_prompt += "\n" + get_context(prompt)

    conv_id = get_or_create_conversation()
    history = get_conversation_history(conv_id, 8)

    messages = [{"role": "system", "content": full_system_prompt}]
    messages.extend(history)
    messages.append({"role": "user", "content": prompt})

    ollama_tools = []
    for tool_name in agent.allowed_tools:
        rt = global_registry.get(tool_name)
        if not rt:
            continue
        
        props = {}
        required = []
        for p_name, p_def in rt.definition.parameters.items():
            param_map = {
                "type": p_def.type,
                "description": p_def.description,
            }
            if p_def.enum:
                param_map["enum"] = p_def.enum
            props[p_name] = param_map
            if p_def.required:
                required.append(p_name)
        
        ollama_tools.append({
            "type": "function",
            "function": {
                "name": rt.definition.name,
                "description": rt.definition.description,
                "parameters": {
                    "type": "object",
                    "properties": props,
                    "required": required
                }
            }
        })

    model = await get_available_model()
    req_body = {
        "model": model,
        "messages": messages,
        "stream": False,
    }
    if ollama_tools:
        req_body["tools"] = ollama_tools

    chat_resp = await query_ollama_chat(req_body)

    if chat_resp.get("message", {}).get("tool_calls"):
        tool_call = chat_resp["message"]["tool_calls"][0]
        tool_name = tool_call["function"]["name"]
        tool_args_map = tool_call["function"]["arguments"]
        
        logger.info(f"[{agent.name}] Invoking tool: {tool_name} with args: {tool_args_map}")

        t_obj = global_registry.get(tool_name)
        if t_obj and t_obj.definition.requires_hitl:
            log_tool_execution(agent.name, tool_name, tool_args_map, "pending", None, 0)
            return f"CONFIRMATION_REQUIRED:{tool_name}:{json.dumps(tool_args_map)}"

        start_time = time.time()
        try:
            result = await global_registry.invoke(tool_name, tool_args_map)
            duration = int((time.time() - start_time) * 1000)
            log_tool_execution(agent.name, tool_name, tool_args_map, "success", None, duration)
            save_interaction(prompt, result)
            return result
        except Exception as e:
            duration = int((time.time() - start_time) * 1000)
            log_tool_execution(agent.name, tool_name, tool_args_map, "failed", str(e), duration)
            return f"Error executing tool {tool_name}: {e}"
    
    reply = chat_resp.get("message", {}).get("content", "")
    save_interaction(prompt, reply)
    return reply
