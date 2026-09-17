import type { MinecraftServer } from "../domain/minecraft-server";
import type { MinecraftServerProvider } from "../providers/minecraft-server-provider";

export class GetMinecraftServerUseCase {
    constructor(private readonly minecraftServerProvider: MinecraftServerProvider) { }

    async execute(serverId: string): Promise<MinecraftServer> {
        return this.minecraftServerProvider.get(serverId);
    }
}
