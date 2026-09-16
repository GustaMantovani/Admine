import type { CraftyCreateServer } from "./dtos/crafty-create-server";
import type { CraftyCreateServerResult } from "./dtos/crafty-create-server-result";
import type { CraftyResponse } from "./dtos/crafty-response";
import type { CraftyServer } from "./dtos/crafty-server";

export class CraftyApiClient {
    constructor(
        private readonly baseUrl: string,
        private readonly token: string,
    ) { }

    async createServer(body: CraftyCreateServer): Promise<CraftyCreateServerResult> {
        return this.request("POST", "/api/v2/servers", body);
    }

    async listServers(): Promise<CraftyServer[]> {
        return this.request("GET", "/api/v2/servers");
    }

    private async request<T>(method: string, path: string, body?: unknown, headers?: Record<string, unknown>): Promise<T> {
        const response = await fetch(`${this.baseUrl}${path}`, {
            method,
            headers: {
                Authorization: `Bearer ${this.token}`,
                "Content-Type": "application/json",
                ...headers
            },
            body: body === undefined ? undefined : JSON.stringify(body),
        });

        if (!response.ok) {
            throw new CraftyApiError(response.status, await response.text());
        }

        const json = (await response.json()) as CraftyResponse<T>;
        return json.data;
    }
}

export class CraftyApiError extends Error {
    constructor(
        public readonly status: number,
        public readonly body: string,
    ) {
        super(`Crafty API responded with ${status}: ${body}`);
    }
}