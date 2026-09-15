import { Configuration as CraftyApiClientConfiguration } from "../../../../dependencies/crafty-api-client";
import { ServersApi } from "../../../../dependencies/crafty-api-client";
import type { CreateMinecraftServerRequest } from "../../types/create-minecraft-server-request";
import { MinecraftServer } from "../../domain/minecraft-server";
import { toCraftyCreateServerRequest } from "./mappers";
import type { MinecraftServerProvider } from "../minecraft-server-provider";


export class CraftyMinecraftServerProvider implements MinecraftServerProvider {
    private readonly serversApi: ServersApi;

    constructor(private readonly craftyApiClientConfiguration: CraftyApiClientConfiguration) {
        this.serversApi = new ServersApi(this.craftyApiClientConfiguration);
    };

    async createServer(createServerRequest: CreateMinecraftServerRequest): Promise<string> {
        const craftyCreateServerRequest = toCraftyCreateServerRequest(createServerRequest);
        return (await this.serversApi.apiV2ServersPost({ apiV2ServersPostRequest: craftyCreateServerRequest }))!.data!.newServerId!;
    };


    async listServers(): Promise<MinecraftServer[]> {
        const serverListCraftyResponse = await this.serversApi.apiV2ServersGet();
        console.log(null, 4, serverListCraftyResponse.data)
        return [];
    }


}