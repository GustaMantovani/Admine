from classes import AdmineMessage, RedisPubSub, BotDiscord, MeuBot
import redis
import threading
import json
import discord
import os

def processarMensagemDiscord(discord: BotDiscord, pubsub: RedisPubSub):
    while True:
        mensagem = discord.ouvirMensagem()
        if(mensagem == "ouvindo"):
            pubsub.enviarMensagem()

def processarMensagemPubSub(discord: BotDiscord, pubsub: RedisPubSub):
    while True:
        dados = pubsub.ouvirMensagem()["data"].decode("utf-8")
        mensagem = AdmineMessage.from_json_to_object(dados)
        match mensagem.getTags():
            case ["server_up"]:
                #pubsub.enviarMensagem(mensagem)
                print(mensagem.getMessage())
            case ["server_up","down"]:
                pubsub.enviarMensagem(mensagem)
            case _:
                print("mensagem inválida")
                
    



gilsepi = AdmineMessage(["server_up"],"ola mundo")

bot_discord = BotDiscord("comandos", ["Gustavo","Gilsepi"])

meuPubSub = RedisPubSub("localhost",6379,["teste"],["teste"])


@bot_discord.tree.command(name="serverup",description="Subir server na nuvem")
async def subirServer(interaction:discord.Interaction):
    await interaction.response.send_message(f"subindo server obrigado {interaction.user.mention}!")
    #pubsub = RedisPubSub("localhost",6379,["teste"],["teste"])
    #top = AdmineMessage(["server_up"],"ola mundo")
    #pubsub.enviarMensagem(top)
    

thread = threading.Thread(target=processarMensagemPubSub, args=(bot_discord, meuPubSub))
thread.start()

DISCORD_BOT_TOKEN = "colocar o token do bot"#os.getenv("DISCORD_BOT_TOKEN")

bot_discord.run(DISCORD_BOT_TOKEN)


#thread = threading.Thread(target=processarMensagemPubSub, args=(bot_discord, meuPubSub))
#thread.start()
#print("ok")
#meuPubSub.enviarMensagem(gilsepi)

#thread2 = threading.Thread(target=processarMensagemDiscord(discord,meuPubSub))
#thread2.start()
