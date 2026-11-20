import { posix } from "node:path";
import type { CreateServerSpec } from "../../domain/create-server-spec";
import type {
	GameServerFile,
	GameServerFileEntry,
	GameServerStats,
} from "../../domain/game-server";
import type { GameTypeId } from "../../domain/game-types/game-type";
import type { CraftyWebhookEvent } from "./crafty-config";
import type {
	CraftyCreateServer,
	CraftyCreateWebhook,
	CraftyDirectory,
	CraftyFile,
	CraftyFileCreateRequest,
	CraftyFileDeleteRequest,
	CraftyFileOperationRequest,
	CraftyFilePatchRequest,
	CraftyFileRenameRequest,
	CraftyServerStats,
	CraftyWebhookTrigger,
} from "./dtos";

type SpecOf<T extends GameTypeId> = Extract<CreateServerSpec, { type: T }>;

const CREATE_MAPPERS: {
	[T in GameTypeId]: (
		spec: SpecOf<T>,
		serversHost: string,
	) => CraftyCreateServer;
} = {
	java: (spec, serversHost) => ({
		name: spec.name,
		monitoring_type: "minecraft_java",
		minecraft_java_monitoring_data: {
			host: serversHost,
			port: spec.ports.game,
		},
		create_type: "minecraft_java",
		minecraft_java_create_data: {
			create_type: "download_jar",
			download_jar_create_data: {
				category: "mc_java_servers",
				type: "paper",
				version: spec.version,
				mem_min: spec.memory.minimum,
				mem_max: spec.memory.maximum,
				server_properties_port: spec.ports.game,
			},
		},
	}),
	bedrock: (spec, serversHost) => ({
		name: spec.name,
		monitoring_type: "minecraft_bedrock",
		minecraft_bedrock_monitoring_data: {
			host: serversHost,
			port: spec.ports.game,
		},
		create_type: "minecraft_bedrock",
		minecraft_bedrock_create_data: {
			create_type: "download_exe",
			download_exe_create_data: { agree_to_eula: true },
		},
	}),
	hytale: (spec, serversHost) => ({
		name: spec.name,
		monitoring_type: "hytale",
		hytale_monitoring_data: { host: serversHost, port: spec.ports.game },
		create_type: "hytale",
		hytale_create_data: {
			create_type: "download_exe",
			download_exe_create_data: {
				agree_to_eula: true,
				mem_min: spec.memory.minimum,
				mem_max: spec.memory.maximum,
			},
		},
	}),
};

export function toCraftyCreateServerRequest(
	spec: CreateServerSpec,
	serversHost: string,
): CraftyCreateServer {
	const mapper = CREATE_MAPPERS[spec.type] as (
		spec: CreateServerSpec,
		serversHost: string,
	) => CraftyCreateServer;
	return mapper(spec, serversHost);
}

const CRAFTY_WEBHOOKS: Record<
	CraftyWebhookEvent,
	{ trigger: CraftyWebhookTrigger; name: string }
> = {
	started: { trigger: "start_server", name: "Server Started" },
	stopped: { trigger: "stop_server", name: "Server Stopped" },
	crashed: { trigger: "crash_detected", name: "Server Crashed" },
	backupCompleted: {
		trigger: "backup_server",
		name: "Server Backup Completed",
	},
	executableUpdated: {
		trigger: "jar_update",
		name: "Server Executable Updated",
	},
	commandReceived: { trigger: "send_command", name: "Server Command Received" },
	killed: { trigger: "kill", name: "Server Killed" },
};

export function toCraftyCreateWebhookRequest(
	event: CraftyWebhookEvent,
	url: string,
): CraftyCreateWebhook {
	const { trigger, name } = CRAFTY_WEBHOOKS[event];
	return {
		webhook_type: "Discord",
		name,
		url,
		bot_name: "Admine",
		trigger: [trigger],
		body: "{{ event_type }} | {{ server_name }} | {{ server_id }} | {{ time_iso }}",
		color: "#005cd1",
		enabled: true,
	};
}

function fromCraftyText(value: string): string | undefined {
	return value === "" || value === "False" ? undefined : value;
}

export function toGameServerStats(dto: CraftyServerStats): GameServerStats {
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

export function isCraftyDirectory(
	dto: CraftyFile | CraftyDirectory,
): dto is CraftyDirectory {
	return "root_path" in dto;
}

export function toGameServerFile(
	path: string,
	dto: CraftyFile,
): GameServerFile {
	return {
		path,
		content: dto.content,
		modifiedEpoch: dto.attributes.modified_epoch,
	};
}

export function toGameServerFileEntries(
	dto: CraftyDirectory,
): GameServerFileEntry[] {
	return Object.entries(dto).flatMap(([name, entry]) =>
		"dir" in entry
			? [{ name, path: entry.path, type: entry.dir ? "directory" : "file" }]
			: [],
	);
}

export function toCraftyFilePatchRequest(
	file: GameServerFile,
): CraftyFilePatchRequest {
	return {
		path: file.path,
		contents: file.content,
		modified_epoch: file.modifiedEpoch,
	};
}

export function toCraftyFileCreateRequest(
	path: string,
	directory: boolean,
): CraftyFileCreateRequest {
	return { parent: posix.dirname(path), name: posix.basename(path), directory };
}

export function toCraftyFileRenameRequest(
	path: string,
	newName: string,
): CraftyFileRenameRequest {
	return { path, new_name: newName };
}

export function toCraftyFileOperationRequest(
	sourcePaths: string[],
	targetDirectory: string,
): CraftyFileOperationRequest {
	return {
		file_system_objects: sourcePaths.map((sourcePath) => ({
			source_path: sourcePath,
			target_path: targetDirectory,
		})),
	};
}

export function toCraftyFileDeleteRequest(
	paths: string[],
): CraftyFileDeleteRequest {
	return { file_system_objects: paths.map((path) => ({ filename: path })) };
}

export function fromCraftyLogLine(line: string): string {
	return line
		.replaceAll("&lt;", "<")
		.replaceAll("&gt;", ">")
		.replaceAll("&quot;", '"')
		.replaceAll("&#x27;", "'")
		.replaceAll("&amp;", "&");
}
