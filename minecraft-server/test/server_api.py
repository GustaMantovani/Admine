from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import Dict

app = FastAPI()


class CommandRequest(BaseModel):
    command: str


@app.get("/info")
async def get_info():
    try:
        #simulando uma exceção
        #x = 1 / 0
        
        return {
            "payload": {
                "minecraftVersion": "1.20.4",
                "javaVersion": "17",
                "modEngine": "Fabric",
                "maxPlayers": 20,
                "seed": "123456789"
            }
        }
    except Exception:
        raise HTTPException(status_code=500, detail={"message": "error"})


@app.post("/command")
async def post_command(request: CommandRequest):
    try:
        #simulando uma exceção
        #x = 1 / 0

        # Você pode adicionar lógica aqui para lidar com o comando
        return {"payload": {"result": f"executed: {request.command}"}}
    except Exception:
        raise HTTPException(status_code=500, detail={"message": "error"})


@app.get("/status")
async def get_status():
    try:
        #simulando uma exceção
        #x = 1 / 0
        
        return {
            "payload": {
                "health": "healthy",
                "status": "online",
                "description": "Servidor rodando normalmente",
                "uptime": "2h 15m",
                "onlinePlayers": 5,
                "tps": 19.98
            }
        }
    except Exception:
        raise HTTPException(status_code=500, detail={"message": "error"})


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8080)
