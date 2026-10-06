import type { RootDatabase } from "lmdb";
import type { RconPasswordStorage } from "./rcon-password-storage";

const LMDB_RCON_PASSWORD_KEY = "admine-minecraft-server-provider-rcon-password";

export class LmdbRconPasswordStorage implements RconPasswordStorage {
	constructor(private readonly db: RootDatabase) {}

	async getPassword(): Promise<string> {
		return this.db.get(LMDB_RCON_PASSWORD_KEY);
	}

	async setPassword(RconPassword: string): Promise<void> {
		await this.db.put(LMDB_RCON_PASSWORD_KEY, RconPassword);
	}
}
