FROM alpine:latest

# Argumentos para especificar a versão do Java e ID da rede ZeroTier
ARG JAVA_VERSION
ARG NETWORK_ID

# Variáveis de ambiente
ENV JAVA_HOME=/usr/lib/jvm/java-${JAVA_VERSION}-openjdk
ENV PATH=$JAVA_HOME/bin:$PATH
ENV NETWORK_ID=${NETWORK_ID}

# Atualiza e instala dependências essenciais
RUN echo "https://dl-cdn.alpinelinux.org/alpine/v3.15/community" >> /etc/apk/repositories && \
    apk update && \
    apk add --no-cache \
    openjdk${JAVA_VERSION} \
    curl \
    zerotier-one && \
    mkdir -p /srv/minecraft

# Abre as portas necessárias
EXPOSE 9993/udp
EXPOSE 25565/tcp

# Define o diretório de trabalho
WORKDIR /srv/minecraft

# Inicia o serviço ZeroTier e ingressa na rede especificada, se o NETWORK_ID estiver definido
RUN zerotier-one & sleep 1 && \
    if [ -n "$NETWORK_ID" ]; then zerotier-cli join $NETWORK_ID; fi

#FORGE
RUN curl -O https://maven.minecraftforge.net/net/minecraftforge/forge/1.20.1-47.3.0/forge-1.20.1-47.3.0-installer.jar && \
    java -jar forge-1.20.1-47.3.0-installer.jar --installServer && \
    rm forge-1.20.1-47.3.0-installer.jar

# Comando para iniciar o serviço ZeroTier e o servidor Minecraft
CMD zerotier-one & sleep 1 && ./run.sh
