import type { TokenStorage } from "../../../persistence/token-storage/token-storage";
import type { ManagedCraftyConfig } from "../crafty-config";
import type { CraftyLoginResponse } from "../dtos/crafty-login-response";
import type { CraftyAuthProvider } from "./crafty-auth-provider";

export class ManagedCraftyAuthProvider implements CraftyAuthProvider {
	constructor(
		private readonly config: ManagedCraftyConfig,
		private readonly tokenStorage: TokenStorage,
	) {}

	async getToken(): Promise<string> {
		const storedToken = await this.tokenStorage.getToken();

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

			await this.tokenStorage.setToken(token);
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

		await this.tokenStorage.setToken(token);
		return token;
	}
}
