import json

class AdmineMessage:
    def __init__(self, tags: list[str], message: str):
        self.tags = tags
        self.message = message

    @classmethod
    def from_json_to_object(cls, json_str):
        data = json.loads(json_str)
        return cls(**data)
    
    def from_object_to_json(self):
        return json.dumps({"tags": self.get_tags(), "message": self.get_message()})