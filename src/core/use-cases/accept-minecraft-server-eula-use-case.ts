import type { MinecraftServerProvider } from "../providers/minecraft-server-provider";

export class AcceptMinecraftServerEulaUseCase {
    constructor(private readonly minecraftServerProvider: MinecraftServerProvider) { }

    async execute(serverId: string): Promise<void> {
        return this.minecraftServerProvider.acceptEula(serverId);
    }
}
