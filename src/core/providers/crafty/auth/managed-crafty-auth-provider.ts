import type { ManagedCraftyDeployment } from "../crafty-config";
import type { CraftyLoginResponse } from "../dtos/crafty-login-response";
import type { CraftyAuthTokenRepository } from "../persistence/crafty-auth-token-repository";
import type { CraftyAuthProvider } from "./crafty-auth-provider";

export class ManagedCraftyAuthProvider implements CraftyAuthProvider {
	constructor(
		private readonly apiBaseUrl: string,
		private readonly deployment: ManagedCraftyDeployment,
		private readonly repository: CraftyAuthTokenRepository,
	) {}

	async getToken(): Promise<string> {
		const storedToken = await this.repository.get();

		if (!storedToken) {
			const response = await fetch(`${this.apiBaseUrl}/auth/login`, {
				method: "POST",
				body: JSON.stringify({
					username: this.deployment.username,
					password: this.deployment.password,
				}),
				headers: { "Content-Type": "application/json" },
			});

			const body = (await response.json()) as CraftyLoginResponse;
			const token = body.data.token;

			await this.repository.set(token);
			return token;
		}

		return storedToken;
	}

	async getNewToken(): Promise<string> {
		const response = await fetch(`${this.apiBaseUrl}/auth/login`, {
			method: "POST",
			body: JSON.stringify({
				username: this.deployment.username,
				password: this.deployment.password,
			}),
			headers: { "Content-Type": "application/json" },
		});

		const body = (await response.json()) as CraftyLoginResponse;
		const token = body.data.token;

		await this.repository.set(token);
		return token;
	}
}
