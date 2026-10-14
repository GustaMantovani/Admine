import type { Database } from "bun:sqlite";
import type { GameServerPorts } from "../domain/game-server-ports";

export interface GameServerPortsRepository {
	create(gameServer: GameServerPorts): Promise<void>;
	delete(id: string): Promise<void>;
	getUsedGamePorts(): Promise<number[]>;
	getUsedRconPorts(): Promise<number[]>;
}

export class SqliteGameServerPortsRepository
	implements GameServerPortsRepository
{
	constructor(private readonly db: Database) {}

	async create(gameServer: GameServerPorts): Promise<void> {
		this.db
			.query(
				"INSERT INTO GameServerPorts (id, gamePort, rconPort) VALUES (?, ?, ?)",
			)
			.run(gameServer.id, gameServer.port, gameServer.rconPort ?? null);
	}

	async delete(id: string): Promise<void> {
		this.db.query("DELETE FROM GameServerPorts WHERE id = ?").run(id);
	}

	async getUsedGamePorts(): Promise<number[]> {
		const rows = this.db
			.query("SELECT DISTINCT gamePort FROM GameServerPorts")
			.all() as { gamePort: number }[];

		return rows.map((r) => r.gamePort);
	}

	async getUsedRconPorts(): Promise<number[]> {
		const rows = this.db
			.query(
				"SELECT DISTINCT rconPort FROM GameServerPorts WHERE rconPort IS NOT NULL",
			)
			.all() as { rconPort: number }[];

		return rows.map((r) => r.rconPort);
	}
}
