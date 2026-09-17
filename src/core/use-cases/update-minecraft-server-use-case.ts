import type { MinecraftServerProvider } from "../providers/minecraft-server-provider";

export class UpdateMinecraftServerUseCase {
    constructor(private readonly minecraftServerProvider: MinecraftServerProvider) { }

    async execute(serverId: string): Promise<void> {
        return this.minecraftServerProvider.update(serverId);
    }
}
