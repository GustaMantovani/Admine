import type { MinecraftServer } from "../domain/minecraft-server";
import type { MinecraftServerProvider } from "../providers/minecraft-server-provider";

export class ListMinecraftServersUseCase {
    constructor(private readonly minecraftServerProvider: MinecraftServerProvider) {}

    async execute(): Promise<MinecraftServer[]> {
        const servers = this.minecraftServerProvider.listServers();
        return servers;
    }
}