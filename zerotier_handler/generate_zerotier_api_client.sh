#!/bin/sh

openapi-generator-cli generate \
    -i plugin-redoc-0.yaml \
    -g rust \
    -o zerotier_api_client \
    -p topLevelApiClient=true \
    -p packageName=zerotier_api_client \
    -p packageVersion=1.0.0
