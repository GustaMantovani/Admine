import type { MinecraftServerProvider } from "../providers/minecraft-server-provider";

export class StopAllMinecraftServersUseCase {
    constructor(private readonly minecraftServerProvider: MinecraftServerProvider) { }

    async execute(): Promise<number> {
        return this.minecraftServerProvider.stopAll();
    }
}
