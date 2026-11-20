import type { Database } from "bun:sqlite";
import { CraftyApiError } from "./crafty-api-client";
import type {
	ManagedCraftyDeployment,
	SelfHostedCraftyDeployment,
} from "./crafty-config";
import type { CraftyCredsFile, CraftyLoginResponse } from "./dtos";

export interface CraftyAuthProvider {
	getToken(): Promise<string>;
	getNewToken(): Promise<string>;
}

export class StaticCraftyAuthProvider implements CraftyAuthProvider {
	constructor(private readonly deployment: SelfHostedCraftyDeployment) {}

	async getToken(): Promise<string> {
		return this.deployment.apiToken;
	}

	async getNewToken(): Promise<string> {
		return this.deployment.apiToken;
	}
}

export type CraftyCredentials = { username: string; password: string };

export class ManagedCraftyAuthProvider implements CraftyAuthProvider {
	constructor(
		private readonly apiBaseUrl: string,
		private readonly credentials: CraftyCredentials,
		private readonly tokens: CraftyAuthTokenRepository,
	) {}

	async getToken(): Promise<string> {
		return this.tokens.get() ?? this.getNewToken();
	}

	async getNewToken(): Promise<string> {
		const response = await fetch(`${this.apiBaseUrl}/auth/login`, {
			method: "POST",
			body: JSON.stringify(this.credentials),
			headers: { "Content-Type": "application/json" },
		});
		if (!response.ok)
			throw new CraftyApiError(response.status, await response.text());

		const body = (await response.json()) as CraftyLoginResponse;
		this.tokens.set(body.data.token);
		return body.data.token;
	}
}

export class CraftyAuthTokenRepository {
	constructor(private readonly db: Database) {}

	get(): string | undefined {
		const row = this.db
			.query("SELECT token FROM servers_crafty_auth_token WHERE id = 1")
			.get() as { token: string } | null;
		return row?.token;
	}

	set(token: string): void {
		this.db
			.query(
				"INSERT INTO servers_crafty_auth_token (id, token) VALUES (1, ?) ON CONFLICT (id) DO UPDATE SET token = excluded.token",
			)
			.run(token);
	}
}

export async function resolveCraftyCredentials(
	deployment: ManagedCraftyDeployment,
	credsFilePath: string,
): Promise<CraftyCredentials> {
	if (deployment.username !== undefined && deployment.password !== undefined)
		return { username: deployment.username, password: deployment.password };

	const creds = (await Bun.file(credsFilePath).json()) as CraftyCredsFile;
	return {
		username: deployment.username ?? creds.username,
		password: deployment.password ?? creds.password,
	};
}
