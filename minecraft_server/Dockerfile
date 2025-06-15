FROM archlinux:latest

ARG JAVA_VERSION
ARG FORGE_VERSION
ARG NETWORK_ID

ENV JAVA_VERSION=${JAVA_VERSION:-17}
ENV FORGE_VERSION=${FORGE_VERSION:-1.20.1-47.4.0}
ENV NETWORK_ID=${NETWORK_ID}

RUN pacman -Syu --noconfirm --needed --noprogressbar && \
    pacman -S --noconfirm --needed --noprogressbar jre${JAVA_VERSION}-openjdk-headless zerotier-one && \
    pacman -S --noconfirm --needed --noprogressbar curl git unzip wget && \
    if [ -n "$(pacman -Qdtq)" ]; then pacman -Rsn $(pacman -Qdtq) --noconfirm; fi

RUN zerotier-one & sleep 1 && \
    if [ -n "$NETWORK_ID" ]; then zerotier-cli join $NETWORK_ID; fi

EXPOSE 9993/udp
EXPOSE 25565/tcp

WORKDIR /srv/minecraft

RUN curl -O https://maven.minecraftforge.net/net/minecraftforge/forge/${FORGE_VERSION}/forge-${FORGE_VERSION}-installer.jar && \
    java -jar forge-${FORGE_VERSION}-installer.jar --installServer

CMD zerotier-one & ./run.sh