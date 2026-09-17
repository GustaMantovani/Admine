import type { MinecraftServerProvider } from "../providers/minecraft-server-provider";

export class GetMinecraftServerLogsUseCase {
    constructor(private readonly minecraftServerProvider: MinecraftServerProvider) { }

    async execute(serverId: string, lastLines?: number): Promise<string[]> {
        const lines = await this.minecraftServerProvider.logs(serverId);
        return lastLines === undefined ? lines : lines.slice(-lastLines);
    }
}
