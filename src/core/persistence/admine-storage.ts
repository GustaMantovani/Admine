import type { ProviderName } from "../types/provider-name";

export interface AdmineStorage {
	// Tokens are scoped by provider, so one provider never overwrites another's
	getProviderToken(providerName: ProviderName): Promise<string | undefined>;
	setProviderToken(providerName: ProviderName, token: string): Promise<void>;

	getRconPassword(): Promise<string | undefined>;
	setRconPassword(password: string): Promise<void>;
}
