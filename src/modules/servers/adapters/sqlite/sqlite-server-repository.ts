import type { Database } from "bun:sqlite";
import type { ServerRecord } from "../../domain/game-server";
import type { GameTypeId } from "../../domain/game-types/game-type";
import type { ServerRepository } from "../../ports/server-repository";

type Row = {
	id: string;
	name: string;
	type: string;
	version: string | null;
	min_memory: number | null;
	max_memory: number | null;
	orchestrator_ref: string;
	game_port: number;
	created_at: string;
};

export class SqliteServerRepository implements ServerRepository {
	constructor(private readonly db: Database) {}

	async insert(record: ServerRecord): Promise<void> {
		this.db
			.query(
				`INSERT INTO servers_records
					(id, name, type, version, min_memory, max_memory, orchestrator_ref, game_port, created_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			)
			.run(
				record.id,
				record.name,
				record.type,
				record.version ?? null,
				record.memory?.minimum ?? null,
				record.memory?.maximum ?? null,
				record.orchestratorRef,
				record.gamePort,
				record.createdAt,
			);
	}

	async findById(id: string): Promise<ServerRecord | undefined> {
		const row = this.db
			.query("SELECT * FROM servers_records WHERE id = ?")
			.get(id) as Row | null;
		return row ? toRecord(row) : undefined;
	}

	async findByName(name: string): Promise<ServerRecord | undefined> {
		const row = this.db
			.query("SELECT * FROM servers_records WHERE name = ?")
			.get(name) as Row | null;
		return row ? toRecord(row) : undefined;
	}

	async list(): Promise<ServerRecord[]> {
		const rows = this.db
			.query("SELECT * FROM servers_records ORDER BY created_at")
			.all() as Row[];
		return rows.map(toRecord);
	}

	async delete(id: string): Promise<void> {
		this.db.query("DELETE FROM servers_records WHERE id = ?").run(id);
	}
}

function toRecord(row: Row): ServerRecord {
	return {
		id: row.id,
		name: row.name,
		type: row.type as GameTypeId,
		version: row.version ?? undefined,
		memory:
			row.min_memory !== null && row.max_memory !== null
				? { minimum: row.min_memory, maximum: row.max_memory }
				: undefined,
		orchestratorRef: row.orchestrator_ref,
		gamePort: row.game_port,
		createdAt: row.created_at,
	};
}
