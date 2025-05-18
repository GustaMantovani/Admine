FROM ubuntu:latest

ARG JAVA_VERSION
ARG NETWORK_ID

ENV GRAALVM_HOME=/opt/graalvm
ENV JAVA_HOME=$GRAALVM_HOME
ENV PATH=$JAVA_HOME/bin:$PATH
ENV NETWORK_ID=${NETWORK_ID}

RUN apt update && \
    apt install -y curl zip unzip libz-dev libstdc++6 ca-certificates && \
    rm -rf /var/lib/apt/lists/*

RUN curl -L -o graalvm.tar.gz https://download.oracle.com/graalvm/${JAVA_VERSION}/latest/graalvm-jdk-${JAVA_VERSION}_linux-x64_bin.tar.gz && \
    mkdir -p /opt && \
    tar -xzf graalvm.tar.gz -C /opt && \
    mv /opt/graalvm-jdk-* $GRAALVM_HOME && \
    rm graalvm.tar.gz

RUN curl -s https://install.zerotier.com | bash

RUN zerotier-one & sleep 1 && \
    if [ -n "$NETWORK_ID" ]; then zerotier-cli join $NETWORK_ID; fi

EXPOSE 9993/udp
EXPOSE 25565/tcp

WORKDIR /srv/minecraft

RUN curl -O https://maven.minecraftforge.net/net/minecraftforge/forge/1.20.1-47.3.0/forge-1.20.1-47.3.0-installer.jar && \
    java -jar forge-1.20.1-47.3.0-installer.jar --installServer