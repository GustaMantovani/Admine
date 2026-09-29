import type { TokenStorage } from "../../../persistence/token-storage/token-storage";
import type { CraftyConfig } from "../crafty-config";
import type { CraftyLoginResponse } from "../dtos/crafty-login-response";
import type { CraftyAuthProvider } from "./crafty-auth-provider";

export class ManagedCraftyAuthProvider implements CraftyAuthProvider {
	constructor(
		private readonly config: CraftyConfig,
		private readonly tokenStorage: TokenStorage,
	) {}

	async getToken(): Promise<string> {
		const storedToken = await this.tokenStorage.getToken();

		if (!storedToken) {
			const response = await fetch(`${this.config.apiBaseUrl}/auth/login`, {
				method: "POST",
				body: JSON.stringify({
					username: this.config.userName,
					password: this.config.password,
				}),
				headers: { "Content-Type": "application/json" },
			});

			const body = (await response.json()) as CraftyLoginResponse;
			const token = body.data.token;

			this.tokenStorage.setToken(token);
			return token;
		}

		return storedToken;
	}
}
