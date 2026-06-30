from typing import Dict, Any, Callable, List, Optional
import inspect
from pydantic import BaseModel

class ToolParameter(BaseModel):
    type: str
    description: str
    enum: Optional[List[str]] = None
    required: bool = False

class ToolDefinition(BaseModel):
    name: str
    description: str
    parameters: Dict[str, ToolParameter]
    requires_hitl: bool = False

class Tool:
    def __init__(self, definition: ToolDefinition, handler: Callable):
        self.definition = definition
        self.handler = handler

class ToolRegistry:
    def __init__(self):
        self._tools: Dict[str, Tool] = {}

    def register(self, definition: ToolDefinition, handler: Callable):
        self._tools[definition.name] = Tool(definition, handler)

    def get(self, name: str) -> Optional[Tool]:
        return self._tools.get(name)

    def list_definitions(self) -> List[ToolDefinition]:
        return [t.definition for t in self._tools.values()]

    async def invoke(self, name: str, args: Dict[str, Any]) -> str:
        tool = self.get(name)
        if not tool:
            raise ValueError(f"Tool {name} not found")
        
        # Determine if handler is async
        if inspect.iscoroutinefunction(tool.handler):
            return await tool.handler(**args)
        else:
            return tool.handler(**args)

global_registry = ToolRegistry()
