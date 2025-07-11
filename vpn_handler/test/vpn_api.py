from fastapi import FastAPI, HTTPException
from pydantic import BaseModel

app = FastAPI()

class AuthRequest(BaseModel):
    member_id: str


@app.get("/server-ip")
async def get_server_ip():
    try:
        return {"server_ip": "123.123.123.123"}
    except Exception:
        raise HTTPException(status_code=500, detail={"message": "error"})


@app.post("/auth-member")
async def auth_member(request: AuthRequest):
    try:
        if request.member_id != "valid_member":
            return {"message": "member not found"}
        return {"message": "authenticated"}
    except Exception:
        raise HTTPException(status_code=500, detail={"message": "error"})


@app.get("/vpn-id")
async def get_vpn_id():
    try:
        return {"vpn_id": "af;dfj;lkj1jk5132l4k5"}
    except Exception:
        raise HTTPException(status_code=500, detail={"message": "error"})


if name == "main":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=9090)