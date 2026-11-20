import * as nodePath from "node:path";
import type { BunFile } from "bun";
import type { CraftyAuthProvider } from "./crafty-auth";
import type {
	CraftyAction,
	CraftyActionExecutionResponse,
	CraftyCreateServer,
	CraftyCreateServerResponse,
	CraftyCreateWebhook,
	CraftyCreateWebhookResponse,
	CraftyDirectory,
	CraftyFile,
	CraftyFileCreateRequest,
	CraftyFileDeleteRequest,
	CraftyFileOperation,
	CraftyFileOperationRequest,
	CraftyFilePatchRequest,
	CraftyFileRenameRequest,
	CraftyResponse,
	CraftyServer,
	CraftyServerDetail,
	CraftyServerStats,
	CraftyStopAllResponse,
} from "./dtos";

export class CraftyApiError extends Error {
	constructor(
		public readonly status: number,
		public readonly body: string,
	) {
		super(
			`Crafty API responded with ${status}: ${CraftyApiError.summarize(body)}`,
		);
		this.name = "CraftyApiError";
	}

	private static summarize(body: string): string {
		try {
			const json = JSON.parse(body) as { error?: string; error_data?: string };
			if (json.error !== undefined)
				return json.error_data
					? `${json.error} (${json.error_data})`
					: json.error;
		} catch {}

		const lines = body.trim().split("\n");
		return lines[lines.length - 1] ?? body;
	}
}

const UPLOAD_CHUNK_SIZE = 10 * 1024 * 1024;

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
	): Promise<{
		contentStream: ReadableStream<Uint8Array>;
		fileNameOnServer?: string;
	}> {
		const response = await this.request(
			"GET",
			`/servers/${serverId}/files/${encodeURIComponent(pathOnServer)}/download`,
			undefined,
			undefined,
			{ raw: true },
		);

		const fileNameOnServer = response.headers
			.get("content-disposition")
			?.match(/filename="([^"]+)"/)?.[1];

		if (response.body === null) throw new Error("No response body");
		return { contentStream: response.body, fileNameOnServer };
	}

	async uploadServerFile(
		serverId: string,
		path: string,
		file: BunFile,
	): Promise<void> {
		const apiPath = `/servers/${serverId}/files/upload`;
		const fileId = Bun.randomUUIDv7();
		const fileName = nodePath.basename(path);
		const location = nodePath.dirname(path) || ".";
		const chunks = splitFile(UPLOAD_CHUNK_SIZE, file);
		const uploadHeaders = {
			"Content-Type": "application/octet-stream",
			chunked: true,
			fileId,
			fileName,
			fileSize: file.size,
			totalChunks: chunks.length,
			location,
		};

		await this.request("POST", apiPath, undefined, uploadHeaders, {
			raw: true,
		});

		for (const [index, chunk] of chunks.entries()) {
			const response = await this.request<{ status: string }>(
				"POST",
				apiPath,
				chunk,
				{
					...uploadHeaders,
					chunkId: index,
					chunkHash: await hashChunk(chunk),
				},
			);

			if (index === chunks.length - 1 && response.status !== "completed")
				throw new Error(
					"Every chunk was sent, but Crafty didn't confirm the upload",
				);
		}
	}

	private async request(
		method: string,
		path: string,
		body: unknown,
		extraHeaders: Record<string, string | boolean | number> | undefined,
		options: { raw: true },
	): Promise<Response>;
	private async request<T>(
		method: string,
		path: string,
		body?: unknown,
		extraHeaders?: Record<string, string | boolean | number>,
	): Promise<T>;

	private async request(
		method: string,
		path: string,
		body?: unknown,
		extraHeaders?: Record<string, string | boolean | number>,
		options?: { raw: true },
	): Promise<unknown> {
		const isText = typeof body === "string";
		const isBinary =
			body instanceof Blob ||
			body instanceof ArrayBuffer ||
			ArrayBuffer.isView(body) ||
			body instanceof FormData ||
			body instanceof URLSearchParams ||
			body instanceof ReadableStream;

		const url = `${this.baseUrl}${path}`;
		const request = {
			method,
			headers: {
				Authorization: `Bearer ${await this.auth.getToken()}`,
				"Content-Type": isText ? "text/plain" : "application/json",
				...extraHeaders,
			},
			body:
				body === undefined || isText || isBinary
					? (body as Bun.BodyInit | undefined)
					: JSON.stringify(body),
		};

		let response = await fetch(url, request);

		if (response.status === 403) {
			request.headers.Authorization = `Bearer ${await this.auth.getNewToken()}`;
			response = await fetch(url, request);
		}

		if (!response.ok)
			throw new CraftyApiError(response.status, await response.text());
		if (options?.raw) return response;

		try {
			return await response.json();
		} catch {
			return undefined;
		}
	}
}

function splitFile(chunkSize: number, file: BunFile): Blob[] {
	const chunks: Blob[] = [];
	for (let offset = 0; offset < file.size; offset += chunkSize)
		chunks.push(file.slice(offset, offset + chunkSize));
	return chunks;
}

async function hashChunk(chunk: Blob): Promise<string> {
	return new Bun.CryptoHasher("sha256")
		.update(await chunk.arrayBuffer())
		.digest("hex");
}
