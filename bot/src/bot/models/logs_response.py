class LogsResponse:
    def __init__(self, lines: list[str], total: int):
        self.lines = lines
        self.total = total

    @classmethod
    def from_json(cls, json_data: dict) -> "LogsResponse":
        lines = json_data.get("lines", [])
        return cls(lines=lines, total=int(json_data.get("total", len(lines))))
