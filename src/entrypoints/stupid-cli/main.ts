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
	| "unknow";
const command: CLICommand = getNextStdinArg() as CLICommand | "unknow";

const craftyApiClient = new CraftyApiClient(
	CRAFTY_API_BASE_PATH,
	CRAFTY_API_TOKEN,
);
const craftyProvider = new CraftyMinecraftServerProvider(craftyApiClient);
const minecraftServerUseCases = new MinecraftServerUseCases(craftyProvider);

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

	unknow: async () => console.log(`Unknown command: ${command}`),
};

await commands[command]();
