import type { CreateMinecraftServerRequest } from "../../types/create-minecraft-server-request";
import { MinecraftServer } from "../../domain/minecraft-server";
import { toCraftyCreateServerRequest, toMinecraftServer } from "./mappers";
import type { MinecraftServerProvider } from "../minecraft-server-provider";
import type { CraftyApiClient } from "./crafty-api-client";


export class CraftyMinecraftServerProvider implements MinecraftServerProvider {
    constructor(private readonly craftyApiClient: CraftyApiClient) {};

    async createServer(createServerRequest: CreateMinecraftServerRequest): Promise<string> {
        const craftyCreateServerRequest = toCraftyCreateServerRequest(createServerRequest);
        const result = await this.craftyApiClient.createServer(craftyCreateServerRequest);
        return result.new_server_id;
    };


    async listServers(): Promise<MinecraftServer[]> {
        const servers = await this.craftyApiClient.listServers();
        return servers.map(toMinecraftServer);
    }


}
