import type { ApiV2ServersPostRequest, CreateBedrockServer, CreateHytaleServer, CreateJavaServer } from "../../../../dependencies/crafty-api-client";
import type { CreateMinecraftServerRequest } from "../../types/create-minecraft-server-request";

const DEFAULT_JAVA_PORT = 25565;
const DEFAULT_BEDROCK_PORT = 19132;
const DEFAULT_HYTALE_PORT = 5520;
const DEFAULT_MINIMUM_MEMORY = 1;
const DEFAULT_MAXIMUM_MEMORY = 2;

export function resolvePort(request: CreateMinecraftServerRequest): number {
    if (request.port != null) {
        return request.port;
    }

    switch (request.type) {
        case "java":
            return DEFAULT_JAVA_PORT;
        case "bedrock":
            return DEFAULT_BEDROCK_PORT;
        case "hytale":
            return DEFAULT_HYTALE_PORT;
    }
}

export function toCraftyCreateServerRequest(request: CreateMinecraftServerRequest): ApiV2ServersPostRequest {
    switch (request.type) {
        case "java":
            return toCraftyCreateJavaServer(request);
        case "bedrock":
            return toCraftyCreateBedrockServer(request);
        case "hytale":
            return toCraftyCreateHytaleServer(request);
    }
}

function toCraftyCreateJavaServer(request: CreateMinecraftServerRequest): CreateJavaServer {
    const port = resolvePort(request);

    return {
        name: request.name,
        monitoringType: "minecraft_java",
        minecraftJavaMonitoringData: {
            host: "127.0.0.1",
            port,
        },
        createType: "minecraft_java",
        minecraftJavaCreateData: {
            createType: "download_jar",
            downloadJarCreateData: {
                category: "mc_java_servers",
                type: "paper",
                version: request.version,
                memMin: request.minimumMemory ?? DEFAULT_MINIMUM_MEMORY,
                memMax: request.maximumMemory ?? DEFAULT_MAXIMUM_MEMORY,
                serverPropertiesPort: port,
            },
        },
    };
}

function toCraftyCreateBedrockServer(request: CreateMinecraftServerRequest): CreateBedrockServer {
    return {
        name: request.name,
        monitoringType: "minecraft_bedrock",
        minecraftBedrockMonitoringData: {
            host: "127.0.0.1",
            port: resolvePort(request),
        },
        createType: "minecraft_bedrock",
        minecraftBedrockCreateData: {
            createType: "download_exe",
            downloadExeCreateData: {
                agreeToEula: true,
            },
        },
    };
}

function toCraftyCreateHytaleServer(request: CreateMinecraftServerRequest): CreateHytaleServer {
    return {
        name: request.name,
        monitoringType: "hytale",
        hytaleMonitoringData: {
            host: "127.0.0.1",
            port: resolvePort(request),
        },
        createType: "hytale",
        hytaleCreateData: {
            createType: "download_exe",
            downloadExeCreateData: {
                agreeToEula: true,
                memMin: request.minimumMemory ?? DEFAULT_MINIMUM_MEMORY,
                memMax: request.maximumMemory ?? DEFAULT_MAXIMUM_MEMORY,
            },
        },
    };
}
