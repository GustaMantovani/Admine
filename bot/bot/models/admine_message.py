import json

class AdmineMessage:
    def __init__(self, tags: list[str], message: str):
        self._tags = tags
        self._message = message

    def get_tags(self) -> list[str]:
        return self._tags
    
    def set_tags(self, tags: list[str]):
        self._tags = tags

    def get_message(self) -> str:
        return self._message
    
    def set_message(self, message: str):
        self._message = message

    @classmethod
    def from_json_to_object(cls, json_str):
        data = json.loads(json_str)
        return cls(**data)
    
    def from_object_to_json(self):
        return json.dumps({"tags": self.get_tags(), "message": self.get_message()})