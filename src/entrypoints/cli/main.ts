import path from "node:path";
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
import { SqliteCraftyAuthTokenRepository } from "../../core/providers/crafty/persistence/crafty-auth-token-repository";
import { SqliteCraftyRconPortRepository } from "../../core/providers/crafty/persistence/crafty-rcon-port-repository";
import { ensureCraftyInfra } from "../../infra/crafty-docker";
import { getDatabase } from "../../infra/database";
import { loadCliConfig, withCraftyCreds } from "./config";

const config = await loadCliConfig();

let argCounter = 2;
function getNextStdinArg(): string | undefined {
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
	| "delete-all"
	| "stop-all"
	| "files"
	| "help"
	| "unknown";
const command: CliCommand = getNextStdinArg() as CliCommand | "unknown";

type FilesSubcommand =
	| "ls"
	| "cat"
	| "write"
	| "touch"
	| "mkdir"
	| "rename"
	| "mv"
	| "cp"
	| "rm"
	| "download"
	| "upload";

// Dependencies
const db = getDatabase(".data/admine.sqlite");

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
			new SqliteCraftyAuthTokenRepository(db),
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
	new SqliteCraftyRconPortRepository(db),
);

const gameServerService = new GameServerService(
	craftyProvider,
	config.minecraftServer.rcon,
);

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
      delete <id> [--files]            Delete the server (--files also deletes its files)
      delete-all [--files]             Delete every server (--files also deletes their files)

    Files (paths are relative to the server's root directory):
      files ls <id> [path]             List a directory (default: the server's root)
      files cat <id> <path>            Print a file
      files write <id> <path>          Replace an existing file's content with stdin
      files touch <id> <path>          Create an empty file
      files mkdir <id> <path>          Create a directory
      files rename <id> <path> <name>  Rename a file or directory
      files mv <id> <dir> <paths...>   Move files or directories into a directory
      files cp <id> <dir> <paths...>   Copy files or directories into a directory
      files rm <id> <paths...>         Delete files or directories
      files download <id> <path> [dest] Download a file to a local path
      files upload <id> <local> <path> Upload a local file to a path on the server`;
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
			port:
				portArg !== undefined
					? parseInt(portArg, 10)
					: await gameServerService.getRandomUnusedGameServerPort(typeArg), // Aqui chamar o serviço que gera a porta evitando colisão
			rconPort:
				rconPortArg !== undefined
					? parseInt(rconPortArg, 10)
					: await gameServerService.getRandomUnusedRconPort(typeArg),
		};

		const id = await gameServerService.create(
			newGameServerData,
			config.minecraftServer.acceptEulaOnCreate,
		);
		console.log("Server created:", { ...newGameServerData, id });
	},

	ls: async () => console.table(await gameServerService.list()),

	info: async () => {
		const serverIdArg = getNextStdinArg()!;
		console.table(await gameServerService.get(serverIdArg));
	},

	stats: async () => {
		const serverIdArg = getNextStdinArg()!;
		console.table(await gameServerService.stats(serverIdArg));
	},

	logs: async () => {
		const serverIdArg = getNextStdinArg()!;
		const linesArg = getNextStdinArg();
		const lines = await gameServerService.logs(
			serverIdArg,
			linesArg ? Number(linesArg) : undefined,
		);
		console.log(lines.join("\n"));
	},

	// Control
	start: async () => {
		const serverIdArg = getNextStdinArg()!;
		await gameServerService.start(serverIdArg);
	},

	stop: async () => {
		const serverIdArg = getNextStdinArg()!;
		await gameServerService.stop(serverIdArg);
	},

	restart: async () => {
		const serverIdArg = getNextStdinArg()!;
		await gameServerService.restart(serverIdArg);
	},

	kill: async () => {
		const serverIdArg = getNextStdinArg()!;
		await gameServerService.kill(serverIdArg);
	},

	eula: async () => {
		const serverIdArg = getNextStdinArg()!;
		await gameServerService.acceptEula(serverIdArg);
	},

	"stop-all": async () => {
		await gameServerService.stopAll();
	},

	// Maintenance
	backup: async () => {
		const serverIdArg = getNextStdinArg()!;
		await gameServerService.backup(serverIdArg);
	},

	update: async () => {
		const serverIdArg = getNextStdinArg()!;
		await gameServerService.update(serverIdArg);
	},

	clone: async () => {
		const serverIdArg = getNextStdinArg()!;
		await gameServerService.clone(serverIdArg);
	},

	cmd: async () => {
		const serverIdArg = getNextStdinArg()!;
		const consoleCommandArg = Bun.argv.slice(argCounter).join(" ");
		console.log(
			await gameServerService.sendCommand(serverIdArg, consoleCommandArg),
		);
	},

	delete: async () => {
		const serverIdArg = getNextStdinArg()!;
		const removeFilesArg = getNextStdinArg() === "--files";
		await gameServerService.delete(serverIdArg, removeFilesArg);
	},

	"delete-all": async () => {
		const removeFilesArg = getNextStdinArg() === "--files";
		const servers = await gameServerService.list();
		for (const server of servers) {
			await gameServerService.delete(server.id, removeFilesArg);
			console.log("Server deleted:", server.id, server.name);
		}
		console.log(`Deleted ${servers.length} server(s)`);
	},

	// Files
	files: async () => {
		const subcommandArg = getNextStdinArg() as FilesSubcommand;
		const serverIdArg = getNextStdinArg()!;

		const subcommands: Record<FilesSubcommand, () => Promise<void>> = {
			ls: async () => {
				const pathArg = getNextStdinArg() ?? ".";
				console.table(await gameServerService.listFiles(serverIdArg, pathArg));
			},

			cat: async () => {
				const pathArg = getNextStdinArg()!;
				const file = await gameServerService.readFile(serverIdArg, pathArg);
				process.stdout.write(file.content);
			},

			write: async () => {
				const pathArg = getNextStdinArg()!;
				const file = await gameServerService.readFile(serverIdArg, pathArg);
				const content = await Bun.stdin.text();
				await gameServerService.writeFile(serverIdArg, { ...file, content });
			},

			touch: async () => {
				const pathArg = getNextStdinArg()!;
				await gameServerService.createFile(serverIdArg, pathArg);
			},

			mkdir: async () => {
				const pathArg = getNextStdinArg()!;
				await gameServerService.createDirectory(serverIdArg, pathArg);
			},

			rename: async () => {
				const pathArg = getNextStdinArg()!;
				const newNameArg = getNextStdinArg()!;
				await gameServerService.renameFile(serverIdArg, pathArg, newNameArg);
			},

			mv: async () => {
				const targetDirectoryArg = getNextStdinArg()!;
				const pathsArg = Bun.argv.slice(argCounter);
				await gameServerService.moveFiles(
					serverIdArg,
					pathsArg,
					targetDirectoryArg,
				);
			},

			cp: async () => {
				const targetDirectoryArg = getNextStdinArg()!;
				const pathsArg = Bun.argv.slice(argCounter);
				await gameServerService.copyFiles(
					serverIdArg,
					pathsArg,
					targetDirectoryArg,
				);
			},

			rm: async () => {
				const pathsArg = Bun.argv.slice(argCounter);
				await gameServerService.deleteFiles(serverIdArg, pathsArg);
			},

			download: async () => {
				const serverPathArg = getNextStdinArg()!;
				const downloadPathArg = getNextStdinArg();

				const { contentStream, fileNameOnServer } =
					await gameServerService.downloadFile(serverIdArg, serverPathArg);

				const destFile = Bun.file(
					downloadPathArg || fileNameOnServer || path.basename(serverPathArg),
				);

				if (await destFile.exists()) {
					console.error(`File ${destFile.name} already exists`);
					process.exit(1);
				}

				const writer = destFile.writer();

				for await (const chunk of contentStream) writer.write(chunk);

				await writer.end();
			},

			upload: async () => {
				const localFilePathArg = getNextStdinArg()!;
				const pathOnServerArg =
					getNextStdinArg() || path.basename(localFilePathArg);

				const localFilePointer = Bun.file(localFilePathArg);

				await gameServerService.uploadFile(
					serverIdArg,
					pathOnServerArg,
					localFilePointer,
				);
			},
		};

		const subcommand = subcommands[subcommandArg];
		if (subcommand === undefined) {
			console.log(`Unknown files command: ${subcommandArg}`);
			cliHelp();
			process.exit(1);
		}

		await subcommand();
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
