import type { MinecraftServerProvider } from "../providers/minecraft-server-provider";

export class StopMinecraftServerUseCase {
    constructor(private readonly minecraftServerProvider: MinecraftServerProvider) { }

    async execute(serverId: string): Promise<void> {
        return this.minecraftServerProvider.stop(serverId);
    }
}
