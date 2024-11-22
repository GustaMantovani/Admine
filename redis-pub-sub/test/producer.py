import os
import redis
import time

def main():
    # Obtém os nomes dos canais das variáveis de ambiente
    channel1 = os.getenv('CHANNEL1', 'default_channel1')
    channel2 = os.getenv('CHANNEL2', 'default_channel2')

    # Conecta ao Redis
    client = redis.StrictRedis(host='localhost', port=6379, db=0)

    while True:

        msg = input("Digite a mensagem: ")

        client.publish(channel1, f'Mensagem {msg} para {channel1}')
        client.publish(channel2, f'Mensagem {msg} para {channel2}')
        print(f"Mensagens enviadas para {channel1} e {channel2}")

if __name__ == "__main__":
    main()