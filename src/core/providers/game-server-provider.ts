import type { CreateGameServerRequest } from "../domain/create-game-server-request";
import type { GameServer } from "../domain/game-server";
import type { GameServerFile } from "../domain/game-server-file";
import type { GameServerFileEntry } from "../domain/game-server-file-entry";
import type { GameServerStats } from "../domain/game-server-stats";
import type { ProviderName } from "./provider-name";

export interface GameServerProvider {
	readonly name: ProviderName;
	create(createServerRequest: CreateGameServerRequest): Promise<string>;
	list(): Promise<GameServer[]>;
	info(id: string): Promise<GameServer>;
	stats(id: string): Promise<GameServerStats>;
	logs(id: string): Promise<string[]>;
	start(id: string): Promise<void>;
	stop(id: string): Promise<void>;
	restart(id: string): Promise<void>;
	kill(id: string): Promise<void>;
	backup(id: string): Promise<void>;
	update(id: string): Promise<void>;
	clone(id: string): Promise<void>;
	acceptEula(id: string): Promise<void>;
	sendCommand(id: string, command: string): Promise<void>;
	// File paths are relative to the server's root directory
	listFiles(id: string, path: string): Promise<GameServerFileEntry[]>;
	getFile(id: string, path: string): Promise<GameServerFile>;
	writeFile(id: string, file: GameServerFile): Promise<void>;
	createFile(id: string, path: string): Promise<void>;
	createDirectory(id: string, path: string): Promise<void>;
	renameFile(id: string, path: string, newName: string): Promise<void>;
	// Each item keeps its name inside targetDirectory
	moveFiles(
		id: string,
		sourcePaths: string[],
		targetDirectory: string,
	): Promise<void>;
	copyFiles(
		id: string,
		sourcePaths: string[],
		targetDirectory: string,
	): Promise<void>;
	deleteFiles(id: string, paths: string[]): Promise<void>;
	delete(id: string, removeFiles: boolean): Promise<void>;
	stopAll(): Promise<number>;
}
