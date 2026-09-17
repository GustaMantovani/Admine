import type { MinecraftServerProvider } from "../providers/minecraft-server-provider";

export class DeleteMinecraftServerUseCase {
    constructor(private readonly minecraftServerProvider: MinecraftServerProvider) { }

    async execute(serverId: string, removeFiles: boolean): Promise<void> {
        return this.minecraftServerProvider.delete(serverId, removeFiles);
    }
}
