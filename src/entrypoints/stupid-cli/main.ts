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
import {
	CRAFTY_API_BASE_PATH,
	CRAFTY_API_TOKEN,
	DEFAULT_MINECRAFT_SERVER_TYPE,
} from "../../core/constants";
import { CraftyApiClient } from "../../core/minecraft/crafty/crafty-api-client";
import { CraftyMinecraftServerProvider } from "../../core/minecraft/crafty/crafty-minecraft-server-provider";
import type { MinecraftServerType } from "../../core/types/minecraft-server-type";
import { open } from 'lmdb';
import { LmdbTokenStorage } from "../../core/persistence/token-storage/lmdb-token-storage";
import type { CraftyConfig } from "../../core/minecraft/crafty/crafty-config";
import { StaticCraftyAuthProvider } from "../../core/minecraft/crafty/auth/static-crafty-auth-provider";
import { ManagedCraftyAuthProvider } from "../../core/minecraft/crafty/auth/managed-crafty-auth-provider";

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

if (command === 'help') {
  cliHelp();
  process.exit(0);
};

// Crafty setup

const config: CraftyConfig = {
        apiBaseUrl: "https://localhost:8443/api/v2",
        isSelfHosted: false,
        apiToken: "token-de-teste",
        userName: "admin",
        password: "e&s%fF$c@7kPs7sdLdiyv*6F3qvydP&tl*!$F8WR6K$o4Zs3N!6l2nzXt&47z8NP",
        defaultCredsFilePath: "./crafty-data/config/default.json",
        infraSetupCommand: 'docker run --name crafty_container --detach -p 8443:8443 -p 8123:8123 -p 19132:19132/udp -p 25500-25600:25500-25600 -e TZ=Etc/UTC -v "/$(pwd)/.data/crafty/backups:/crafty/backups" -v "/$(pwd)/.data/crafty/logs:/crafty/logs" -v "/$(pwd)/.data/crafty/servers:/crafty/servers" -v "/$(pwd)/.data/crafty/config:/crafty/app/config" -v "/$(pwd)/.data/crafty/import:/crafty/import" arcadiatechnology/crafty-4:latest',
        infraCheckDelay: 10000,
        infraCheckCommand: "docker inspect -f '{{.State.Running}}' crafty_container | grep -q true",
};

let lmdbKv = open('.data/admine', {});
const tokenStorage = new LmdbTokenStorage(lmdbKv);

const authProvider = config.isSelfHosted ? new StaticCraftyAuthProvider(config) : new ManagedCraftyAuthProvider(config, tokenStorage);

const craftyApiClient = new CraftyApiClient(
	CRAFTY_API_BASE_PATH,
	authProvider
);
const craftyProvider = new CraftyMinecraftServerProvider(craftyApiClient, config);
await craftyProvider.setup();
const minecraftServerService = new MinecraftServerService(craftyProvider);

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

		const id = await minecraftServerService.create(newMinecraftServerData);
		console.log("Server created:", { ...newMinecraftServerData, id });
	},

	list: async () => console.table(await minecraftServerService.list()),

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
		await minecraftServerService.sendCommand(serverIdArg, consoleCommandArg);
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

  unknow: async () => {
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
