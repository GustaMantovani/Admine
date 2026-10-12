import type { CreateGameServerRequest } from "../domain/create-game-server-request";
import type { GameServer } from "../domain/game-server";
import type { GameServerStats } from "../domain/game-server-stats";
import type { GameServerProvider } from "../providers/game-server-provider";
import { sendCommand } from "../rcon/rcon-utils";

export class GameServerService {
	constructor(private readonly gameServerProvider: GameServerProvider) {}

	async create(
		createServerRequest: CreateGameServerRequest,
		acceptEula?: boolean,
	): Promise<string> {
		const newServerId =
			await this.gameServerProvider.create(createServerRequest);

		if (acceptEula) await this.acceptEula(newServerId);

		return newServerId;
	}

	async list(): Promise<GameServer[]> {
		return this.gameServerProvider.list();
	}

	async get(serverId: string): Promise<GameServer> {
		return this.gameServerProvider.info(serverId);
	}

	async stats(serverId: string): Promise<GameServerStats> {
		return this.gameServerProvider.stats(serverId);
	}

	async logs(serverId: string, lastLines?: number): Promise<string[]> {
		const lines = await this.gameServerProvider.logs(serverId);
		return lastLines === undefined ? lines : lines.slice(-lastLines);
	}

	async start(serverId: string): Promise<void> {
		return this.gameServerProvider.start(serverId);
	}

	async stop(serverId: string): Promise<void> {
		return this.gameServerProvider.stop(serverId);
	}

	async restart(serverId: string): Promise<void> {
		return this.gameServerProvider.restart(serverId);
	}

	async kill(serverId: string): Promise<void> {
		return this.gameServerProvider.kill(serverId);
	}

	async backup(serverId: string): Promise<void> {
		return this.gameServerProvider.backup(serverId);
	}

	async update(serverId: string): Promise<void> {
		return this.gameServerProvider.update(serverId);
	}

	async clone(serverId: string): Promise<void> {
		return this.gameServerProvider.clone(serverId);
	}

	async acceptEula(serverId: string): Promise<void> {
		return this.gameServerProvider.acceptEula(serverId);
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
		return this.gameServerProvider.delete(serverId, removeFiles);
	}

	async stopAll(): Promise<number> {
		return this.gameServerProvider.stopAll();
	}
}
