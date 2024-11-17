FROM ubuntu:latest

# Argumento para especificar a versão do Java
ARG JAVA_VERSION=17

RUN apt-get update && \
    apt-get install -y wget software-properties-common curl && \
    rm -rf /var/lib/apt/lists/*

# Adiciona o repositório do OpenJDK e instala a versão escolhida
RUN add-apt-repository ppa:openjdk-r/ppa && \
    apt-get update && \
    apt-get install -y openjdk-${JAVA_VERSION}-jdk && \
    rm -rf /var/lib/apt/lists/*

# Define a variável de ambiente JAVA_HOME
ENV JAVA_HOME /usr/lib/jvm/java-${JAVA_VERSION}-openjdk-amd64
ENV PATH $JAVA_HOME/bin:$PATH

# Instalação do ZeroTier
RUN curl -s https://install.zerotier.com | bash

# Inicia o serviço ZeroTier e ingressa na rede especificada, se o NETWORK_ID estiver definido
ARG NETWORK_ID=52b337794f4e7176
RUN zerotier-one & sleep 1 && \
    if [ -n "$NETWORK_ID" ]; then zerotier-cli join $NETWORK_ID; fi

# Abre a porta padrão do ZeroTier (9993/UDP) para permitir a comunicação
EXPOSE 9993/udp
EXPOSE 25565/tcp

# ============= Fabric ============= #

WORKDIR /srv/minecraft
RUN curl -OJ https://meta.fabricmc.net/v2/versions/loader/1.20.1/0.16.5/1.0.1/server/jar


COPY ./config/eula.txt /srv/minecraft/
COPY ./config/server.properties /srv/minecraft/
COPY ./config/user_jvm_args.txt /srv/minecraft/

CMD zerotier-one & sleep 1 && java -Xmx7G -jar fabric-server-mc.1.20.1-loader.0.16.5-launcher.1.0.1.jar nogui

# ============= Forge ============= #

#RUN curl -O https://maven.minecraftforge.net/net/minecraftforge/forge/1.20.1-47.3.0/forge-1.20.1-47.3.0-installer.jar
#RUN java -jar forge-1.20.1-47.3.0-installer.jar --installServer

# COPY ./config/eula.txt /srv/minecraft/
# COPY ./config/server.properties /srv/minecraft/
# COPY ./config/user_jvm_args.txt /srv/minecraft/

#CMD zerotier-one & sleep 1 && ./run.sh