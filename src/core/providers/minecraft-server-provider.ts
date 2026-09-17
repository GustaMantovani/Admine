import type { MinecraftServer } from "../domain/minecraft-server";
import type { CreateMinecraftServerRequest } from "../types/create-minecraft-server-request";

export interface MinecraftServerProvider {
    create(createServerRequest: CreateMinecraftServerRequest): Promise<string>;
    list(): Promise<MinecraftServer[]>;
    start(id: string): Promise<void>;
}