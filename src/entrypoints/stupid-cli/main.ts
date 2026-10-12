import {
	adjectives,
	animals,
	colors,
	uniqueNamesGenerator,
} from "unique-names-generator";
import {
	type CreateMinecraftServerRequest,
	MinecraftServerService,
} from "../../core";
import { CraftyApiClient } from "../../core/minecraft/crafty/crafty-api-client";
import { CraftyMinecraftServerProvider } from "../../core/minecraft/crafty/crafty-minecraft-server-provider";
import { ensureCraftyInfra } from "./crafty-infra";
import type { MinecraftServerType } from "../../core/types/minecraft-server-type";
import { open } from 'lmdb';
import { LmdbAdmineStorage } from "../../core/persistence/lmdb-admine-storage";
import { StaticCraftyAuthProvider } from "../../core/minecraft/crafty/auth/static-crafty-auth-provider";
import { ManagedCraftyAuthProvider } from "../../core/minecraft/crafty/auth/managed-crafty-auth-provider";
import type { CraftyAuthProvider } from "../../core/minecraft/crafty/auth/crafty-auth-provider";
import { loadCliConfig, withCraftyCreds } from "./config";
import { RconClient } from "../../core/minecraft/rcon/rcon-client";

const config = await loadCliConfig();

let argCounter = 2;
export function getNextStdinArg(): string | undefined {
  return Bun.argv[argCounter++];
}

type CliCommand =
	| "create"
	| "ls"
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
	| "unknown";
const command: CliCommand = getNextStdinArg() as CliCommand | "unknown";

// Dependencies
let lmdbKv = open('.data/admine', {});
const storage = new LmdbAdmineStorage(lmdbKv);
const craftyConfig = config.craftyConfig;
let authProvider: CraftyAuthProvider;
switch (craftyConfig.deployType) {
	case 'self-hosted':
		authProvider = new StaticCraftyAuthProvider(craftyConfig);
		break;
	case 'managed':
		await ensureCraftyInfra(craftyConfig);
		authProvider = new ManagedCraftyAuthProvider(
			await withCraftyCreds(craftyConfig),
			storage,
		);
}
const craftyApiClient = new CraftyApiClient(
	craftyConfig.apiBaseUrl,
	authProvider
);
const craftyProvider = new CraftyMinecraftServerProvider(
	craftyApiClient,
	config.serversHost,
	craftyConfig.webhooks,
	config.minecraftServer.rcon
);

const rconClient = new RconClient(config.minecraftServer.rcon);

const minecraftServerService = new MinecraftServerService(craftyProvider, rconClient);

function cliHelp() {
    const HELP_STRING = `Usage: bun run cli <command> [args]

    Servers:
      create [name] [type] [version]   Create a server (type: java | bedrock | hytale, default: java)
      list                             List servers
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

const commands: Record<CliCommand, () => Promise<void>> = {
	create: async () => {
		const nameArg = getNextStdinArg();
		const typeArg =
			(getNextStdinArg()?.toLocaleLowerCase() as MinecraftServerType) ||
			config.minecraftServer.type;
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
			version: versionArg || config.minecraftServer.version,
			minimumMemory: config.minecraftServer.minimumMemory,
			maximumMemory: config.minecraftServer.maximumMemory,
			properties: config.minecraftServer.properties
		};

		const id = await minecraftServerService.create(newMinecraftServerData, config.minecraftServer.acceptEulaOnCreate);
		console.log("Server created:", { ...newMinecraftServerData, id });
	},

	ls: async () => console.table(await minecraftServerService.list()),

	info: async () => {
		const serverIdArg = getNextStdinArg()!;
		console.table(await minecraftServerService.get(serverIdArg));
	},

	stats: async () => {
		const serverIdArg = getNextStdinArg()!;
		console.table(await minecraftServerService.stats(serverIdArg));
	},

	logs: async () => {
		const serverIdArg = getNextStdinArg()!;
		const linesArg = getNextStdinArg();
		const lines = await minecraftServerService.logs(
			serverIdArg,
			linesArg ? Number(linesArg) : undefined,
		);
		console.log(lines.join("\n"));
	},

	start: async () => {
		const serverIdArg = getNextStdinArg()!;
		await minecraftServerService.start(serverIdArg);
	},

	stop: async () => {
		const serverIdArg = getNextStdinArg()!;
		await minecraftServerService.stop(serverIdArg);
	},

	restart: async () => {
		const serverIdArg = getNextStdinArg()!;
		await minecraftServerService.restart(serverIdArg);
	},

	kill: async () => {
		const serverIdArg = getNextStdinArg()!;
		await minecraftServerService.kill(serverIdArg);
	},

	backup: async () => {
		const serverIdArg = getNextStdinArg()!;
		await minecraftServerService.backup(serverIdArg);
	},

	update: async () => {
		const serverIdArg = getNextStdinArg()!;
		await minecraftServerService.update(serverIdArg);
	},

	clone: async () => {
		const serverIdArg = getNextStdinArg()!;
		await minecraftServerService.clone(serverIdArg);
	},

	eula: async () => {
		const serverIdArg = getNextStdinArg()!;
		await minecraftServerService.acceptEula(serverIdArg);
	},

	cmd: async () => {
		const serverIdArg = getNextStdinArg()!;
		const consoleCommandArg = Bun.argv.slice(argCounter).join(" ");
		console.log(await minecraftServerService.sendCommand(serverIdArg, consoleCommandArg));
	},

	delete: async () => {
		const serverIdArg = getNextStdinArg()!;
		const removeFilesArg = getNextStdinArg() === "--files";
		await minecraftServerService.delete(serverIdArg, removeFilesArg);
	},

	"stop-all": async () => {
		await minecraftServerService.stopAll();
	},

  help: async () => {
    cliHelp();
  },

  unknown: async () => {
    console.log(`Unknown command: ${command}`);
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
