import type { BunFile } from "bun";
import { PropertiesEditor } from "properties-file/editor";
import type { CreateGameServerRequest } from "../domain/create-game-server-request";
import { GameServerActionError } from "../domain/errors";
import type { GameServer } from "../domain/game-server";
import type { GameServerFile } from "../domain/game-server-file";
import type { GameServerFileDownload } from "../domain/game-server-file-download";
import type { GameServerFileEntry } from "../domain/game-server-file-entry";
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
import { randomOutside } from "../utils";

export class GameServerService {
	constructor(
		private readonly gameServerProvider: GameServerProvider,
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
		const servers = await this.gameServerProvider.list();
		const usedPorts = servers.flatMap((s) =>
			s.port !== undefined ? [s.port] : [],
		);
		return (
			randomOutside(
				GAME_PORT_RANGES[serverType].min,
				GAME_PORT_RANGES[serverType].max,
				usedPorts,
			) ?? DEFAULT_GAME_PORTS[serverType]
		);
	}

	async getRandomUnusedRconPort(
		serverType: GameServerType,
	): Promise<number | undefined> {
		if (serverType !== "java") return undefined;

		const servers = await this.gameServerProvider.list();
		const usedPorts = servers.flatMap((s) =>
			s.rconPort !== undefined ? [s.rconPort] : [],
		);
		return (
			randomOutside(
				RCON_PORT_RANGES[serverType].min,
				RCON_PORT_RANGES[serverType].max,
				usedPorts,
			) ?? DEFAULT_RCON_PORTS[serverType]
		);
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

	async sendCommand(serverId: string, command: string): Promise<string> {
		const server = await this.gameServerProvider.info(serverId);
		if (!this.rconConfig.enabled || server.rconPort === undefined) {
			throw new GameServerActionError(
				`RCON is not available for server ${serverId}`,
			);
		}

		return sendCommand(
			server.admineIpAddress,
			server.rconPort,
			this.rconConfig.password,
			command,
		);
	}

	async delete(serverId: string, removeFiles: boolean): Promise<void> {
		return this.gameServerProvider.delete(serverId, removeFiles);
	}

	async stopAll(): Promise<number> {
		return this.gameServerProvider.stopAll();
	}

	async listFiles(
		serverId: string,
		path: string,
	): Promise<GameServerFileEntry[]> {
		return this.gameServerProvider.listFiles(serverId, path);
	}

	async readFile(serverId: string, path: string): Promise<GameServerFile> {
		return this.gameServerProvider.getFile(serverId, path);
	}

	// Pass the file as read: the provider rejects the write if it changed since
	async writeFile(serverId: string, file: GameServerFile): Promise<void> {
		return this.gameServerProvider.writeFile(serverId, file);
	}

	async createFile(serverId: string, path: string): Promise<void> {
		return this.gameServerProvider.createFile(serverId, path);
	}

	async createDirectory(serverId: string, path: string): Promise<void> {
		return this.gameServerProvider.createDirectory(serverId, path);
	}

	async renameFile(
		serverId: string,
		path: string,
		newName: string,
	): Promise<void> {
		return this.gameServerProvider.renameFile(serverId, path, newName);
	}

	async moveFiles(
		serverId: string,
		sourcePaths: string[],
		targetDirectory: string,
	): Promise<void> {
		return this.gameServerProvider.moveFiles(
			serverId,
			sourcePaths,
			targetDirectory,
		);
	}

	async copyFiles(
		serverId: string,
		sourcePaths: string[],
		targetDirectory: string,
	): Promise<void> {
		return this.gameServerProvider.copyFiles(
			serverId,
			sourcePaths,
			targetDirectory,
		);
	}

	async deleteFiles(serverId: string, paths: string[]): Promise<void> {
		return this.gameServerProvider.deleteFiles(serverId, paths);
	}

	async downloadFile(
		serverId: string,
		path: string,
	): Promise<GameServerFileDownload> {
		return this.gameServerProvider.downloadFile(serverId, path);
	}

	async uploadFile(
		serverId: string,
		path: string,
		file: BunFile,
	): Promise<void> {
		return this.gameServerProvider.uploadFile(serverId, path, file);
	}
}
