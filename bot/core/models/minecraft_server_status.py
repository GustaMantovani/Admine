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
            online_players: Optional[int] = None,
            tps: Optional[float] = None,
    ):
        self.health = health
        self.status = status
        self.description = description
        self.uptime = uptime
        self.online_players = online_players
        self.tps = tps

    @classmethod
    def from_json(cls, json_data: dict):
        health_value = json_data.get("health", "unknown")
        status_value = json_data.get("status", "unknown")
        try:
            health = HealthStatus(health_value.lower())
        except ValueError:
            health = HealthStatus.UNKNOWN
        try:
            status = ServerStatus(status_value.lower())
        except ValueError:
            status = ServerStatus.UNKNOWN

        return cls(
            health=health,
            status=status,
            description=json_data.get("description"),
            uptime=json_data.get("uptime"),
            online_players=json_data.get("onlinePlayers"),
            tps=json_data.get("tps"),
        )

    def to_json(self) -> dict:
        return {
            "health": self.health.value,
            "status": self.status.value,
            "description": self.description,
            "uptime": self.uptime,
            "onlinePlayers": self.online_players,
            "tps": self.tps,
        }
    
    def __str__(self):
        return (
            f"Status do Servidor: {self.status.value}\n"
            f"Saúde: {self.health.value}\n"
            f"Descrição: {self.description}\n"
            f"Uptime: {self.uptime}\n"
            f"Jogadores Online: {self.online_players}\n"
            f"TPS: {self.tps}"
        )