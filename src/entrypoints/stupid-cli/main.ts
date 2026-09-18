import { parseArgs } from "node:util";
import { MinecraftServerUseCases, type CreateMinecraftServerRequest } from "../../core";
import { CRAFTY_API_BASE_PATH, CRAFTY_API_TOKEN, DEFAULT_MINECRAFT_SERVER_TYPE } from "../../core/constants";
import { CraftyMinecraftServerProvider } from "../../core/providers/crafty/crafty-minecraft-server-provider";
import { CraftyApiClient, CraftyApiError } from "../../core/providers/crafty/crafty-api-client";
import { MinecraftServerType } from "../../core/domain/minecraft-server-type";
import { MinecraftServerActionError, UnknownMinecraftServerType } from "../../core/errors";
import { adjectives, animals, colors, uniqueNamesGenerator } from "unique-names-generator";

class UsageError extends Error { }

const { positionals, values: flags } = parseArgs({
    args: Bun.argv.slice(2),
    allowPositionals: true,
    options: {
        files: { type: "boolean", default: false },
        lines: { type: "string" },
        yes: { type: "boolean", short: "y", default: false },
    },
});

const [command, ...args] = positionals;

function requireArg(index: number, name: string): string {
    const value = args[index];
    if (value === undefined) throw new UsageError(`Missing argument <${name}> for "${command}"`);
    return value;
}

const craftyApiClient = new CraftyApiClient(
    CRAFTY_API_BASE_PATH,
    CRAFTY_API_TOKEN,
);
const craftyProvider = new CraftyMinecraftServerProvider(craftyApiClient);
const minecraftServers = new MinecraftServerUseCases(craftyProvider);

const USAGE_MESSAGE = `Usage: bun run cli <command> [args]

  create [name] [type] [version]   Create a server (type: java | bedrock | hytale)
  list                             List servers and their status
  info <id>                        Show a server
  stats <id>                       Show CPU, memory and players of a server
  logs <id> [--lines N]            Show the server console output
  start | stop | restart | kill <id>
  backup <id>                      Start a backup
  update <id>                      Update the server executable
  clone <id>                       Clone a stopped server
  eula <id>                        Accept the Minecraft EULA (Java servers won't start without it)
  cmd <id> <command...>            Send a command to the server console
  delete <id> [--files] [-y]       Delete a server (--files also removes its files)
  stop-all                         Stop every running server`;

const commands: Record<string, () => Promise<void>> = {
    'create': async () => {
        const [nameArg, typeArg, versionArg] = args;

        const newMinecraftServerData: CreateMinecraftServerRequest = {
            name: nameArg || uniqueNamesGenerator({
                dictionaries: [adjectives, animals, colors],
                separator: ' ',
                length: 2
            }),
            type: typeArg?.toLocaleLowerCase() as MinecraftServerType || DEFAULT_MINECRAFT_SERVER_TYPE,
            version: versionArg,
        };

        const id = await minecraftServers.create(newMinecraftServerData);
        console.log("Server created:", { ...newMinecraftServerData, id });
    },

    'list': async () => console.table(await minecraftServers.list()),

    'info': async () => console.log(await minecraftServers.get(requireArg(0, "id"))),

    'stats': async () => console.log(await minecraftServers.stats(requireArg(0, "id"))),

    'logs': async () => {
        const lastLines = flags.lines === undefined ? undefined : Number(flags.lines);
        if (lastLines !== undefined && (!Number.isInteger(lastLines) || lastLines <= 0)) {
            throw new UsageError(`--lines must be a positive integer, got "${flags.lines}"`);
        }

        const lines = await minecraftServers.logs(requireArg(0, "id"), lastLines);
        console.log(lines.join("\n"));
    },

    'start': async () => {
        await minecraftServers.start(requireArg(0, "id"));
        console.table(await minecraftServers.list());
    },

    'stop': async () => {
        await minecraftServers.stop(requireArg(0, "id"));
        console.log("Stop requested");
    },

    'restart': async () => {
        await minecraftServers.restart(requireArg(0, "id"));
        console.log("Restart requested");
    },

    'kill': async () => {
        await minecraftServers.kill(requireArg(0, "id"));
        console.log("Kill requested");
    },

    'backup': async () => {
        await minecraftServers.backup(requireArg(0, "id"));
        console.log("Backup started");
    },

    'update': async () => {
        await minecraftServers.update(requireArg(0, "id"));
        console.log("Update started");
    },

    'clone': async () => {
        await minecraftServers.clone(requireArg(0, "id"));
        console.log("Clone started");
    },

    'eula': async () => {
        await minecraftServers.acceptEula(requireArg(0, "id"));
        console.log("EULA accepted");
    },

    'cmd': async () => {
        const id = requireArg(0, "id");
        const consoleCommand = args.slice(1).join(" ");
        if (consoleCommand === "") throw new UsageError(`Missing argument <command> for "cmd"`);

        await minecraftServers.sendCommand(id, consoleCommand);
        console.log(`Sent: ${consoleCommand}`);
    },

    'delete': async () => {
        const id = requireArg(0, "id");

        if (!flags.yes) {
            const server = await minecraftServers.get(id);
            const what = flags.files ? "and ALL its files" : "(files are kept on disk)";
            const answer = prompt(`Delete server "${server.name}" (${id}) ${what}? Type the server name to confirm:`);
            if (answer !== server.name) {
                console.log("Aborted");
                return;
            }
        }

        await minecraftServers.delete(id, flags.files);
        console.log("Server deleted");
    },

    'stop-all': async () => {
        const count = await minecraftServers.stopAll();
        console.log(`Stopping ${count} server(s)`);
    },
}

const handler = commands[command!];

if (handler === undefined) {
    if (command !== undefined) console.error(`Unknown command: ${command}\n`);
    console.log(USAGE_MESSAGE);
    process.exitCode = 1;
} else {
    try {
        await handler();
    } catch (error) {
        if (
            error instanceof UsageError ||
            error instanceof CraftyApiError ||
            error instanceof MinecraftServerActionError ||
            error instanceof UnknownMinecraftServerType
        ) {
            console.error(`Error: ${error.message}`);
            process.exitCode = 1;
        } else {
            throw error;
        }
    }
}
