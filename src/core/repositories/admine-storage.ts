import type { ProviderName } from "../providers/provider-name";

export interface AdmineStorage {
	getProviderToken(providerName: ProviderName): Promise<string | undefined>;
	setProviderToken(providerName: ProviderName, token: string): Promise<void>;
}
