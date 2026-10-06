import { PropertiesEditor } from "properties-file/editor";
import { MinecraftServerActionError } from "../../errors";
import type { CreateMinecraftServerRequest } from "../../types/create-minecraft-server-request";
import type { MinecraftServer } from "../../types/minecraft-server";
import type { MinecraftServerStats } from "../../types/minecraft-server-stats";
import type {
	MinecraftServerWebhookTrigger,
	MinecraftServerWebhookUrls,
} from "../../types/minecraft-server-webhook-trigger";
import { CRAFTY_PROVIDER_NAME } from "../../types/provider-name";
import type { MinecraftServerProvider } from "../minecraft-server-provider";
import type { RconConfig } from "../rcon-config";
import type { CraftyApiClient } from "./crafty-api-client";
import type { CraftyAction } from "./dtos/crafty-action";
import type { CraftyActionExecutionResponse } from "./dtos/crafty-action-execution-response";
import {
	fromCraftyLogLine,
	isSupportedCraftyServerType,
	toCraftyCreateServerRequest,
	toCraftyCreateWebhookRequest,
	toMinecraftServer,
	toMinecraftServerStats,
} from "./mappers";

export class CraftyMinecraftServerProvider implements MinecraftServerProvider {
	readonly name = CRAFTY_PROVIDER_NAME;

	constructor(
		private readonly craftyApiClient: CraftyApiClient,
		private readonly serversHost: string,
		private readonly webhookUrls: MinecraftServerWebhookUrls = {},
		private readonly rconConfig: RconConfig,
	) {}

	async create(
		createServerRequest: CreateMinecraftServerRequest,
	): Promise<string> {
		const craftyCreateServerRequest = toCraftyCreateServerRequest(
			createServerRequest,
			this.serversHost,
		);
		const result = await this.craftyApiClient.createServer(
			craftyCreateServerRequest,
		);
		const serverId = result.data.new_server_id;

		const serverPropertiesFile = await this.craftyApiClient.getServerFiles(
			serverId,
			"server.properties",
		);

		const propertiesEditor = new PropertiesEditor(
			serverPropertiesFile.data.content,
		);

		propertiesEditor.insert("enable-rcon", "true");
		propertiesEditor.insert("rcon.password", this.rconConfig.password);
		propertiesEditor.insert("rcon.port", this.rconConfig.port);
		propertiesEditor.insert("broadcast-rcon-to-ops", "false");

		// User configurable properties
		const createMSReqProps = createServerRequest.properties;
		if (createMSReqProps) {
			for (const [propertie, value] of Object.entries(createMSReqProps)) {
				propertiesEditor.insert(propertie, String(value));
			}
		}

		const newProperties = propertiesEditor.format({
			removeComments: true, // Strip all comments
			removeBlankLines: true, // Strip all blank lines
			removeLeadingWhitespace: true, // Strip indentation
			deduplicateKeys: true, // Keep only last occurrence
			separatorChar: "=", // Standardize separator
			separatorLeading: "", // Space before =
			separatorTrailing: "", // Space after =
		});

		this.craftyApiClient.updateServerFile(serverId, {
			contents: newProperties,
			path: "server.properties",
			modified_epoch: serverPropertiesFile.data.attributes.modified_epoch,
		});

		await this.createWebhooks(serverId);
		return serverId;
	}

	async list(): Promise<MinecraftServer[]> {
		const servers = await this.craftyApiClient.listServers();

		return servers.data
			.filter((dto) => isSupportedCraftyServerType(dto.type))
			.map((dto) => toMinecraftServer(dto, this.serversHost));
	}

	async info(id: string): Promise<MinecraftServer> {
		const [server, stats] = await Promise.all([
			this.craftyApiClient.getServer(id),
			this.craftyApiClient.getServerStats(id),
		]);
		return toMinecraftServer(server.data, this.serversHost, stats.data);
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

	private async createWebhooks(serverId: string): Promise<void> {
		const webhooks = Object.entries(this.webhookUrls) as [
			MinecraftServerWebhookTrigger,
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
			throw new MinecraftServerActionError(
				`Crafty failed to ${operation}: ${result.error ?? "unknown error"}${details}`,
			);
		}
	}
}
