FROM ubuntu:latest

# Argumento para especificar a versão do Java e id da rede zerotier
ARG JAVA_VERSION
ARG NETWORK_ID

# Define a variável de ambiente JAVA_HOME
ENV JAVA_HOME /usr/lib/jvm/java-${JAVA_VERSION}-openjdk-amd64
ENV PATH $JAVA_HOME/bin:$PATH
ENV NETWORK_ID=${NETWORK_ID}

# Atualiza e instala dependências essenciais
RUN apt update && \
    apt install -y software-properties-common curl && \
    rm -rf /var/lib/apt/lists/*

# Adiciona o repositório do OpenJDK e instala a versão escolhida
RUN add-apt-repository ppa:openjdk-r/ppa && \
    apt update && \
    apt install -y openjdk-${JAVA_VERSION}-jdk && \
    rm -rf /var/lib/apt/lists/*

# Instalação do ZeroTier
RUN curl -s https://install.zerotier.com | bash

# Inicia o serviço ZeroTier e ingressa na rede especificada, se o NETWORK_ID estiver definido
RUN zerotier-one & sleep 1 && \
    if [ -n "$NETWORK_ID" ]; then zerotier-cli join $NETWORK_ID; fi

# Abre a porta padrão do ZeroTier (9993/UDP) para permitir a comunicação
EXPOSE 9993/udp
EXPOSE 25565/tcp

# Servidor Minecraft
WORKDIR /srv/minecraft

#FORGE
RUN curl -O https://maven.minecraftforge.net/net/minecraftforge/forge/1.20.1-47.3.0/forge-1.20.1-47.3.0-installer.jar
RUN java -jar forge-1.20.1-47.3.0-installer.jar --installServer

# Comando para iniciar o serviço ZeroTier e o servidor Minecraft
CMD zerotier-one & ./run.sh