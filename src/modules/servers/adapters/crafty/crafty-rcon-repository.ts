import type { Database } from "bun:sqlite";

export type CraftyRcon = { port: number; password: string };

export class CraftyRconRepository {
	constructor(private readonly db: Database) {}

	get(orchestratorRef: string): CraftyRcon | undefined {
		const row = this.db
			.query(
				"SELECT rcon_port, rcon_password FROM servers_crafty_rcon WHERE orchestrator_ref = ?",
			)
			.get(orchestratorRef) as {
			rcon_port: number;
			rcon_password: string;
		} | null;
		return row
			? { port: row.rcon_port, password: row.rcon_password }
			: undefined;
	}

	set(orchestratorRef: string, rcon: CraftyRcon): void {
		this.db
			.query(
				`INSERT INTO servers_crafty_rcon (orchestrator_ref, rcon_port, rcon_password) VALUES (?, ?, ?)
				ON CONFLICT (orchestrator_ref) DO UPDATE SET rcon_port = excluded.rcon_port, rcon_password = excluded.rcon_password`,
			)
			.run(orchestratorRef, rcon.port, rcon.password);
	}

	delete(orchestratorRef: string): void {
		this.db
			.query("DELETE FROM servers_crafty_rcon WHERE orchestrator_ref = ?")
			.run(orchestratorRef);
	}
}
