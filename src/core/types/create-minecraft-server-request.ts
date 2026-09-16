import type { MinecraftServerType } from "../domain/minecraft-server-type";

export type CreateMinecraftServerRequest = {
    name: string;
    readonly type: MinecraftServerType;
    version: string;
    minimumMemory?: number;
    maximumMemory?: number;
    port?: number;
};
