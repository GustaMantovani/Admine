import type { MinecraftServerStats } from "../domain/minecraft-server-stats";
import type { MinecraftServerProvider } from "../providers/minecraft-server-provider";

export class GetMinecraftServerStatsUseCase {
    constructor(private readonly minecraftServerProvider: MinecraftServerProvider) { }

    async execute(serverId: string): Promise<MinecraftServerStats> {
        return this.minecraftServerProvider.stats(serverId);
    }
}
