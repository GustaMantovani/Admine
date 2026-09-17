import type { MinecraftServerProvider } from "../providers/minecraft-server-provider";

export class SendMinecraftServerCommandUseCase {
    constructor(private readonly minecraftServerProvider: MinecraftServerProvider) { }

    async execute(serverId: string, command: string): Promise<void> {
        return this.minecraftServerProvider.sendCommand(serverId, command);
    }
}
