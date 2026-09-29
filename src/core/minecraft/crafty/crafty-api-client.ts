import { CraftyApiError } from "../../errors";
import type { CraftyAuthProvider } from "./auth/crafty-auth-provider";
import type { CraftyAction } from "./dtos/crafty-action";
import type { CraftyActionExecutionResult } from "./dtos/crafty-action-execution-result";
import type { CraftyCreateServer } from "./dtos/crafty-create-server";
import type { CraftyCreateServerResult } from "./dtos/crafty-create-server-result";
import type { CraftyResponse } from "./dtos/crafty-response";
import type { CraftyServer } from "./dtos/crafty-server";
import type { CraftyServerDetail } from "./dtos/crafty-server-detail";
import type { CraftyServerStats } from "./dtos/crafty-server-stats";
import type { CraftyServerStatus } from "./dtos/crafty-server-status";
import type { CraftyStopAllResult } from "./dtos/crafty-stop-all-result";

export class CraftyApiClient {
	constructor(
		private readonly baseUrl: string,
		private readonly auth: CraftyAuthProvider,
	) {}

	async createServer(
		body: CraftyCreateServer,
	): Promise<CraftyResponse<CraftyCreateServerResult>> {
		return this.request("POST", "/servers", body);
	}

	async listServers(): Promise<CraftyResponse<CraftyServer[]>> {
		return this.request("GET", "/servers");
	}

	async getServer(
		serverId: string,
	): Promise<CraftyResponse<CraftyServerDetail>> {
		return this.request("GET", `/servers/${serverId}`);
	}

	async getServersStatus(): Promise<CraftyResponse<CraftyServerStatus[]>> {
		return this.request("GET", "/servers/status");
	}

	async getServerStats(
		serverId: string,
	): Promise<CraftyResponse<CraftyServerStats>> {
		return this.request("GET", `/servers/${serverId}/stats`);
	}

	async getServerLogs(serverId: string): Promise<CraftyResponse<string[]>> {
		return this.request("GET", `/servers/${serverId}/logs`);
	}

	async sendAction(
		serverId: string,
		action: CraftyAction,
	): Promise<CraftyActionExecutionResult> {
		return this.request("POST", `/servers/${serverId}/action/${action}`);
	}

	// Crafty reads the stdin body as raw text, not JSON
	async sendCommand(
		serverId: string,
		command: string,
	): Promise<CraftyActionExecutionResult> {
		return this.request("POST", `/servers/${serverId}/stdin`, command);
	}

	async deleteServer(
		serverId: string,
		removeFiles: boolean,
	): Promise<CraftyActionExecutionResult> {
		return this.request(
			"DELETE",
			`/servers/${serverId}${removeFiles ? "?files=true" : ""}`,
		);
	}

	async stopAllServers(): Promise<CraftyResponse<CraftyStopAllResult>> {
		return this.request("POST", "/servers/stop_all");
	}

	private async request<T>(
		method: string,
		path: string,
		body?: unknown,
	): Promise<T> {
		const isText = typeof body === "string";

		const response = await fetch(`${this.baseUrl}${path}`, {
			method,
			headers: {
				Authorization: `Bearer ${await this.auth.getToken()}`,
				"Content-Type": isText ? "text/plain" : "application/json",
			},
			body:
				body === undefined ? undefined : isText ? body : JSON.stringify(body),
		});

		if (!response.ok) {
			throw new CraftyApiError(response.status, await response.text());
		}

		return (await response.json()) as T;
	}
}
