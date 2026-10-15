import type { Database } from "bun:sqlite";

export interface CraftyAuthTokenRepository {
	get(): Promise<string | undefined>;
	set(token: string): Promise<void>;
}

export class SqliteCraftyAuthTokenRepository
	implements CraftyAuthTokenRepository
{
	constructor(private readonly db: Database) {}

	async get(): Promise<string | undefined> {
		const row = this.db
			.query("SELECT token FROM CraftyAuthToken WHERE id = 1")
			.get() as { token: string } | null;
		return row?.token;
	}

	async set(token: string): Promise<void> {
		this.db
			.query(
				"INSERT INTO CraftyAuthToken (id, token) VALUES (1, ?) ON CONFLICT (id) DO UPDATE SET token = excluded.token",
			)
			.run(token);
	}
}
