#!/bin/sh

pip install -r requirements.txt

if [ ! $? -eq 0 ]; then
    . /etc/os-release

    case $NAME in
        "Ubuntu")
            sudo apt update
            sudo apt install python3-redis python3-dotenv
            ;;
        "Arch Linux"|"Manjaro Linux")
            sudo pacman -Syu python-redis python-dotenv
            ;;
        *)
            echo "Unsupported OS"
            exit 1
            ;;
    esac
fi