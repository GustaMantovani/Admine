from abc import ABC, abstractmethod
import redis
import json
import discord
from discord import app_commands
#=========================================================
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
        

#=========================================================
class Messager(ABC):
    def __init__(self,canal: str, administradores: list[str]):
        self._canal = canal
        self._administradores = administradores
    
    def getCanal(self)-> str:
        return self._canal

    def setCanal(self,canal: str):
        self._canal = canal
    
    def getAdministradores(self)->list[str]:
        return self._administradores
    
    def setAdministradores(self,administradores: list[str]):
        self._administradores = administradores
    
    @abstractmethod
    def enviarMensagem(self):
        pass

    @abstractmethod  
    def ouvirMensagem(self):
        pass

#=========================================================
class PubSub(ABC):
    def __init__(self, host: str, porta: int, canaisInscrito: list[str], canaisProdutor: list[str]):
        self._host = host
        self._porta = porta
        self._canaisInscrito = canaisInscrito
        self._canaisProdutor = canaisProdutor
    
    def getHost(self)->str:
        return self._host
    
    def setHost(self, host: str):
        self._host = host
    
    def getPorta(self)->int:
        return self._porta
    
    def setPorta(self,porta: int):
        self._porta = porta
    
    def getCanaisInscrito(self)->list[str]:
        return self._canaisInscrito
    
    def getCanaisProdutor(self)->list[str]:
        return self._canaisProdutor
    
    @abstractmethod
    def enviarMensagem(self):
        pass

    @abstractmethod
    def ouvirMensagem(self):
        pass

#=========================================================
class BotDiscord(Messager, discord.Client):
    def __init__(self, canal: str, administradores: list[str]):
        intents = discord.Intents.all()
        Messager.__init__(self, canal, administradores)
        discord.Client.__init__(self,intents=intents)
        self.tree = app_commands.CommandTree(self)

    async def setup_hook(self):
        await self.tree.sync()

    async def on_ready(self):
        print("bot ligado")

    def enviarMensagem(self):
        print("qualquer coisa")

    def ouvirMensagem(self, pubsub):
        return "ouvindo"

#=========================================================

class MeuBot(discord.Client):
    def __init__(self):
        intents = discord.Intents.all()
        super().__init__(command_prefix = "!", intents=intents)
        self.tree = app_commands.CommandTree(self)

    async def setup_hook(self):
        await self.tree.sync()

    async def on_ready(self):
        print("bot ligado")

#=========================================================
class RedisPubSub(PubSub):
    def __init__(self, host, porta, canaisInscrito, canaisProdutor):
        super().__init__(host, porta, canaisInscrito, canaisProdutor)
        cliente = redis.StrictRedis(host, porta, db=0)
        self._cliente = cliente
        self._pubsub = cliente.pubsub()

    def enviarMensagem(self, mensagem: AdmineMessage):
        dados = mensagem.from_objetc_to_json()
        self._cliente.publish("teste", dados)

    def ouvirMensagem(self):
        self._pubsub.subscribe(self.getCanaisInscrito())
        for message in self._pubsub.listen():  # Itera sobre o gerador
            if message["type"] == "message":  # Filtra mensagens reais
                return message