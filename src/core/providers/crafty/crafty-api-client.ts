import type { CraftyAuthProvider } from "./auth/crafty-auth-provider";
import { CraftyApiError } from "./crafty-api-error";
import type { CraftyAction } from "./dtos/crafty-action";
import type { CraftyActionExecutionResponse } from "./dtos/crafty-action-execution-response";
import type { CraftyCreateServer } from "./dtos/crafty-create-server";
import type { CraftyCreateServerResponse } from "./dtos/crafty-create-server-response";
import type { CraftyCreateWebhook } from "./dtos/crafty-create-webhook";
import type { CraftyCreateWebhookResponse } from "./dtos/crafty-create-webhook-response";
import type { CraftyDirectory } from "./dtos/crafty-directory";
import type { CraftyFile } from "./dtos/crafty-file";
import type { CraftyFileCreateRequest } from "./dtos/crafty-file-create-request";
import type { CraftyFileDeleteRequest } from "./dtos/crafty-file-delete-request";
import type { CraftyFileOperation } from "./dtos/crafty-file-operation";
import type { CraftyFileOperationRequest } from "./dtos/crafty-file-operation-request";
import type { CraftyFilePatchRequest } from "./dtos/crafty-file-patch-request";
import type { CraftyFileRenameRequest } from "./dtos/crafty-file-rename-request";
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

	// Crafty answers with the file's content or, for a directory, its listing
	async getServerFiles(
		serverId: string,
		filePath: string,
	): Promise<CraftyResponse<CraftyFile | CraftyDirectory>> {
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

	async createServerFile(
		serverId: string,
		body: CraftyFileCreateRequest,
	): Promise<CraftyActionExecutionResponse> {
		return this.request("PUT", `/servers/${serverId}/files/create`, body);
	}

	async renameServerFile(
		serverId: string,
		body: CraftyFileRenameRequest,
	): Promise<CraftyActionExecutionResponse> {
		return this.request("PATCH", `/servers/${serverId}/files/create`, body);
	}

	async runServerFileOperation(
		serverId: string,
		operation: CraftyFileOperation,
		body: CraftyFileOperationRequest,
	): Promise<CraftyActionExecutionResponse> {
		return this.request(
			"POST",
			`/servers/${serverId}/files/${operation}`,
			body,
		);
	}

	async deleteServerFiles(
		serverId: string,
		body: CraftyFileDeleteRequest,
	): Promise<CraftyActionExecutionResponse> {
		return this.request("DELETE", `/servers/${serverId}/files`, body);
	}

	async downloadServerFile(
		serverId: string,
		pathOnServer: string,
	): Promise<ReadableStream<Uint8Array>> {
		const response = await this.request(
			"GET",
			`/servers/${serverId}/files/${pathOnServer}/download`,
			undefined,
			{ raw: true },
		);
		return response.body!;
	}

	// With raw: true, returns the Response untouched instead of parsing its JSON
	private async request(
		method: string,
		path: string,
		body: unknown,
		options: { raw: true },
	): Promise<Response>;
	private async request<T>(
		method: string,
		path: string,
		body?: unknown,
	): Promise<T>;
	private async request(
		method: string,
		path: string,
		body?: unknown,
		options?: { raw: true },
	): Promise<unknown> {
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

		if (options?.raw) {
			return response;
		}

		let parsedResponse: unknown;
		try {
			// Try to parse response
			parsedResponse = await response.json();
		} catch {}

		return parsedResponse;
	}
}
