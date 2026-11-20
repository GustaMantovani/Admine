import { randomBytes } from "node:crypto";
import type { BunFile } from "bun";
import pLimit from "p-limit";
import type { RconClient } from "../../../../lib/rcon";
import type { CreateServerSpec } from "../../domain/create-server-spec";
import {
	InvalidServerPorts,
	OrchestratorActionError,
} from "../../domain/errors";
import type {
	CommandOutput,
	GameServerFile,
	GameServerFileDownload,
	GameServerFileEntry,
	GameServerStats,
} from "../../domain/game-server";
import type {
	GameProperties,
	GameTypeId,
} from "../../domain/game-types/game-type";
import type { ServerPorts } from "../../domain/port-purpose";
import type {
	BackupCapability,
	FileCapability,
	OrchestratorCapability,
	OrchestratorPortRequest,
	OrchestratorRef,
	OrchestratorTarget,
	ServerOrchestrator,
	UpdateCapability,
} from "../../ports/server-orchestrator";
import type { CraftyApiClient } from "./crafty-api-client";
import type { CraftyConfig, CraftyWebhookEvent } from "./crafty-config";
import type { CraftyRconRepository } from "./crafty-rcon-repository";
import type { CraftyAction, CraftyActionExecutionResponse } from "./dtos";
import {
	fromCraftyLogLine,
	isCraftyDirectory,
	toCraftyCreateServerRequest,
	toCraftyCreateWebhookRequest,
	toCraftyFileCreateRequest,
	toCraftyFileDeleteRequest,
	toCraftyFileOperationRequest,
	toCraftyFilePatchRequest,
	toCraftyFileRenameRequest,
	toGameServerFile,
	toGameServerFileEntries,
	toGameServerStats,
} from "./mappers";

const RCON_GAMES: ReadonlySet<GameTypeId> = new Set(["java"]);
const RCON_PORT_RANGE = { min: 26000, max: 26999 };

export class CraftyServerOrchestrator
	implements
		ServerOrchestrator,
		FileCapability,
		BackupCapability,
		UpdateCapability
{
	readonly capabilities: ReadonlySet<OrchestratorCapability> = new Set([
		"files",
		"backup",
		"update",
	]);
	private readonly limit: ReturnType<typeof pLimit>;

	constructor(
		private readonly api: CraftyApiClient,
		private readonly config: CraftyConfig,
		private readonly rconRepository: CraftyRconRepository,
		private readonly rconClient: RconClient,
	) {
		this.limit = pLimit(config.maxConcurrentRequests);
	}

	supports(type: GameTypeId): boolean {
		return type === "java" || type === "bedrock" || type === "hytale";
	}

	portsFor(type: GameTypeId): OrchestratorPortRequest[] {
		return RCON_GAMES.has(type)
			? [{ purpose: "rcon", protocol: "tcp", defaultRange: RCON_PORT_RANGE }]
			: [];
	}

	propertiesFor(type: GameTypeId, ports: ServerPorts): GameProperties {
		if (!RCON_GAMES.has(type)) return {};
		if (ports.rcon === undefined)
			throw new InvalidServerPorts("Crafty needs an RCON port for this server");

		return {
			"enable-rcon": true,
			"rcon.port": ports.rcon,
			"rcon.password": randomBytes(24).toString("base64url"),
			"broadcast-rcon-to-ops": false,
		};
	}

	async create(spec: CreateServerSpec): Promise<OrchestratorRef> {
		const result = await this.api.createServer(
			toCraftyCreateServerRequest(spec, this.config.serversHost),
		);
		const ref = result.data.new_server_id;

		try {
			const password = spec.orchestratorProperties["rcon.password"];
			if (spec.ports.rcon !== undefined && typeof password === "string")
				this.rconRepository.set(ref, { port: spec.ports.rcon, password });

			if (this.config.webhooks.enabled) await this.createWebhooks(ref);
		} catch (error) {
			await this.api.deleteServer(ref, true).catch(() => undefined);
			this.rconRepository.delete(ref);
			throw new OrchestratorActionError(
				`Crafty failed to set up server ${spec.name}`,
				{
					cause: error,
				},
			);
		}

		return ref;
	}

	async delete(
		target: OrchestratorTarget,
		removeFiles: boolean,
	): Promise<void> {
		const result = await this.api.deleteServer(target.ref, removeFiles);
		this.throwIfFailed(result, `delete server ${target.ref}`);
		this.rconRepository.delete(target.ref);
	}

	async stats(target: OrchestratorTarget): Promise<GameServerStats> {
		const result = await this.limit(() => this.api.getServerStats(target.ref));
		return toGameServerStats(result.data);
	}

	async logs(target: OrchestratorTarget): Promise<string[]> {
		const result = await this.api.getServerLogs(target.ref);
		return result.data.map(fromCraftyLogLine);
	}

	start(target: OrchestratorTarget): Promise<void> {
		return this.runAction(target, "start_server");
	}

	stop(target: OrchestratorTarget): Promise<void> {
		return this.runAction(target, "stop_server");
	}

	restart(target: OrchestratorTarget): Promise<void> {
		return this.runAction(target, "restart_server");
	}

	kill(target: OrchestratorTarget): Promise<void> {
		return this.runAction(target, "kill_server");
	}

	acceptEula(target: OrchestratorTarget): Promise<void> {
		return this.runAction(target, "eula");
	}

	backup(target: OrchestratorTarget): Promise<void> {
		return this.runAction(target, "backup_server");
	}

	update(target: OrchestratorTarget): Promise<void> {
		return this.runAction(target, "update_executable");
	}

	async sendCommand(
		target: OrchestratorTarget,
		command: string,
	): Promise<CommandOutput> {
		const rcon = this.rconRepository.get(target.ref);
		if (rcon !== undefined) {
			const output = await this.rconClient.send(
				{
					host: this.config.serversHost,
					port: rcon.port,
					password: rcon.password,
				},
				command,
			);
			return { output };
		}

		const result = await this.api.sendCommand(target.ref, command);
		this.throwIfFailed(result, `send command to server ${target.ref}`);
		return {};
	}

	async stopAll(): Promise<number> {
		const result = await this.api.stopAllServers();
		return result.data.count;
	}

	async listFiles(
		target: OrchestratorTarget,
		path: string,
	): Promise<GameServerFileEntry[]> {
		const result = await this.api.getServerFiles(target.ref, path);
		if (!isCraftyDirectory(result.data))
			throw new OrchestratorActionError(`${path} is not a directory`);
		return toGameServerFileEntries(result.data);
	}

	async readFile(
		target: OrchestratorTarget,
		path: string,
	): Promise<GameServerFile> {
		const result = await this.api.getServerFiles(target.ref, path);
		if (isCraftyDirectory(result.data))
			throw new OrchestratorActionError(`${path} is a directory, not a file`);
		return toGameServerFile(path, result.data);
	}

	writeFile(target: OrchestratorTarget, file: GameServerFile): Promise<void> {
		return this.api.updateServerFile(
			target.ref,
			toCraftyFilePatchRequest(file),
		);
	}

	async createFile(target: OrchestratorTarget, path: string): Promise<void> {
		const result = await this.api.createServerFile(
			target.ref,
			toCraftyFileCreateRequest(path, false),
		);
		this.throwIfFailed(result, `create file ${path} on server ${target.ref}`);
	}

	async createDirectory(
		target: OrchestratorTarget,
		path: string,
	): Promise<void> {
		const result = await this.api.createServerFile(
			target.ref,
			toCraftyFileCreateRequest(path, true),
		);
		this.throwIfFailed(
			result,
			`create directory ${path} on server ${target.ref}`,
		);
	}

	async renameFile(
		target: OrchestratorTarget,
		path: string,
		newName: string,
	): Promise<void> {
		const result = await this.api.renameServerFile(
			target.ref,
			toCraftyFileRenameRequest(path, newName),
		);
		this.throwIfFailed(result, `rename ${path} on server ${target.ref}`);
	}

	async moveFiles(
		target: OrchestratorTarget,
		sourcePaths: string[],
		targetDirectory: string,
	): Promise<void> {
		const result = await this.api.runServerFileOperation(
			target.ref,
			"move",
			toCraftyFileOperationRequest(sourcePaths, targetDirectory),
		);
		this.throwIfFailed(result, `move files on server ${target.ref}`);
	}

	async copyFiles(
		target: OrchestratorTarget,
		sourcePaths: string[],
		targetDirectory: string,
	): Promise<void> {
		const result = await this.api.runServerFileOperation(
			target.ref,
			"copy",
			toCraftyFileOperationRequest(sourcePaths, targetDirectory),
		);
		this.throwIfFailed(result, `copy files on server ${target.ref}`);
	}

	async deleteFiles(
		target: OrchestratorTarget,
		paths: string[],
	): Promise<void> {
		const result = await this.api.deleteServerFiles(
			target.ref,
			toCraftyFileDeleteRequest(paths),
		);
		this.throwIfFailed(result, `delete files on server ${target.ref}`);
	}

	downloadFile(
		target: OrchestratorTarget,
		path: string,
	): Promise<GameServerFileDownload> {
		return this.api.downloadServerFile(target.ref, path);
	}

	uploadFile(
		target: OrchestratorTarget,
		path: string,
		file: BunFile,
	): Promise<void> {
		return this.api.uploadServerFile(target.ref, path, file);
	}

	async readConfigFile(
		target: OrchestratorTarget,
		path: string,
	): Promise<string | undefined> {
		const result = await this.api
			.getServerFiles(target.ref, path)
			.catch(() => undefined);
		if (result === undefined) return undefined;
		if (isCraftyDirectory(result.data))
			throw new OrchestratorActionError(`${path} is a directory, not a file`);
		return result.data.content;
	}

	async writeConfigFile(
		target: OrchestratorTarget,
		path: string,
		content: string,
	): Promise<void> {
		if ((await this.readConfigFile(target, path)) === undefined) {
			const created = await this.api.createServerFile(
				target.ref,
				toCraftyFileCreateRequest(path, false),
			);
			this.throwIfFailed(created, `create ${path} on server ${target.ref}`);
		}
		await this.api.updateServerFile(target.ref, { path, contents: content });
	}

	private async createWebhooks(ref: OrchestratorRef): Promise<void> {
		const webhooks = Object.entries(this.config.webhooks.urls) as [
			CraftyWebhookEvent,
			string,
		][];
		for (const [event, url] of webhooks)
			await this.api.createWebhook(
				ref,
				toCraftyCreateWebhookRequest(event, url),
			);
	}

	private async runAction(
		target: OrchestratorTarget,
		action: CraftyAction,
	): Promise<void> {
		const result = await this.api.sendAction(target.ref, action);
		this.throwIfFailed(result, `run ${action} on server ${target.ref}`);
	}

	private throwIfFailed(
		result: CraftyActionExecutionResponse,
		operation: string,
	): void {
		if (result.status !== "error") return;
		const details = result.error_data ? ` (${result.error_data})` : "";
		throw new OrchestratorActionError(
			`Crafty failed to ${operation}: ${result.error ?? "unknown error"}${details}`,
		);
	}
}
