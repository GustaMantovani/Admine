FROM ubuntu:latest

LABEL org.opencontainers.image.description="Minecraft Server with Forge and GraalVM"

ARG JAVA_VERSION
ARG FORGE_VERSION
ARG MINECRAFT_VERSION
ARG ATM_SERVER_FILES_DOWNLOAD_URL

ENV JAVA_VERSION=${JAVA_VERSION:-21}
ENV FORGE_VERSION=${FORGE_VERSION:-47.4.0}
ENV MINECRAFT_VERSION=${MINECRAFT_VERSION:-1.20.1}
ENV ATM_SERVER_FILES_DOWNLOAD_URL=${ATM_SERVER_FILES_DOWNLOAD_URL:-}
ENV GRAALVM_HOME=/opt/graalvm  
ENV JAVA_HOME=$GRAALVM_HOME
ENV PATH=$JAVA_HOME/bin:$PATH

RUN apt update && \
    apt install -y curl zip unzip libz-dev libstdc++6 ca-certificates && \
    rm -rf /var/lib/apt/lists/*

RUN curl -L -o graalvm.tar.gz https://download.oracle.com/graalvm/${JAVA_VERSION}/latest/graalvm-jdk-${JAVA_VERSION}_linux-x64_bin.tar.gz && \
    mkdir -p /opt && \
    tar -xzf graalvm.tar.gz -C /opt && \
    mv /opt/graalvm-jdk-* $GRAALVM_HOME && \
    rm graalvm.tar.gz

RUN curl -s https://install.zerotier.com | bash

EXPOSE 9993/udp
EXPOSE 25565/tcp

WORKDIR /srv/minecraft

RUN curl -o atm.zip ${ATM_SERVER_FILES_DOWNLOAD_URL} && \
    unzip atm.zip -d serverfiles && \
    rm atm.zip && \
    mv serverfiles/*/* . && \
    rm -rf serverfiles && \
    chmod +x startserver.sh && \
    java -jar forge-${MINECRAFT_VERSION}-${FORGE_VERSION}-installer.jar --installServer

CMD zerotier-one & sleep 1 && if [ -n "$NETWORK_ID" ]; then zerotier-cli join $NETWORK_ID; fi && ./startserver.sh