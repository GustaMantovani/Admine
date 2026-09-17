import type { CreateMinecraftServerRequest } from "../../types/create-minecraft-server-request";
import { MinecraftServer } from "../../domain/minecraft-server";
import { toCraftyCreateServerRequest, toMinecraftServer } from "./mappers";
import type { MinecraftServerProvider } from "../minecraft-server-provider";
import type { CraftyApiClient } from "./crafty-api-client";


export class CraftyMinecraftServerProvider implements MinecraftServerProvider {
    constructor(private readonly craftyApiClient: CraftyApiClient) { };

    async create(createServerRequest: CreateMinecraftServerRequest): Promise<string> {
        const craftyCreateServerRequest = toCraftyCreateServerRequest(createServerRequest);
        const result = await this.craftyApiClient.createServer(craftyCreateServerRequest);
        return result.new_server_id;
    };


    async list(): Promise<MinecraftServer[]> {
        const [servers, statuses] = await Promise.all([
            this.craftyApiClient.listServers(),
            this.craftyApiClient.getServersStatus(),
        ]);
        const statusById = new Map(statuses.map(status => [status.id, status]));
        return servers.map(dto => toMinecraftServer(dto, statusById.get(dto.server_id)));
    }

    async start(id: string): Promise<void> {
        const result = await this.craftyApiClient.sendAction(id, 'start_server');
        if (result.status === 'error') {
            throw new Error(`Crafty failed to start server ${id}: ${result.error} | ${result.error_data}`);
        }
    }

}
