import { Configuration } from "../../../dependencies/crafty-api-client";
import { CreateMinecraftServerUseCase, type CreateMinecraftServerRequest } from "../../core";
import { CraftyMinecraftServerProvider } from "../../core/providers/crafty/crafty-minecraft-server-provider";

const command = Bun.argv[2];

const craftyApiClientConfiguration = new Configuration({
    accessToken: process.env.CRAFTY_API_TOKEN,
    basePath: process.env.CRAFTY_API_BASE_PATH ?? "https://localhost:8443",
});
const craftyProvider = new CraftyMinecraftServerProvider(craftyApiClientConfiguration);
const createMinecraftServerUseCase = new CreateMinecraftServerUseCase(craftyProvider);

switch (command) {
    case "create":

        const newMinecraftServerData: CreateMinecraftServerRequest = {
            name: `Example Server ${Bun.randomUUIDv7()}`,
            type: "java",
            version: "26.3",
        };

        const result = await createMinecraftServerUseCase.execute(newMinecraftServerData);
        console.log("Server created:", result);
        break;
    default:
        console.log(`Unknown command: ${command}`);
        break;
}