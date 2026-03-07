class ResourceUsage:
    def __init__(
        self,
        cpu_usage: float,
        memory_used: int,
        memory_total: int,
        memory_used_percent: float,
        disk_used: int,
        disk_total: int,
        disk_used_percent: float,
    ):
        self.cpu_usage = cpu_usage
        self.memory_used = memory_used
        self.memory_total = memory_total
        self.memory_used_percent = memory_used_percent
        self.disk_used = disk_used
        self.disk_total = disk_total
        self.disk_used_percent = disk_used_percent

    @classmethod
    def from_json(cls, json_data: dict) -> "ResourceUsage":
        return cls(
            cpu_usage=float(json_data.get("cpu_usage", 0)),
            memory_used=int(json_data.get("memory_used", 0)),
            memory_total=int(json_data.get("memory_total", 0)),
            memory_used_percent=float(json_data.get("memory_used_percent", 0)),
            disk_used=int(json_data.get("disk_used", 0)),
            disk_total=int(json_data.get("disk_total", 0)),
            disk_used_percent=float(json_data.get("disk_used_percent", 0)),
        )
