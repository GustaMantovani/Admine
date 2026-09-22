import {
	DEFAULT_BEDROCK_PORT,
	DEFAULT_HYTALE_PORT,
	DEFAULT_JAVA_PORT,
	DEFAULT_MAXIMUM_MEMORY,
	DEFAULT_MINECRAFT_SERVER_VERSION,
	DEFAULT_MINIMUM_MEMORY,
	SERVERS_HOST_FOR_ADMINE,
} from "../../constants";
import { UnknownMinecraftServerType } from "../../errors";
import type { CreateMinecraftServerRequest } from "../../types/create-minecraft-server-request";
import { MinecraftServer } from "../../types/minecraft-server";
import type { MinecraftServerStats } from "../../types/minecraft-server-stats";
import { MinecraftServerType } from "../../types/minecraft-server-type";
import type { CraftyCreateBedrockServer } from "./dtos/crafty-create-bedrock-server";
import type { CraftyCreateHytaleServer } from "./dtos/crafty-create-hytale-server";
import type { CraftyCreateJavaServer } from "./dtos/crafty-create-java-server";
import type { CraftyCreateServer } from "./dtos/crafty-create-server";
import type { CraftyServer } from "./dtos/crafty-server";
import type { CraftyServerStats } from "./dtos/crafty-server-stats";

export function resolvePort(request: CreateMinecraftServerRequest): number {
	if (request.port != null) {
		return request.port;
	}

	switch (request.type) {
		case "java":
			return DEFAULT_JAVA_PORT;
		case "bedrock":
			return DEFAULT_BEDROCK_PORT;
		case "hytale":
			return DEFAULT_HYTALE_PORT;
		default:
			return DEFAULT_JAVA_PORT;
	}
}

export function toCraftyCreateServerRequest(
	request: CreateMinecraftServerRequest,
): CraftyCreateServer {
	switch (request.type) {
		case "java":
			return toCraftyCreateJavaServer(request);
		case "bedrock":
			return toCraftyCreateBedrockServer(request);
		case "hytale":
			return toCraftyCreateHytaleServer(request);
		default:
			return toCraftyCreateJavaServer(request);
	}
}

function toCraftyCreateJavaServer(
	request: CreateMinecraftServerRequest,
): CraftyCreateJavaServer {
	const port = resolvePort(request);

	return {
		name: request.name,
		monitoring_type: "minecraft_java",
		minecraft_java_monitoring_data: {
			host: SERVERS_HOST_FOR_ADMINE,
			port,
		},
		create_type: "minecraft_java",
		minecraft_java_create_data: {
			create_type: "download_jar",
			download_jar_create_data: {
				category: "mc_java_servers",
				type: "paper",
				version: request.version ?? DEFAULT_MINECRAFT_SERVER_VERSION,
				mem_min: request.minimumMemory ?? DEFAULT_MINIMUM_MEMORY,
				mem_max: request.maximumMemory ?? DEFAULT_MAXIMUM_MEMORY,
				server_properties_port: port,
			},
		},
	};
}

function toCraftyCreateBedrockServer(
	request: CreateMinecraftServerRequest,
): CraftyCreateBedrockServer {
	return {
		name: request.name,
		monitoring_type: "minecraft_bedrock",
		minecraft_bedrock_monitoring_data: {
			host: SERVERS_HOST_FOR_ADMINE,
			port: resolvePort(request),
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
): CraftyCreateHytaleServer {
	return {
		name: request.name,
		monitoring_type: "hytale",
		hytale_monitoring_data: {
			host: SERVERS_HOST_FOR_ADMINE,
			port: resolvePort(request),
		},
		create_type: "hytale",
		hytale_create_data: {
			create_type: "download_exe",
			download_exe_create_data: {
				agree_to_eula: true,
				mem_min: request.minimumMemory ?? DEFAULT_MINIMUM_MEMORY,
				mem_max: request.maximumMemory ?? DEFAULT_MAXIMUM_MEMORY,
			},
		},
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

function mapServerTypesToEnum(dto: string): MinecraftServerType {
	const mapped = CRAFTY_SERVER_TYPES[dto];

	if (mapped === undefined)
		throw new UnknownMinecraftServerType(
			`Cant map server type from string ${dto}`,
		);

	return mapped;
}

export function toMinecraftServerStatus(
	running: boolean | undefined,
): MinecraftServer["status"] {
	if (running === undefined) return "unknow";
	return running ? "on" : "off";
}

export function toMinecraftServer(
	dto: CraftyServer,
	running: boolean | undefined,
): MinecraftServer {
	return new MinecraftServer(
		dto.server_id,
		dto.server_name,
		mapServerTypesToEnum(dto.type),
		toMinecraftServerStatus(running),
		undefined,
		SERVERS_HOST_FOR_ADMINE,
		[],
		dto.server_port,
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
