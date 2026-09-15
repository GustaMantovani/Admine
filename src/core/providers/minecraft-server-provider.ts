import type { MinecraftServer } from "../domain/minecraft-server";
import type { CreateMinecraftServerRequest } from "../types/create-minecraft-server-request";

export interface MinecraftServerProvider {
    createServer(createServerRequest: CreateMinecraftServerRequest): Promise<string>;
    listServers(): Promise<MinecraftServer[]>;
}