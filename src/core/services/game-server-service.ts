import { PropertiesEditor } from "properties-file/editor";
import type { CreateGameServerRequest } from "../domain/create-game-server-request";
import type { GameServer } from "../domain/game-server";
import type { GameServerStats } from "../domain/game-server-stats";
import type { GameServerType } from "../domain/game-server-type";
import {
	DEFAULT_GAME_PORTS,
	DEFAULT_RCON_PORTS,
	GAME_PORT_RANGES,
	RCON_PORT_RANGES,
} from "../domain/server-ports";
import type { GameServerProvider } from "../providers/game-server-provider";
import type { RconConfig } from "../rcon/rcon-config";
import { sendCommand } from "../rcon/rcon-utils";
import type { GameServerPortsRepository } from "../repositories/game-server-ports-repository";
import { randomOutside } from "../utils";

export class GameServerService {
	constructor(
		private readonly gameServerProvider: GameServerProvider,
		private readonly gameServerPortsRepository: GameServerPortsRepository,
		private readonly rconConfig: RconConfig,
	) {}

	private buildServerProperties(
		request: CreateGameServerRequest,
	): Record<string, string> {
		const properties: Record<string, string> = {};

		if (this.rconConfig.enabled && request.type === "java") {
			properties["enable-rcon"] = "true";
			properties["rcon.password"] = this.rconConfig.password;
			if (request.rconPort !== undefined)
				properties["rcon.port"] = String(request.rconPort);
			properties["broadcast-rcon-to-ops"] = "false";
		}

		// User configurable properties
		for (const [key, value] of Object.entries(request.properties ?? {})) {
			properties[key] = String(value);
		}

		return properties;
	}

	private async writeServerProperties(
		serverId: string,
		properties: Record<string, string>,
	): Promise<void> {
		const file = await this.gameServerProvider.getFile(
			serverId,
			"server.properties",
		);

		const propertiesEditor = new PropertiesEditor(file.content);
		for (const [property, value] of Object.entries(properties)) {
			propertiesEditor.insert(property, value);
		}

		const content = propertiesEditor.format({
			removeComments: true, // Strip all comments
			removeBlankLines: true, // Strip all blank lines
			removeLeadingWhitespace: true, // Strip indentation
			deduplicateKeys: true, // Keep only last occurrence
			separatorChar: "=", // Standardize separator
			separatorLeading: "", // Space before =
			separatorTrailing: "", // Space after =
		});

		await this.gameServerProvider.writeFile(serverId, { ...file, content });
	}

	async getRandomUnusedGameServerPort(
		serverType: GameServerType,
	): Promise<number> {
		const usedPorts = await this.gameServerPortsRepository.getUsedGamePorts();
		let randomPort = randomOutside(
			GAME_PORT_RANGES[serverType].min,
			GAME_PORT_RANGES[serverType].max,
			usedPorts,
		);
		if (randomPort === null) randomPort = DEFAULT_GAME_PORTS[serverType];

		return randomPort;
	}

	async getRandomUnusedRconPort(
		serverType: GameServerType,
	): Promise<number | undefined> {
		if (serverType !== "java") return undefined;

		const usedPorts = await this.gameServerPortsRepository.getUsedRconPorts();
		let randomPort = randomOutside(
			RCON_PORT_RANGES[serverType].min,
			RCON_PORT_RANGES[serverType].max,
			usedPorts,
		);
		if (randomPort === null) randomPort = DEFAULT_RCON_PORTS[serverType];

		return randomPort;
	}

	async create(
		createServerRequest: CreateGameServerRequest,
		acceptEula?: boolean,
	): Promise<string> {
		const newServerId =
			await this.gameServerProvider.create(createServerRequest);

		await this.writeServerProperties(
			newServerId,
			this.buildServerProperties(createServerRequest),
		);

		if (acceptEula) await this.acceptEula(newServerId);

		await this.gameServerPortsRepository.create({
			id: newServerId,
			port: createServerRequest.port,
			rconPort: createServerRequest.rconPort,
		});

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
		await this.gameServerProvider.delete(serverId, removeFiles);
		await this.gameServerPortsRepository.delete(serverId);
	}

	async stopAll(): Promise<number> {
		return this.gameServerProvider.stopAll();
	}
}
