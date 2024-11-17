import os
import redis

def main():
    # Obtém os nomes dos canais das variáveis de ambiente
    channel1 = os.getenv('CHANNEL1', 'default_channel1')
    channel2 = os.getenv('CHANNEL2', 'default_channel2')

    # Conecta ao Redis
    client = redis.StrictRedis(host='localhost', port=6379, db=0)
    pubsub = client.pubsub()
    pubsub.subscribe([channel1, channel2])

    print(f"Inscrito nos canais {channel1} e {channel2}")

    for message in pubsub.listen():
        if message['type'] == 'message':
            print(f"Recebido do {message['channel'].decode()}: {message['data'].decode()}")

if __name__ == "__main__":
    main()