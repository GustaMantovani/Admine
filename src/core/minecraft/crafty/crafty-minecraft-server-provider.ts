import { $ } from "bun";
import { MinecraftServerActionError } from "../../errors";
import type { CreateMinecraftServerRequest } from "../../types/create-minecraft-server-request";
import type { MinecraftServer } from "../../types/minecraft-server";
import type { MinecraftServerStats } from "../../types/minecraft-server-stats";
import type { MinecraftServerProvider } from "../minecraft-server-provider";
import type { CraftyApiClient } from "./crafty-api-client";
import type { CraftyInfraConfig } from "./crafty-config";
import type { CraftyAction } from "./dtos/crafty-action";
import type { CraftyActionExecutionResult } from "./dtos/crafty-action-execution-result";
import {
	fromCraftyLogLine,
	isSupportedCraftyServerType,
	toCraftyCreateServerRequest,
	toMinecraftServer,
	toMinecraftServerStats,
} from "./mappers";

export class CraftyMinecraftServerProvider implements MinecraftServerProvider {
	constructor(
		private readonly craftyApiClient: CraftyApiClient,
		private readonly infraConfig?: CraftyInfraConfig,
	) {}

	async setup() {
		// Infra setup
		if (!this.infraConfig) return;

		const firstInfraCheckExitCode = (
			await $`${{ raw: this.infraConfig.infraCheckCommand }}`.nothrow().quiet()
		).exitCode;
		if (firstInfraCheckExitCode === 0) return;

		await $`${{ raw: this.infraConfig.infraSetupCommand }}`;

		Bun.sleepSync(this.infraConfig.infraCheckDelay);

		await $`${{ raw: this.infraConfig.infraCheckCommand }}`;
	}

	async create(
		createServerRequest: CreateMinecraftServerRequest,
	): Promise<string> {
		const craftyCreateServerRequest =
			toCraftyCreateServerRequest(createServerRequest);
		const result = await this.craftyApiClient.createServer(
			craftyCreateServerRequest,
		);
		return result.data.new_server_id;
	}

	async list(): Promise<MinecraftServer[]> {
		const servers = await this.craftyApiClient.listServers();

		return servers.data
			.filter((dto) => isSupportedCraftyServerType(dto.type))
			.map((dto) => toMinecraftServer(dto));
	}

	async info(id: string): Promise<MinecraftServer> {
		const [server, stats] = await Promise.all([
			this.craftyApiClient.getServer(id),
			this.craftyApiClient.getServerStats(id),
		]);
		return toMinecraftServer(server.data, stats.data);
	}

	async stats(id: string): Promise<MinecraftServerStats> {
		const result = await this.craftyApiClient.getServerStats(id);
		return toMinecraftServerStats(result.data);
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

	async delete(id: string, removeFiles: boolean): Promise<void> {
		const result = await this.craftyApiClient.deleteServer(id, removeFiles);
		this.throwIfFailed(result, `delete server ${id}`);
	}

	async stopAll(): Promise<number> {
		const result = await this.craftyApiClient.stopAllServers();
		return result.data.count;
	}

	private async runAction(id: string, action: CraftyAction): Promise<void> {
		const result = await this.craftyApiClient.sendAction(id, action);
		this.throwIfFailed(result, `run ${action} on server ${id}`);
	}

	// Some Crafty failures come back as HTTP 200 with status "error"
	private throwIfFailed(
		result: CraftyActionExecutionResult,
		operation: string,
	): void {
		if (result.status === "error") {
			const details = result.error_data ? ` (${result.error_data})` : "";
			throw new MinecraftServerActionError(
				`Crafty failed to ${operation}: ${result.error ?? "unknown error"}${details}`,
			);
		}
	}
}
