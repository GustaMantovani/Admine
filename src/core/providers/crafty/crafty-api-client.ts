import type { CraftyAuthProvider } from "./auth/crafty-auth-provider";
import { CraftyApiError } from "./crafty-api-error";
import type { CraftyAction } from "./dtos/crafty-action";
import type { CraftyActionExecutionResponse } from "./dtos/crafty-action-execution-response";
import type { CraftyCreateServer } from "./dtos/crafty-create-server";
import type { CraftyCreateServerResponse } from "./dtos/crafty-create-server-response";
import type { CraftyCreateWebhook } from "./dtos/crafty-create-webhook";
import type { CraftyCreateWebhookResponse } from "./dtos/crafty-create-webhook-response";
import type { CraftyFile } from "./dtos/crafty-file";
import type { CraftyFilePatchRequest } from "./dtos/crafty-file-patch-request";
import type { CraftyResponse } from "./dtos/crafty-response";
import type { CraftyServer } from "./dtos/crafty-server";
import type { CraftyServerDetail } from "./dtos/crafty-server-detail";
import type { CraftyServerStats } from "./dtos/crafty-server-stats";
import type { CraftyServerStatus } from "./dtos/crafty-server-status";
import type { CraftyStopAllResponse } from "./dtos/crafty-stop-all-response";

export class CraftyApiClient {
	constructor(
		private readonly baseUrl: string,
		private readonly auth: CraftyAuthProvider,
	) {}

	async createServer(
		body: CraftyCreateServer,
	): Promise<CraftyResponse<CraftyCreateServerResponse>> {
		return this.request("POST", "/servers", body);
	}

	async createWebhook(
		serverId: string,
		body: CraftyCreateWebhook,
	): Promise<CraftyResponse<CraftyCreateWebhookResponse>> {
		return this.request("POST", `/servers/${serverId}/webhook`, body);
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
	): Promise<CraftyActionExecutionResponse> {
		return this.request("POST", `/servers/${serverId}/action/${action}`);
	}

	// Crafty reads the stdin body as raw text, not JSON
	async sendCommand(
		serverId: string,
		command: string,
	): Promise<CraftyActionExecutionResponse> {
		return this.request("POST", `/servers/${serverId}/stdin`, command);
	}

	async deleteServer(
		serverId: string,
		removeFiles: boolean,
	): Promise<CraftyActionExecutionResponse> {
		return this.request(
			"DELETE",
			`/servers/${serverId}${removeFiles ? "?files=true" : ""}`,
		);
	}

	async stopAllServers(): Promise<CraftyResponse<CraftyStopAllResponse>> {
		return this.request("POST", "/servers/stop_all");
	}

	async getServerFiles(
		serverId: string,
		filePath: string,
	): Promise<CraftyFile> {
		return this.request("POST", `/servers/${serverId}/files`, {
			path: filePath,
		});
	}

	async updateServerFile(
		serverId: string,
		filePatch: CraftyFilePatchRequest,
	): Promise<void> {
		return this.request("PATCH", `/servers/${serverId}/files`, filePatch);
	}

	private async request<T>(
		method: string,
		path: string,
		body?: unknown,
	): Promise<T> {
		const isText = typeof body === "string";

		const url = `${this.baseUrl}${path}`;
		const request = {
			method,
			headers: {
				Authorization: `Bearer ${await this.auth.getToken()}`,
				"Content-Type": isText ? "text/plain" : "application/json",
			},
			body:
				body === undefined ? undefined : isText ? body : JSON.stringify(body),
		};

		let response = await fetch(url, request);

		// Get new token and retry
		if (response.status === 403) {
			request.headers.Authorization = await this.auth.getNewToken();
			response = await fetch(url, request);
		}

		if (!response.ok) {
			throw new CraftyApiError(response.status, await response.text());
		}

		return (await response.json()) as T;
	}
}
