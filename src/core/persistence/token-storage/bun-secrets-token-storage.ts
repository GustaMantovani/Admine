import type { TokenStorage } from "./token-storage";

const SECRET_SERVICE = "admine";
const SECRET_NAME = "minecraft-server-provider-token";

export class BunSecretsTokenStorage implements TokenStorage {
	async getToken(): Promise<string | undefined> {
		return (
			(await Bun.secrets.get({ service: SECRET_SERVICE, name: SECRET_NAME })) ??
			undefined
		);
	}

	async setToken(token: string): Promise<void> {
		await Bun.secrets.set({
			service: SECRET_SERVICE,
			name: SECRET_NAME,
			value: token,
		});
	}
}
