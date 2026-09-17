import { CreateMinecraftServerUseCase, type CreateMinecraftServerRequest } from "../../core";
import { CRAFTY_API_BASE_PATH, CRAFTY_API_TOKEN, DEFAULT_MINECRAFT_SERVER_TYPE, DEFAULT_MINECRAFT_SERVER_VERSION } from "../../core/constants";
import { CraftyMinecraftServerProvider } from "../../core/providers/crafty/crafty-minecraft-server-provider";
import { CraftyApiClient } from "../../core/providers/crafty/crafty-api-client";
import { MinecraftServerType } from "../../core/domain/minecraft-server-type";
import { adjectives, animals, colors, uniqueNamesGenerator } from "unique-names-generator";
import { ListMinecraftServersUseCase } from "../../core/use-cases/list-minecraft-servers-use-case";
import { StartMinecraftServerUseCase } from "../../core/use-cases/start-minecraft-server-use-case";

let argCounter = 2;

export function getNextStdinArg(): string | undefined {
    return Bun.argv[argCounter++];
}

type CLICommand = 'create' | 'start' | 'list' | 'start' | 'unknow';
const command: CLICommand = getNextStdinArg() as CLICommand | 'unknow';

const craftyApiClient = new CraftyApiClient(
    CRAFTY_API_BASE_PATH,
    CRAFTY_API_TOKEN,
);
const craftyProvider = new CraftyMinecraftServerProvider(craftyApiClient);
const createMinecraftServerUseCase = new CreateMinecraftServerUseCase(craftyProvider);
const listMinecraftServerUseCase = new ListMinecraftServersUseCase(craftyProvider);
const startMinecraftServerUseCase = new StartMinecraftServerUseCase(craftyProvider);


const commands: Record<CLICommand, () => Promise<void>> = {
    'create': async () => {
        const nameArg = getNextStdinArg();
        const typeArg = getNextStdinArg()?.toLocaleLowerCase() as MinecraftServerType || DEFAULT_MINECRAFT_SERVER_TYPE;
        const versionArg = getNextStdinArg();

        const newMinecraftServerData: CreateMinecraftServerRequest = {
            name: nameArg || uniqueNamesGenerator({
                dictionaries: [adjectives, animals, colors],
                separator: ' ',
                length: 2
            }),
            type: typeArg,
            version: versionArg,
        };

        const id = await createMinecraftServerUseCase.execute(newMinecraftServerData);
        console.log("Server created:", { ...newMinecraftServerData, id });
    },

    'list': async () => console.table(await listMinecraftServerUseCase.execute()),

    'start': async () => {
        const serverIdArg = getNextStdinArg()!;
        await startMinecraftServerUseCase.execute(serverIdArg);
        console.table(await listMinecraftServerUseCase.execute());
    },

    'unknow': async () => console.log(`Unknown command: ${command}`)
}

await commands[command]();