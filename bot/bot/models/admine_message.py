import json

class AdmineMessage:
    def __init__(self, tags: list[str], message: str):
        self._tags = tags
        self._message = message

    def getTags(self) -> list[str]:
        return self._tags
    
    def setTags(self,tags: list[str]):
        self._tags = tags

    def getMessage(self) -> str:
        return self._message
    
    def setMessage(self,message: str):
        self._message = message

    @classmethod
    def from_json_to_object(cls, json_str):
        data = json.loads(json_str)  # Converte JSON para dicionário
        return cls(**data)  # Usa os dados para criar um objeto
    
 
    def from_objetc_to_json(self):
        return json.dumps({"tags": self.getTags(), "message": self.getMessage()})