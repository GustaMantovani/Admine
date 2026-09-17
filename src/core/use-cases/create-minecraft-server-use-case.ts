import type { MinecraftServer } from "../domain/minecraft-server";
import type { MinecraftServerProvider } from "../providers/minecraft-server-provider";
import type { CreateMinecraftServerRequest } from "../types/create-minecraft-server-request";

export class CreateMinecraftServerUseCase {
    constructor(private readonly minecraftServerProvider: MinecraftServerProvider) {}

    async execute(createServerRequest: CreateMinecraftServerRequest): Promise<string> {
        return this.minecraftServerProvider.create(createServerRequest);
    }
}