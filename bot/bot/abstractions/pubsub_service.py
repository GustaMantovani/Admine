from abc import ABC, abstractmethod

class PubSubService(ABC):
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