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

    async createServer(createServerRequest: CreateMinecraftServerRequest): Promise<MinecraftServer> {
        const craftyCreateServerRequest = toCraftyCreateServerRequest(createServerRequest);
        const serverCreateCraftyResponse = await this.serversApi.apiV2ServersPost({ apiV2ServersPostRequest: craftyCreateServerRequest });

        return new MinecraftServer(
            serverCreateCraftyResponse.data!.newServerUuid!,
            createServerRequest.name,
            createServerRequest.type,
            createServerRequest.version,
            '127.0.0.1',
            [],
            craftyCreateServerRequest.port,
            createServerRequest.minimumMemory,
            createServerRequest.maximumMemory,
        );
    };


}