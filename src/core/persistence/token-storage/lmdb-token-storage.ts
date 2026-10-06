import type { RootDatabase } from "lmdb";
import type { TokenStorage } from "./token-storage";

const LMDB_TOKEN_KEY = "admine-minecraft-server-provider-token";

export class LmdbTokenStorage implements TokenStorage {
	constructor(private readonly db: RootDatabase) {}

	async getToken(): Promise<string> {
		return this.db.get(LMDB_TOKEN_KEY);
	}

	async setToken(token: string): Promise<void> {
		await this.db.put(LMDB_TOKEN_KEY, token);
	}
}
