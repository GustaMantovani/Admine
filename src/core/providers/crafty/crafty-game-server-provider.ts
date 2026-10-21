import type { CreateGameServerRequest } from "../../domain/create-game-server-request";
import { GameServerActionError } from "../../domain/errors";
import type { GameServer } from "../../domain/game-server";
import type { GameServerFile } from "../../domain/game-server-file";
import type { GameServerFileEntry } from "../../domain/game-server-file-entry";
import type { GameServerStats } from "../../domain/game-server-stats";
import type { GameServerWebhookTrigger } from "../../domain/game-server-webhook-trigger";
import type { GameServerProvider } from "../game-server-provider";
import { CRAFTY_PROVIDER_NAME } from "../provider-name";
import type { CraftyApiClient } from "./crafty-api-client";
import type { CraftyWebhooksConfig } from "./crafty-config";
import type { CraftyAction } from "./dtos/crafty-action";
import type { CraftyActionExecutionResponse } from "./dtos/crafty-action-execution-response";
import {
	fromCraftyLogLine,
	isCraftyDirectory,
	isSupportedCraftyServerType,
	toCraftyCreateServerRequest,
	toCraftyCreateWebhookRequest,
	toCraftyFileCreateRequest,
	toCraftyFileDeleteRequest,
	toCraftyFileOperationRequest,
	toCraftyFilePatchRequest,
	toCraftyFileRenameRequest,
	toGameServer,
	toGameServerFile,
	toGameServerFileEntries,
	toGameServerStats,
} from "./mappers";
import type { CraftyRconPortRepository } from "./persistence/crafty-rcon-port-repository";

export class CraftyGameServerProvider implements GameServerProvider {
	readonly name = CRAFTY_PROVIDER_NAME;

	constructor(
		private readonly craftyApiClient: CraftyApiClient,
		private readonly serversHost: string,
		private readonly webhooksConfig: CraftyWebhooksConfig,
		private readonly rconPortRepository: CraftyRconPortRepository,
	) {}

	async create(createServerRequest: CreateGameServerRequest): Promise<string> {
		const craftyCreateServerRequest = toCraftyCreateServerRequest(
			createServerRequest,
			this.serversHost,
		);
		const result = await this.craftyApiClient.createServer(
			craftyCreateServerRequest,
		);
		const serverId = result.data.new_server_id;

		if (createServerRequest.rconPort !== undefined)
			await this.rconPortRepository.set(serverId, createServerRequest.rconPort);
		if (this.webhooksConfig.enabled) await this.createWebhooks(serverId);
		return serverId;
	}

	async list(): Promise<GameServer[]> {
		const [servers, rconPorts] = await Promise.all([
			this.craftyApiClient.listServers(),
			this.rconPortRepository.getAll(),
		]);

		return servers.data
			.filter((dto) => isSupportedCraftyServerType(dto.type))
			.map((dto) =>
				toGameServer(
					dto,
					this.serversHost,
					undefined,
					rconPorts.get(dto.server_id),
				),
			);
	}

	async info(id: string): Promise<GameServer> {
		const [server, stats, rconPort] = await Promise.all([
			this.craftyApiClient.getServer(id),
			this.craftyApiClient.getServerStats(id),
			this.rconPortRepository.get(id),
		]);
		return toGameServer(server.data, this.serversHost, stats.data, rconPort);
	}

	async stats(id: string): Promise<GameServerStats> {
		const result = await this.craftyApiClient.getServerStats(id);
		return toGameServerStats(result.data);
	}

	async logs(id: string): Promise<string[]> {
		const result = await this.craftyApiClient.getServerLogs(id);
		return result.data.map(fromCraftyLogLine);
	}

	async start(id: string): Promise<void> {
		return this.runAction(id, "start_server");
	}

	async stop(id: string): Promise<void> {
		return this.runAction(id, "stop_server");
	}

	async restart(id: string): Promise<void> {
		return this.runAction(id, "restart_server");
	}

	async kill(id: string): Promise<void> {
		return this.runAction(id, "kill_server");
	}

	async backup(id: string): Promise<void> {
		return this.runAction(id, "backup_server");
	}

	async update(id: string): Promise<void> {
		return this.runAction(id, "update_executable");
	}

	async clone(id: string): Promise<void> {
		return this.runAction(id, "clone_server");
	}

	async acceptEula(id: string): Promise<void> {
		return this.runAction(id, "eula");
	}

	async sendCommand(id: string, command: string): Promise<void> {
		const result = await this.craftyApiClient.sendCommand(id, command);
		this.throwIfFailed(result, `send command to server ${id}`);
	}

	async listFiles(id: string, path: string): Promise<GameServerFileEntry[]> {
		const result = await this.craftyApiClient.getServerFiles(id, path);
		if (!isCraftyDirectory(result.data))
			throw new GameServerActionError(`${path} is not a directory`);
		return toGameServerFileEntries(result.data);
	}

	async getFile(id: string, path: string): Promise<GameServerFile> {
		const result = await this.craftyApiClient.getServerFiles(id, path);
		if (isCraftyDirectory(result.data))
			throw new GameServerActionError(`${path} is a directory, not a file`);
		return toGameServerFile(path, result.data);
	}

	async writeFile(id: string, file: GameServerFile): Promise<void> {
		return this.craftyApiClient.updateServerFile(
			id,
			toCraftyFilePatchRequest(file),
		);
	}

	async createFile(id: string, path: string): Promise<void> {
		const result = await this.craftyApiClient.createServerFile(
			id,
			toCraftyFileCreateRequest(path, false),
		);
		this.throwIfFailed(result, `create file ${path} on server ${id}`);
	}

	async createDirectory(id: string, path: string): Promise<void> {
		const result = await this.craftyApiClient.createServerFile(
			id,
			toCraftyFileCreateRequest(path, true),
		);
		this.throwIfFailed(result, `create directory ${path} on server ${id}`);
	}

	async renameFile(id: string, path: string, newName: string): Promise<void> {
		const result = await this.craftyApiClient.renameServerFile(
			id,
			toCraftyFileRenameRequest(path, newName),
		);
		this.throwIfFailed(result, `rename ${path} on server ${id}`);
	}

	async moveFiles(
		id: string,
		sourcePaths: string[],
		targetDirectory: string,
	): Promise<void> {
		const result = await this.craftyApiClient.runServerFileOperation(
			id,
			"move",
			toCraftyFileOperationRequest(sourcePaths, targetDirectory),
		);
		this.throwIfFailed(result, `move files on server ${id}`);
	}

	async copyFiles(
		id: string,
		sourcePaths: string[],
		targetDirectory: string,
	): Promise<void> {
		const result = await this.craftyApiClient.runServerFileOperation(
			id,
			"copy",
			toCraftyFileOperationRequest(sourcePaths, targetDirectory),
		);
		this.throwIfFailed(result, `copy files on server ${id}`);
	}

	async deleteFiles(id: string, paths: string[]): Promise<void> {
		const result = await this.craftyApiClient.deleteServerFiles(
			id,
			toCraftyFileDeleteRequest(paths),
		);
		this.throwIfFailed(result, `delete files on server ${id}`);
	}

	async downloadFile(
		id: string,
		path: string,
	): Promise<ReadableStream<Uint8Array>> {
		return this.craftyApiClient.downloadServerFile(id, path);
	}

	async delete(id: string, removeFiles: boolean): Promise<void> {
		const result = await this.craftyApiClient.deleteServer(id, removeFiles);
		this.throwIfFailed(result, `delete server ${id}`);
		await this.rconPortRepository.delete(id);
	}

	async stopAll(): Promise<number> {
		const result = await this.craftyApiClient.stopAllServers();
		return result.data.count;
	}

	private async createWebhooks(serverId: string): Promise<void> {
		const webhooks = Object.entries(this.webhooksConfig.urls) as [
			GameServerWebhookTrigger,
			string,
		][];

		for (const [trigger, url] of webhooks) {
			await this.craftyApiClient.createWebhook(
				serverId,
				toCraftyCreateWebhookRequest(trigger, url),
			);
		}
	}

	private async runAction(id: string, action: CraftyAction): Promise<void> {
		const result = await this.craftyApiClient.sendAction(id, action);
		this.throwIfFailed(result, `run ${action} on server ${id}`);
	}

	// Some Crafty failures come back as HTTP 200 with status "error"
	private throwIfFailed(
		result: CraftyActionExecutionResponse,
		operation: string,
	): void {
		if (result.status === "error") {
			const details = result.error_data ? ` (${result.error_data})` : "";
			throw new GameServerActionError(
				`Crafty failed to ${operation}: ${result.error ?? "unknown error"}${details}`,
			);
		}
	}
}
