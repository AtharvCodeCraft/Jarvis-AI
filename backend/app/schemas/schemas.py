from pydantic import BaseModel

class CommandRequest(BaseModel):
    text: str

class SystemStatusResponse(BaseModel):
    cpu: float
    ram: float
    status: str

class CommandResponse(BaseModel):
    result: str
