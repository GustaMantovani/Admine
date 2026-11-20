import type { BunFile } from "bun";
import type { CreateServerSpec } from "../domain/create-server-spec";
import type {
	CommandOutput,
	GameServerFile,
	GameServerFileDownload,
	GameServerFileEntry,
	GameServerStats,
} from "../domain/game-server";
import type {
	GameProperties,
	GameTypeId,
} from "../domain/game-types/game-type";
import type {
	PortProtocol,
	PortPurpose,
	PortRange,
	ServerPorts,
} from "../domain/port-purpose";

export type OrchestratorRef = string;

export type OrchestratorTarget = { ref: OrchestratorRef; type: GameTypeId };

export type OrchestratorCapability = "files" | "backup" | "update";

export type OrchestratorPortRequest = {
	purpose: PortPurpose;
	protocol: PortProtocol;
	defaultRange: PortRange;
};

export interface ServerOrchestrator {
	readonly capabilities: ReadonlySet<OrchestratorCapability>;
	supports(type: GameTypeId): boolean;
	portsFor(type: GameTypeId): OrchestratorPortRequest[];
	propertiesFor(type: GameTypeId, ports: ServerPorts): GameProperties;
	create(spec: CreateServerSpec): Promise<OrchestratorRef>;
	delete(target: OrchestratorTarget, removeFiles: boolean): Promise<void>;
	stats(target: OrchestratorTarget): Promise<GameServerStats>;
	logs(target: OrchestratorTarget): Promise<string[]>;
	start(target: OrchestratorTarget): Promise<void>;
	stop(target: OrchestratorTarget): Promise<void>;
	restart(target: OrchestratorTarget): Promise<void>;
	kill(target: OrchestratorTarget): Promise<void>;
	acceptEula(target: OrchestratorTarget): Promise<void>;
	sendCommand(
		target: OrchestratorTarget,
		command: string,
	): Promise<CommandOutput>;
	stopAll(): Promise<number>;
	readConfigFile(
		target: OrchestratorTarget,
		path: string,
	): Promise<string | undefined>;
	writeConfigFile(
		target: OrchestratorTarget,
		path: string,
		content: string,
	): Promise<void>;
}

export interface FileCapability {
	listFiles(
		target: OrchestratorTarget,
		path: string,
	): Promise<GameServerFileEntry[]>;
	readFile(target: OrchestratorTarget, path: string): Promise<GameServerFile>;
	writeFile(target: OrchestratorTarget, file: GameServerFile): Promise<void>;
	createFile(target: OrchestratorTarget, path: string): Promise<void>;
	createDirectory(target: OrchestratorTarget, path: string): Promise<void>;
	renameFile(
		target: OrchestratorTarget,
		path: string,
		newName: string,
	): Promise<void>;
	moveFiles(
		target: OrchestratorTarget,
		sourcePaths: string[],
		targetDirectory: string,
	): Promise<void>;
	copyFiles(
		target: OrchestratorTarget,
		sourcePaths: string[],
		targetDirectory: string,
	): Promise<void>;
	deleteFiles(target: OrchestratorTarget, paths: string[]): Promise<void>;
	downloadFile(
		target: OrchestratorTarget,
		path: string,
	): Promise<GameServerFileDownload>;
	uploadFile(
		target: OrchestratorTarget,
		path: string,
		file: BunFile,
	): Promise<void>;
}

export interface BackupCapability {
	backup(target: OrchestratorTarget): Promise<void>;
}

export interface UpdateCapability {
	update(target: OrchestratorTarget): Promise<void>;
}

export function hasFiles(
	orchestrator: ServerOrchestrator,
): orchestrator is ServerOrchestrator & FileCapability {
	return orchestrator.capabilities.has("files");
}

export function hasBackup(
	orchestrator: ServerOrchestrator,
): orchestrator is ServerOrchestrator & BackupCapability {
	return orchestrator.capabilities.has("backup");
}

export function hasUpdate(
	orchestrator: ServerOrchestrator,
): orchestrator is ServerOrchestrator & UpdateCapability {
	return orchestrator.capabilities.has("update");
}
