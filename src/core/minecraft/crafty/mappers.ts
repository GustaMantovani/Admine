import { UnknownMinecraftServerType } from "../../errors";
import type { CreateMinecraftServerRequest } from "../../types/create-minecraft-server-request";
import { MinecraftServer } from "../../types/minecraft-server";
import type { MinecraftServerStats } from "../../types/minecraft-server-stats";
import { MinecraftServerType } from "../../types/minecraft-server-type";
import { MinecraftServerWebhookTrigger } from "../../types/minecraft-server-webhook-trigger";
import { resolveServerPort } from "../server-ports";
import type { CraftyCreateBedrockServer } from "./dtos/crafty-create-bedrock-server";
import type { CraftyCreateHytaleServer } from "./dtos/crafty-create-hytale-server";
import type { CraftyCreateJavaServer } from "./dtos/crafty-create-java-server";
import type { CraftyCreateServer } from "./dtos/crafty-create-server";
import type { CraftyCreateWebhook } from "./dtos/crafty-create-webhook";
import type { CraftyServer } from "./dtos/crafty-server";
import type { CraftyServerStats } from "./dtos/crafty-server-stats";
import type { CraftyWebhookTrigger } from "./dtos/crafty-webhook-trigger";

export function toCraftyCreateServerRequest(
	request: CreateMinecraftServerRequest,
	serversHost: string,
): CraftyCreateServer {
	switch (request.type) {
		case "java":
			return toCraftyCreateJavaServer(request, serversHost);
		case "bedrock":
			return toCraftyCreateBedrockServer(request, serversHost);
		case "hytale":
			return toCraftyCreateHytaleServer(request, serversHost);
		default:
			return toCraftyCreateJavaServer(request, serversHost);
	}
}

function toCraftyCreateJavaServer(
	request: CreateMinecraftServerRequest,
	serversHost: string,
): CraftyCreateJavaServer {
	const port = resolveServerPort(request);

	return {
		name: request.name,
		monitoring_type: "minecraft_java",
		minecraft_java_monitoring_data: {
			host: serversHost,
			port,
		},
		create_type: "minecraft_java",
		minecraft_java_create_data: {
			create_type: "download_jar",
			download_jar_create_data: {
				category: "mc_java_servers",
				type: "paper",
				version: request.version,
				mem_min: request.minimumMemory,
				mem_max: request.maximumMemory,
				server_properties_port: port,
			},
		},
	};
}

function toCraftyCreateBedrockServer(
	request: CreateMinecraftServerRequest,
	serversHost: string,
): CraftyCreateBedrockServer {
	return {
		name: request.name,
		monitoring_type: "minecraft_bedrock",
		minecraft_bedrock_monitoring_data: {
			host: serversHost,
			port: resolveServerPort(request),
		},
		create_type: "minecraft_bedrock",
		minecraft_bedrock_create_data: {
			create_type: "download_exe",
			download_exe_create_data: {
				agree_to_eula: true,
			},
		},
	};
}

function toCraftyCreateHytaleServer(
	request: CreateMinecraftServerRequest,
	serversHost: string,
): CraftyCreateHytaleServer {
	return {
		name: request.name,
		monitoring_type: "hytale",
		hytale_monitoring_data: {
			host: serversHost,
			port: resolveServerPort(request),
		},
		create_type: "hytale",
		hytale_create_data: {
			create_type: "download_exe",
			download_exe_create_data: {
				agree_to_eula: true,
				mem_min: request.minimumMemory,
				mem_max: request.maximumMemory,
			},
		},
	};
}

const CRAFTY_WEBHOOKS: Record<
	MinecraftServerWebhookTrigger,
	{ trigger: CraftyWebhookTrigger; name: string }
> = {
	[MinecraftServerWebhookTrigger.Started]: {
		trigger: "start_server",
		name: "Server Started",
	},
	[MinecraftServerWebhookTrigger.Stopped]: {
		trigger: "stop_server",
		name: "Server Stopped",
	},
	[MinecraftServerWebhookTrigger.Crashed]: {
		trigger: "crash_detected",
		name: "Server Crashed",
	},
	[MinecraftServerWebhookTrigger.BackupCompleted]: {
		trigger: "backup_server",
		name: "Server Backup Completed",
	},
	[MinecraftServerWebhookTrigger.ExecutableUpdated]: {
		trigger: "jar_update",
		name: "Server Executable Updated",
	},
	[MinecraftServerWebhookTrigger.CommandReceived]: {
		trigger: "send_command",
		name: "Server Command Received",
	},
	[MinecraftServerWebhookTrigger.Killed]: {
		trigger: "kill",
		name: "Server Killed",
	},
};

export function toCraftyCreateWebhookRequest(
	trigger: MinecraftServerWebhookTrigger,
	url: string,
): CraftyCreateWebhook {
	const { trigger: craftyTrigger, name } = CRAFTY_WEBHOOKS[trigger];

	return {
		// Crafty has no generic webhook type, so Discord's payload is used
		webhook_type: "Discord",
		name,
		url,
		bot_name: "Admine",
		trigger: [craftyTrigger],
		body: "{{ event_type }} | {{ server_name }} | {{ server_id }} | {{ time_iso }}",
		color: "#005cd1",
		enabled: true,
	};
}

// Crafty also manages types Admine doesn't support yet, such as "steam_cmd"
const CRAFTY_SERVER_TYPES: Record<string, MinecraftServerType> = {
	"minecraft-java": MinecraftServerType.Java,
	"minecraft-bedrock": MinecraftServerType.Bedrock,
	hytale: MinecraftServerType.Hytale,
};

export function isSupportedCraftyServerType(type: string): boolean {
	return type in CRAFTY_SERVER_TYPES;
}

function fromCraftyServerType(dto: string): MinecraftServerType {
	const mapped = CRAFTY_SERVER_TYPES[dto];

	if (mapped === undefined)
		throw new UnknownMinecraftServerType(
			`Can't map server type from string ${dto}`,
		);

	return mapped;
}

export function toMinecraftServerStatus(
	running: boolean | undefined,
): MinecraftServer["status"] {
	if (running === undefined) return "unknown";
	return running ? "on" : "off";
}

export function toMinecraftServer(
	craftyServerDto: CraftyServer,
	serversHost: string,
	craftyServerStatsDto?: CraftyServerStats | undefined,
): MinecraftServer {
	return new MinecraftServer(
		craftyServerDto.server_id,
		craftyServerDto.server_name,
		fromCraftyServerType(craftyServerDto.type),
		toMinecraftServerStatus(craftyServerStatsDto?.running),
		craftyServerStatsDto?.version ?? undefined,
		serversHost,
		[],
		craftyServerDto.server_port,
		undefined,
		undefined,
	);
}

// Crafty stores Python's str(False) in text stats while the server is offline
function fromCraftyText(value: string): string | undefined {
	return value === "" || value === "False" ? undefined : value;
}

export function toMinecraftServerStats(
	dto: CraftyServerStats,
): MinecraftServerStats {
	return {
		running: dto.running,
		crashed: dto.crashed,
		startedAt: fromCraftyText(dto.started),
		cpu: dto.cpu,
		memoryBytes: Number(dto.mem),
		memoryPercent: dto.mem_percent,
		playersOnline: dto.online,
		maxPlayers: dto.max,
		version: fromCraftyText(dto.version),
		worldName: dto.world_name,
		worldSize: dto.world_size,
	};
}

// Crafty runs each log line through Python's html.escape
export function fromCraftyLogLine(line: string): string {
	return line
		.replaceAll("&lt;", "<")
		.replaceAll("&gt;", ">")
		.replaceAll("&quot;", '"')
		.replaceAll("&#x27;", "'")
		.replaceAll("&amp;", "&");
}
