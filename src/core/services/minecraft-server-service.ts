import type { CreateMinecraftServerRequest } from "../domain/create-minecraft-server-request";
import type { MinecraftServer } from "../domain/minecraft-server";
import type { MinecraftServerStats } from "../domain/minecraft-server-stats";
import type { MinecraftServerProvider } from "../providers/minecraft-server-provider";
import { sendCommand } from "../rcon/rcon-utils";

export class MinecraftServerService {
	constructor(
		private readonly minecraftServerProvider: MinecraftServerProvider,
	) {}

	async create(
		createServerRequest: CreateMinecraftServerRequest,
		acceptEula?: boolean,
	): Promise<string> {
		const newServerId =
			await this.minecraftServerProvider.create(createServerRequest);

		if (acceptEula) await this.acceptEula(newServerId);

		return newServerId;
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

	async sendCommand(_serverId: string, command: string): Promise<string> {
		// Aqui eu me viro pra arrumar as credenciais baseado no id do servidor
		return sendCommand(
			"localhost",
			25575,
			"S5Yg4LTmjlZb-8DfwcejJA5fnID9t65cC6zpqiUB7XM",
			command,
		);
	}

	async delete(serverId: string, removeFiles: boolean): Promise<void> {
		return this.minecraftServerProvider.delete(serverId, removeFiles);
	}

	async stopAll(): Promise<number> {
		return this.minecraftServerProvider.stopAll();
	}
}
