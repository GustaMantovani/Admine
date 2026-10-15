import type { Database } from "bun:sqlite";

export interface CraftyRconPortRepository {
	get(serverId: string): Promise<number | undefined>;
	getAll(): Promise<Map<string, number>>;
	set(serverId: string, rconPort: number): Promise<void>;
	delete(serverId: string): Promise<void>;
}

export class SqliteCraftyRconPortRepository
	implements CraftyRconPortRepository
{
	constructor(private readonly db: Database) {}

	async get(serverId: string): Promise<number | undefined> {
		const row = this.db
			.query("SELECT rconPort FROM CraftyRconPort WHERE serverId = ?")
			.get(serverId) as { rconPort: number } | null;
		return row?.rconPort;
	}

	async getAll(): Promise<Map<string, number>> {
		const rows = this.db
			.query("SELECT serverId, rconPort FROM CraftyRconPort")
			.all() as { serverId: string; rconPort: number }[];
		return new Map(rows.map((r) => [r.serverId, r.rconPort]));
	}

	async set(serverId: string, rconPort: number): Promise<void> {
		this.db
			.query(
				"INSERT INTO CraftyRconPort (serverId, rconPort) VALUES (?, ?) ON CONFLICT (serverId) DO UPDATE SET rconPort = excluded.rconPort",
			)
			.run(serverId, rconPort);
	}

	async delete(serverId: string): Promise<void> {
		this.db
			.query("DELETE FROM CraftyRconPort WHERE serverId = ?")
			.run(serverId);
	}
}
