import {
	adjectives,
	animals,
	uniqueNamesGenerator,
} from "unique-names-generator";
import { type CreateGameServerRequest, GameServerService } from "../../core";
import type { GameServerType } from "../../core/domain/game-server-type";
import type { CraftyAuthProvider } from "../../core/providers/crafty/auth/crafty-auth-provider";
import { ManagedCraftyAuthProvider } from "../../core/providers/crafty/auth/managed-crafty-auth-provider";
import { StaticCraftyAuthProvider } from "../../core/providers/crafty/auth/static-crafty-auth-provider";
import { CraftyApiClient } from "../../core/providers/crafty/crafty-api-client";
import { CraftyGameServerProvider } from "../../core/providers/crafty/crafty-game-server-provider";
import { CRAFTY_PROVIDER_NAME } from "../../core/providers/provider-name";
import { getDatabase } from "../../core/repositories/database";
import { SqliteProviderTokenRepository } from "../../core/repositories/provider-token-repository";
import { ensureCraftyInfra } from "../../infra/crafty-docker";
import { loadCliConfig, withCraftyCreds } from "./config";

const config = await loadCliConfig();

let argCounter = 2;
function getNextStdinArg(): string | undefined {
	return Bun.argv[argCounter++];
}

function getRequiredStdinArg(argName: string): string {
	const arg = getNextStdinArg();
	if (arg === undefined) {
		console.log(`Missing argument: <${argName}>`);
		cliHelp();
		process.exit(1);
	}
	return arg;
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
const db = getDatabase(".data/admine.sqlite");
const craftyProviderTokenRepository = new SqliteProviderTokenRepository(
	db,
	CRAFTY_PROVIDER_NAME,
);

const craftyConfig = config.craftyConfig;
let authProvider: CraftyAuthProvider;
switch (craftyConfig.deployType) {
	case "self-hosted":
		authProvider = new StaticCraftyAuthProvider(craftyConfig);
		break;
	case "managed":
		await ensureCraftyInfra(craftyConfig);
		authProvider = new ManagedCraftyAuthProvider(
			await withCraftyCreds(craftyConfig),
			craftyProviderTokenRepository,
		);
}
const craftyApiClient = new CraftyApiClient(
	craftyConfig.apiBaseUrl,
	authProvider,
);
const craftyProvider = new CraftyGameServerProvider(
	craftyApiClient,
	config.serversHost,
	craftyConfig.webhooks,
	config.minecraftServer.rcon,
);

const gameServerService = new GameServerService(craftyProvider);

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
	// Servers
	create: async () => {
		const nameArg = getNextStdinArg();
		const typeArg =
			(getNextStdinArg()?.toLocaleLowerCase() as GameServerType) ||
			config.minecraftServer.type;
		const versionArg = getNextStdinArg();
		const portArg = getNextStdinArg();
		const rconPortArg = getNextStdinArg();

		const newGameServerData: CreateGameServerRequest = {
			name:
				nameArg ||
				uniqueNamesGenerator({
					dictionaries: [adjectives, animals],
					separator: " ",
					length: 2,
				}),
			type: typeArg,
			version: versionArg || config.minecraftServer.version,
			minimumMemory: config.minecraftServer.minimumMemory,
			maximumMemory: config.minecraftServer.maximumMemory,
			properties: config.minecraftServer.properties,
			port: portArg !== undefined ? parseInt(portArg, 10) : 25565, // Aqui chamar o serviço que gera a porta evitando colisão
			rconPort: rconPortArg !== undefined ? parseInt(rconPortArg, 10) : 25575,
		};

		const id = await gameServerService.create(
			newGameServerData,
			config.minecraftServer.acceptEulaOnCreate,
		);
		console.log("Server created:", { ...newGameServerData, id });
	},

	ls: async () => console.table(await gameServerService.list()),

	info: async () => {
		const serverIdArg = getRequiredStdinArg("id");
		console.table(await gameServerService.get(serverIdArg));
	},

	stats: async () => {
		const serverIdArg = getRequiredStdinArg("id");
		console.table(await gameServerService.stats(serverIdArg));
	},

	logs: async () => {
		const serverIdArg = getRequiredStdinArg("id");
		const linesArg = getNextStdinArg();
		const lines = await gameServerService.logs(
			serverIdArg,
			linesArg ? Number(linesArg) : undefined,
		);
		console.log(lines.join("\n"));
	},

	// Control
	start: async () => {
		const serverIdArg = getRequiredStdinArg("id");
		await gameServerService.start(serverIdArg);
	},

	stop: async () => {
		const serverIdArg = getRequiredStdinArg("id");
		await gameServerService.stop(serverIdArg);
	},

	restart: async () => {
		const serverIdArg = getRequiredStdinArg("id");
		await gameServerService.restart(serverIdArg);
	},

	kill: async () => {
		const serverIdArg = getRequiredStdinArg("id");
		await gameServerService.kill(serverIdArg);
	},

	eula: async () => {
		const serverIdArg = getRequiredStdinArg("id");
		await gameServerService.acceptEula(serverIdArg);
	},

	"stop-all": async () => {
		await gameServerService.stopAll();
	},

	// Maintenance
	backup: async () => {
		const serverIdArg = getRequiredStdinArg("id");
		await gameServerService.backup(serverIdArg);
	},

	update: async () => {
		const serverIdArg = getRequiredStdinArg("id");
		await gameServerService.update(serverIdArg);
	},

	clone: async () => {
		const serverIdArg = getRequiredStdinArg("id");
		await gameServerService.clone(serverIdArg);
	},

	cmd: async () => {
		const serverIdArg = getRequiredStdinArg("id");
		const consoleCommandArg = Bun.argv.slice(argCounter).join(" ");
		console.log(
			await gameServerService.sendCommand(serverIdArg, consoleCommandArg),
		);
	},

	delete: async () => {
		const serverIdArg = getRequiredStdinArg("id");
		const removeFilesArg = getNextStdinArg() === "--files";
		await gameServerService.delete(serverIdArg, removeFilesArg);
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
}

await handler();
