import type { ServerRecord } from "../domain/game-server";

export interface ServerRepository {
	insert(record: ServerRecord): Promise<void>;
	findById(id: string): Promise<ServerRecord | undefined>;
	findByName(name: string): Promise<ServerRecord | undefined>;
	list(): Promise<ServerRecord[]>;
	delete(id: string): Promise<void>;
}
