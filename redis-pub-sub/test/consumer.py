import os
import redis

def main():

    channel = input("Digite o nome do canal: ")

    # Conecta ao Redis
    client = redis.StrictRedis(host='localhost', port=6379, db=0)
    pubsub = client.pubsub()
    pubsub.subscribe([channel])

    print(f"Inscrito no canl {channel}")

    for message in pubsub.listen():
        if message['type'] == 'message':
            print(f"Recebido do {message['channel'].decode()}: {message['data'].decode()}")

if __name__ == "__main__":
    main()