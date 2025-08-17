import redis
import json

def main():
    # Conecta ao Redis
    client = redis.StrictRedis(host='localhost', port=6379, db=0)

    channel = input("Digite o nome do canal: ")
    origin = input("Digite a origem: ")

    while True:
        tags_input = input("Digite as tags (separadas por vírgula): ")
        tags = [tag.strip() for tag in tags_input.split(',')]
        message = input("Digite a mensagem: ")

        # Cria a mensagem no formato JSON
        admine_message = {
            "origin": origin,
            "tags": tags,
            "message": message
        }

        # Converte a mensagem para JSON
        admine_message_json = json.dumps(admine_message)

        # Publica a mensagem no canal
        client.publish(channel, admine_message_json)
        print(f"Mensagem enviada para o {channel}")

if __name__ == "__main__":
    main()
