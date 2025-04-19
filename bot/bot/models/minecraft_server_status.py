from enum import Enum
from typing import Optional

class HealthStatus(Enum):
    HEALTHY = "healthy"
    SICK = "sick"
    CRITICAL = "critical"
    UNKNOWN = "unknown"

class ServerStatus(Enum):
    ONLINE = "online"
    OFFLINE = "offline"
    MAINTENANCE = "maintenance"
    UNKNOWN = "unknown"

class MinecraftServerStatus:
    def __init__(
        self,
        health: HealthStatus,
        status: ServerStatus,
        description: str,
        uptime: Optional[str] = None,
        online_players: Optional[int] = 0,
        tps: Optional[float] = None
    ):
        self.health = health
        self.status = status
        self.description = description
        self.uptime = uptime
        self.online_players = online_players
        self.tps = tps

    @classmethod
    def from_json(cls, json_data: dict):
        return cls(
            health=HealthStatus(json_data.get("health").lower()),
            status=ServerStatus(json_data.get("status").lower()),
            description=json_data.get("description"),
            uptime=json_data.get("uptime"),
            online_players=json_data.get("onlinePlayers"),
            tps=json_data.get("tps")
        )

    def to_json(self) -> dict:
        return {
            "health": self.health.value,
            "status": self.status.value,
            "description": self.description,
            "uptime": self.uptime,
            "onlinePlayers": self.online_players,
            "tps": self.tps
        }