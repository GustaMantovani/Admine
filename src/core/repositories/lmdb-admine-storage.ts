import type { RootDatabase } from "lmdb";
import type { ProviderName } from "../providers/provider-name";
import type { AdmineStorage } from "./admine-storage";

const providerTokenKey = (providerName: ProviderName) =>
	`provider-token:${providerName}`;

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
}
