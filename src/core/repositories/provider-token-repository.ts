import type { Database } from "bun:sqlite";
import type { ProviderName } from "../providers/provider-name";

export interface ProviderTokenRepository {
	getToken(): Promise<string | undefined>;
	setToken(token: string): Promise<void>;
}

export class SqliteProviderTokenRepository implements ProviderTokenRepository {
	constructor(
		private readonly db: Database,
		private readonly providerName: ProviderName,
	) {}

	async getToken(): Promise<string | undefined> {
		const row = this.db.query("SELECT token FROM ProviderToken WHERE name = ?").get(this.providerName) as { token: string } | null;
		return row?.token;
	}
	async setToken(token: string): Promise<void> {
		this.db
			.query(
				"INSERT INTO ProviderToken (name, token) VALUES (?, ?) ON CONFLICT (name) DO UPDATE SET token = excluded.token",
			)
			.run(this.providerName, token);
	}
}
