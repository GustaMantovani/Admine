#!/bin/bash

# Atualiza o gerenciador de pacotes
pip install --upgrade pip

# Instala as dependências principais
pip install fastapi uvicorn pydantic

echo "Dependências instaladas com sucesso!"