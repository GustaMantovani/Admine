import { CreateMinecraftServerUseCase, type CreateMinecraftServerRequest } from "../../core";
import { DEFAULT_CRAFTY_API_BASE_PATH, DEFAULT_MINECRAFT_SERVER_TYPE, DEFAULT_MINECRAFT_SERVER_VERSION } from "../../core/constants";
import { CraftyMinecraftServerProvider } from "../../core/providers/crafty/crafty-minecraft-server-provider";
import { CraftyApiClient } from "../../core/providers/crafty/crafty-api-client";
import { MinecraftServerType } from "../../core/domain/minecraft-server-type";
import { adjectives, animals, colors, uniqueNamesGenerator } from "unique-names-generator";
import { ListMinecraftServersUseCase } from "../../core/use-cases/list-minecraft-servers-use-case";

let argCounter = 2;

function getNextArg(): string | undefined {
    return Bun.argv[argCounter++];
}

const command = getNextArg();

const craftyApiClient = new CraftyApiClient(
    process.env.CRAFTY_API_BASE_PATH ?? DEFAULT_CRAFTY_API_BASE_PATH,
    process.env.CRAFTY_API_TOKEN ?? "",
);
const craftyProvider = new CraftyMinecraftServerProvider(craftyApiClient);
const createMinecraftServerUseCase = new CreateMinecraftServerUseCase(craftyProvider);
const listMinecraftServerUseCase = new ListMinecraftServersUseCase(craftyProvider);

switch (command) {
    case "create":

        const nameArg = getNextArg();
        const typeArg = getNextArg()?.toLocaleLowerCase() as MinecraftServerType || DEFAULT_MINECRAFT_SERVER_TYPE;
        const versionArg = getNextArg();

        const newMinecraftServerData: CreateMinecraftServerRequest = {
            name: nameArg || uniqueNamesGenerator({
                dictionaries: [adjectives, animals, colors],
                separator: ' ',
                length: 2
            }),
            type: typeArg,
            version: versionArg || DEFAULT_MINECRAFT_SERVER_VERSION,
        };

        const id = await createMinecraftServerUseCase.execute(newMinecraftServerData);
        console.log("Server created:", {...newMinecraftServerData, id});
        break;
    case "list":
        console.table(await listMinecraftServerUseCase.execute());
        break;
    default:
        console.log(`Unknown command: ${command}`);
        break;
}
