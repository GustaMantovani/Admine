import type { AdmineStorage } from "../../../repositories/admine-storage";
import { CRAFTY_PROVIDER_NAME } from "../../provider-name";
import type { ManagedCraftyConfig } from "../crafty-config";
import type { CraftyLoginResponse } from "../dtos/crafty-login-response";
import type { CraftyAuthProvider } from "./crafty-auth-provider";

export class ManagedCraftyAuthProvider implements CraftyAuthProvider {
	constructor(
		private readonly config: ManagedCraftyConfig,
		private readonly storage: AdmineStorage,
	) {}

	async getToken(): Promise<string> {
		const storedToken =
			await this.storage.getProviderToken(CRAFTY_PROVIDER_NAME);

		if (!storedToken) {
			const response = await fetch(`${this.config.apiBaseUrl}/auth/login`, {
				method: "POST",
				body: JSON.stringify({
					username: this.config.username,
					password: this.config.password,
				}),
				headers: { "Content-Type": "application/json" },
			});

			const body = (await response.json()) as CraftyLoginResponse;
			const token = body.data.token;

			await this.storage.setProviderToken(CRAFTY_PROVIDER_NAME, token);
			return token;
		}

		return storedToken;
	}

	async getNewToken(): Promise<string> {
		const response = await fetch(`${this.config.apiBaseUrl}/auth/login`, {
			method: "POST",
			body: JSON.stringify({
				username: this.config.username,
				password: this.config.password,
			}),
			headers: { "Content-Type": "application/json" },
		});

		const body = (await response.json()) as CraftyLoginResponse;
		const token = body.data.token;

		await this.storage.setProviderToken(CRAFTY_PROVIDER_NAME, token);
		return token;
	}
}
