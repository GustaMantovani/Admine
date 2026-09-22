import {
	adjectives,
	animals,
	colors,
	uniqueNamesGenerator,
} from "unique-names-generator";
import {
	type CreateMinecraftServerRequest,
	MinecraftServerUseCases,
} from "../../core";
import {
	CRAFTY_API_BASE_PATH,
	CRAFTY_API_TOKEN,
	DEFAULT_MINECRAFT_SERVER_TYPE,
} from "../../core/constants";
import { CraftyApiClient } from "../../core/providers/crafty/crafty-api-client";
import { CraftyMinecraftServerProvider } from "../../core/providers/crafty/crafty-minecraft-server-provider";
import type { MinecraftServerType } from "../../core/types/minecraft-server-type";

let argCounter = 2;

export function getNextStdinArg(): string | undefined {
  return Bun.argv[argCounter++];
}

type CLICommand =
	| "create"
	| "list"
	| "info"
	| "stats"
	| "logs"
	| "start"
	| "stop"
	| "restart"
	| "kill"
	| "backup"
	| "update"
	| "clone"
	| "eula"
	| "cmd"
	| "delete"
    | "stop-all"
	| "help"
	| "unknow";
const command: CLICommand = getNextStdinArg() as CLICommand | "unknow";

const craftyApiClient = new CraftyApiClient(
	CRAFTY_API_BASE_PATH,
	CRAFTY_API_TOKEN,
);
const craftyProvider = new CraftyMinecraftServerProvider(craftyApiClient);
const minecraftServerUseCases = new MinecraftServerUseCases(craftyProvider);

function cliHelp() {
    console.log(`Unknown command: ${command}`);
    const HELP_STRING = `Usage: bun run cli <command> [args]

    Servers:
      create [name] [type] [version]   Create a server (type: java | bedrock | hytale, default: java)
      list                             List servers and their status
      info <id>                        Show a server's details
      stats <id>                       Show CPU, memory, players and state of a server
      logs <id> [lines]                Show the server console (last N lines)

    Control:
      start <id>                       Start the server
      stop <id>                        Stop the server
      restart <id>                     Restart the server
      kill <id>                        Force-kill the server process
      eula <id>                        Accept the EULA (Java servers won't start without it)
      stop-all                         Stop every running server

    Maintenance:
      backup <id>                      Start a backup of the server
      update <id>                      Update the server executable
      clone <id>                       Clone a stopped server
      cmd <id> <command...>            Send a command to the server console
      delete <id> [--files]            Delete the server (--files also deletes its files)`;
    console.log(HELP_STRING);
}

const commands: Record<CLICommand, () => Promise<void>> = {
	create: async () => {
		const nameArg = getNextStdinArg();
		const typeArg =
			(getNextStdinArg()?.toLocaleLowerCase() as MinecraftServerType) ||
			DEFAULT_MINECRAFT_SERVER_TYPE;
		const versionArg = getNextStdinArg();

		const newMinecraftServerData: CreateMinecraftServerRequest = {
			name:
				nameArg ||
				uniqueNamesGenerator({
					dictionaries: [adjectives, animals, colors],
					separator: " ",
					length: 2,
				}),
			type: typeArg,
			version: versionArg,
		};

		const id = await minecraftServerUseCases.create(newMinecraftServerData);
		console.log("Server created:", { ...newMinecraftServerData, id });
	},

	list: async () => console.table(await minecraftServerUseCases.list()),

	info: async () => {
		const serverIdArg = getNextStdinArg()!;
		console.table(await minecraftServerUseCases.get(serverIdArg));
	},

	stats: async () => {
		const serverIdArg = getNextStdinArg()!;
		console.table(await minecraftServerUseCases.stats(serverIdArg));
	},

	logs: async () => {
		const serverIdArg = getNextStdinArg()!;
		const linesArg = getNextStdinArg();
		const lines = await minecraftServerUseCases.logs(
			serverIdArg,
			linesArg ? Number(linesArg) : undefined,
		);
		console.log(lines.join("\n"));
	},

	start: async () => {
		const serverIdArg = getNextStdinArg()!;
		await minecraftServerUseCases.start(serverIdArg);
		console.table(await minecraftServerUseCases.list());
	},

	stop: async () => {
		const serverIdArg = getNextStdinArg()!;
		await minecraftServerUseCases.stop(serverIdArg);
		console.table(await minecraftServerUseCases.list());
	},

	restart: async () => {
		const serverIdArg = getNextStdinArg()!;
		await minecraftServerUseCases.restart(serverIdArg);
		console.table(await minecraftServerUseCases.list());
	},

	kill: async () => {
		const serverIdArg = getNextStdinArg()!;
		await minecraftServerUseCases.kill(serverIdArg);
		console.table(await minecraftServerUseCases.list());
	},

	backup: async () => {
		const serverIdArg = getNextStdinArg()!;
		await minecraftServerUseCases.backup(serverIdArg);
		console.table(await minecraftServerUseCases.list());
	},

	update: async () => {
		const serverIdArg = getNextStdinArg()!;
		await minecraftServerUseCases.update(serverIdArg);
		console.table(await minecraftServerUseCases.list());
	},

	clone: async () => {
		const serverIdArg = getNextStdinArg()!;
		await minecraftServerUseCases.clone(serverIdArg);
		console.table(await minecraftServerUseCases.list());
	},

	eula: async () => {
		const serverIdArg = getNextStdinArg()!;
		await minecraftServerUseCases.acceptEula(serverIdArg);
		console.table(await minecraftServerUseCases.list());
	},

	cmd: async () => {
		const serverIdArg = getNextStdinArg()!;
		const consoleCommandArg = Bun.argv.slice(argCounter).join(" ");
		await minecraftServerUseCases.sendCommand(serverIdArg, consoleCommandArg);
	},

	delete: async () => {
		const serverIdArg = getNextStdinArg()!;
		const removeFilesArg = getNextStdinArg() === "--files";
		await minecraftServerUseCases.delete(serverIdArg, removeFilesArg);
		console.table(await minecraftServerUseCases.list());
	},

	"stop-all": async () => {
		await minecraftServerUseCases.stopAll();
		console.table(await minecraftServerUseCases.list());
	},

  help: async () => {
      cliHelp();
  },

  unknow: async () => {
      cliHelp();
      process.exit(1);
  },
};

const handler = commands[command];

if (handler === undefined) {
  cliHelp();
  process.exit(1);
};

await handler();
