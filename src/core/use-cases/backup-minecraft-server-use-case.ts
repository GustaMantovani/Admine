import type { MinecraftServerProvider } from "../providers/minecraft-server-provider";

export class BackupMinecraftServerUseCase {
    constructor(private readonly minecraftServerProvider: MinecraftServerProvider) { }

    async execute(serverId: string): Promise<void> {
        return this.minecraftServerProvider.backup(serverId);
    }
}
