import path from "node:path";
import { z } from "zod";
import { GAME_TYPE_IDS } from "../modules/servers";
import { buildApp } from "./composition-root";

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
	| "eula"
	| "cmd"
	| "delete"
	| "delete-all"
	| "stop-all"
	| "files"
	| "help";

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

const command = getNextStdinArg() as CliCommand | undefined;

function cliHelp() {
	console.log(`Usage: bun run cli <command> [args]

    Servers (<server> is the server's name or id):
      create [name] [type] [version]   Create a server (type: ${GAME_TYPE_IDS.join(" | ")}, default from config)
      ls                               List servers
      info <server>                    Show a server's details and ports
      stats <server>                   Show CPU, memory, players and state of a server
      logs <server> [lines]            Show the server console (last N lines)

    Control:
      start <server>                   Start the server
      stop <server>                    Stop the server
      restart <server>                 Restart the server
      kill <server>                    Force-kill the server process
      eula <server>                    Accept the EULA (Java servers won't start without it)
      stop-all                         Stop every running server

    Maintenance:
      backup <server>                  Start a backup of the server
      update <server>                  Update the server executable
      cmd <server> <command...>        Send a command to the server console
      delete <server> [--files]        Delete the server (--files also deletes its files)
      delete-all [--files]             Delete every server (--files also deletes their files)

    Files (paths are relative to the server's root directory):
      files ls <server> [path]             List a directory (default: the server's root)
      files cat <server> <path>            Print a file
      files write <server> <path>          Replace an existing file's content with stdin
      files touch <server> <path>          Create an empty file
      files mkdir <server> <path>          Create a directory
      files rename <server> <path> <name>  Rename a file or directory
      files mv <server> <dir> <paths...>   Move files or directories into a directory
      files cp <server> <dir> <paths...>   Copy files or directories into a directory
      files rm <server> <paths...>         Delete files or directories
      files download <server> <path> [dest] Download a file to a local path
      files upload <server> <local> <path> Upload a local file to a path on the server`);
}

if (command === undefined || command === "help") {
	cliHelp();
	process.exit(0);
}

const { servers } = await buildApp();

const commands: Record<Exclude<CliCommand, "help">, () => Promise<void>> = {
	create: async () => {
		const name = getNextStdinArg();
		const type = z.enum(GAME_TYPE_IDS).optional().parse(getNextStdinArg());
		const version = getNextStdinArg();
		const created = await servers.create({ name, type, version });
		console.log("Server created:", created);
	},

	ls: async () => console.table(await servers.list()),

	info: async () => {
		const { ports, ...details } = await servers.info(getNextStdinArg()!);
		console.table(details);
		console.table(ports);
	},

	stats: async () => console.table(await servers.stats(getNextStdinArg()!)),

	logs: async () => {
		const server = getNextStdinArg()!;
		const lines = getNextStdinArg();
		console.log(
			(await servers.logs(server, lines ? Number(lines) : undefined)).join(
				"\n",
			),
		);
	},

	start: () => servers.start(getNextStdinArg()!),
	stop: () => servers.stop(getNextStdinArg()!),
	restart: () => servers.restart(getNextStdinArg()!),
	kill: () => servers.kill(getNextStdinArg()!),
	eula: () => servers.acceptEula(getNextStdinArg()!),

	"stop-all": async () => {
		console.log(`Stopped ${await servers.stopAll()} server(s)`);
	},

	backup: () => servers.backup(getNextStdinArg()!),
	update: () => servers.update(getNextStdinArg()!),

	cmd: async () => {
		const server = getNextStdinArg()!;
		const { output } = await servers.sendCommand(
			server,
			Bun.argv.slice(argCounter).join(" "),
		);
		console.log(output ?? "Command sent (this server doesn't return output)");
	},

	delete: async () => {
		const server = getNextStdinArg()!;
		const removeFiles = getNextStdinArg() === "--files";
		const deleted = await servers.delete(server, removeFiles);
		console.log("Server deleted:", deleted.name);
	},

	"delete-all": async () => {
		const removeFiles = getNextStdinArg() === "--files";
		const deleted = await servers.deleteAll(removeFiles);
		for (const server of deleted) console.log("Server deleted:", server.name);
		console.log(`Deleted ${deleted.length} server(s)`);
	},

	files: async () => {
		const subcommandArg = getNextStdinArg() as FilesSubcommand;
		const server = getNextStdinArg()!;

		const subcommands: Record<FilesSubcommand, () => Promise<void>> = {
			ls: async () =>
				console.table(
					await servers.listFiles(server, getNextStdinArg() ?? "."),
				),

			cat: async () => {
				const file = await servers.readFile(server, getNextStdinArg()!);
				process.stdout.write(file.content);
			},

			write: async () => {
				const file = await servers.readFile(server, getNextStdinArg()!);
				await servers.writeFile(server, {
					...file,
					content: await Bun.stdin.text(),
				});
			},

			touch: () => servers.createFile(server, getNextStdinArg()!),
			mkdir: () => servers.createDirectory(server, getNextStdinArg()!),

			rename: async () => {
				const filePath = getNextStdinArg()!;
				await servers.renameFile(server, filePath, getNextStdinArg()!);
			},

			mv: async () => {
				const targetDirectory = getNextStdinArg()!;
				await servers.moveFiles(
					server,
					Bun.argv.slice(argCounter),
					targetDirectory,
				);
			},

			cp: async () => {
				const targetDirectory = getNextStdinArg()!;
				await servers.copyFiles(
					server,
					Bun.argv.slice(argCounter),
					targetDirectory,
				);
			},

			rm: () => servers.deleteFiles(server, Bun.argv.slice(argCounter)),

			download: async () => {
				const serverPath = getNextStdinArg()!;
				const downloadPath = getNextStdinArg();
				const { contentStream, fileNameOnServer } = await servers.downloadFile(
					server,
					serverPath,
				);

				const destFile = Bun.file(
					downloadPath || fileNameOnServer || path.basename(serverPath),
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
				const localFilePath = getNextStdinArg()!;
				const pathOnServer = getNextStdinArg() || path.basename(localFilePath);
				await servers.uploadFile(server, pathOnServer, Bun.file(localFilePath));
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
};

const handler = commands[command];
if (handler === undefined) {
	console.log(`Unknown command: ${command}`);
	cliHelp();
	process.exit(1);
}

await handler();
