import type { RootDatabase } from "lmdb";
import type { ProviderName } from "../types/provider-name";
import type { AdmineStorage } from "./admine-storage";

const providerTokenKey = (providerName: ProviderName) =>
	`provider-token:${providerName}`;
const RCON_PASSWORD_KEY = "rcon-password";

export class LmdbAdmineStorage implements AdmineStorage {
	constructor(private readonly db: RootDatabase) {}

	async getProviderToken(
		providerName: ProviderName,
	): Promise<string | undefined> {
		return this.db.get(providerTokenKey(providerName));
	}

	async setProviderToken(
		providerName: ProviderName,
		token: string,
	): Promise<void> {
		await this.db.put(providerTokenKey(providerName), token);
	}

	async getRconPassword(): Promise<string | undefined> {
		return this.db.get(RCON_PASSWORD_KEY);
	}

	async setRconPassword(password: string): Promise<void> {
		await this.db.put(RCON_PASSWORD_KEY, password);
	}
}
