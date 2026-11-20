import type { ServerMemory } from "./create-server-spec";
import type { GameTypeId } from "./game-types/game-type";
import type { PortPurpose } from "./port-purpose";

export type ServerRecord = {
	id: string;
	name: string;
	type: GameTypeId;
	version?: string;
	memory?: ServerMemory;
	orchestratorRef: string;
	gamePort: number;
	createdAt: string;
};

export type ServerStatus = "on" | "off" | "unknown";

export type ServerOverview = ServerRecord & {
	status: ServerStatus;
	runningVersion?: string;
};

export type ServerDetails = ServerOverview & {
	ports: { port: number; purpose: PortPurpose }[];
};

export type GameServerStats = {
	running: boolean;
	crashed: boolean;
	startedAt?: string;
	cpu: number;
	memoryBytes: number;
	memoryPercent: number;
	playersOnline: number;
	maxPlayers: number;
	version?: string;
	worldName: string;
	worldSize: string;
};

export type CommandOutput = {
	output?: string;
};

export type GameServerFile = {
	path: string;
	content: string;
	modifiedEpoch?: number;
};

export type GameServerFileEntry = {
	name: string;
	path: string;
	type: "file" | "directory";
};

export type GameServerFileDownload = {
	contentStream: ReadableStream<Uint8Array>;
	fileNameOnServer?: string;
};
