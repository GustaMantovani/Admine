import type { MinecraftServerProvider } from "../providers/minecraft-server-provider";
import type { CreateMinecraftServerRequest } from "../types/create-minecraft-server-request";
import type { MinecraftServer } from "../types/minecraft-server";
import type { MinecraftServerStats } from "../types/minecraft-server-stats";

export class MinecraftServerUseCases {
	constructor(
		private readonly minecraftServerProvider: MinecraftServerProvider,
	) {}

	async create(
		createServerRequest: CreateMinecraftServerRequest,
	): Promise<string> {
		return this.minecraftServerProvider.create(createServerRequest);
	}

	async list(): Promise<MinecraftServer[]> {
		return this.minecraftServerProvider.list();
	}

	async get(serverId: string): Promise<MinecraftServer> {
		return this.minecraftServerProvider.info(serverId);
	}

	async stats(serverId: string): Promise<MinecraftServerStats> {
		return this.minecraftServerProvider.stats(serverId);
	}

	async logs(serverId: string, lastLines?: number): Promise<string[]> {
		const lines = await this.minecraftServerProvider.logs(serverId);
		return lastLines === undefined ? lines : lines.slice(-lastLines);
	}

	async start(serverId: string): Promise<void> {
		return this.minecraftServerProvider.start(serverId);
	}

	async stop(serverId: string): Promise<void> {
		return this.minecraftServerProvider.stop(serverId);
	}

	async restart(serverId: string): Promise<void> {
		return this.minecraftServerProvider.restart(serverId);
	}

	async kill(serverId: string): Promise<void> {
		return this.minecraftServerProvider.kill(serverId);
	}

	async backup(serverId: string): Promise<void> {
		return this.minecraftServerProvider.backup(serverId);
	}

	async update(serverId: string): Promise<void> {
		return this.minecraftServerProvider.update(serverId);
	}

	async clone(serverId: string): Promise<void> {
		return this.minecraftServerProvider.clone(serverId);
	}

	async acceptEula(serverId: string): Promise<void> {
		return this.minecraftServerProvider.acceptEula(serverId);
	}

	async sendCommand(serverId: string, command: string): Promise<void> {
		return this.minecraftServerProvider.sendCommand(serverId, command);
	}

	async delete(serverId: string, removeFiles: boolean): Promise<void> {
		return this.minecraftServerProvider.delete(serverId, removeFiles);
	}

	async stopAll(): Promise<number> {
		return this.minecraftServerProvider.stopAll();
	}
}
