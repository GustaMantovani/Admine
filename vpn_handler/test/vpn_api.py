from fastapi import FastAPI, HTTPException
from pydantic import BaseModel

app = FastAPI()


class AuthRequest(BaseModel):
    member_id: str


@app.get("/server-ip")
async def get_server_ip():
    try:
        return {"payload":{"server_ip": "123.123.123.123"}}
    except Exception:
        raise HTTPException(status_code=500, detail={"message": "error"})


@app.post("/auth-member")
async def auth_member(request: AuthRequest):
    try:
        if request.member_id != "123456":
            return {"payload":{"message": "member not found"}}
        return {"payload":{"message": "authenticated"}}
    except Exception:
        raise HTTPException(status_code=500, detail={"message": "error"})


@app.get("/vpn-id")
async def get_vpn_id():
    try:
        return {"payload":{"vpn_id": "af;dfj;lkj1jk5132l4k5"}}
    except Exception:
        raise HTTPException(status_code=500, detail={"message": "error"})





if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="127.0.0.1", port=9090)