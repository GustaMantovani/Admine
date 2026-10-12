import { MinecraftServerService } from "../../core";
import { CraftyApiClient } from "../../core/providers/crafty/crafty-api-client";
import { CraftyMinecraftServerProvider } from "../../core/providers/crafty/crafty-minecraft-server-provider";
import { ensureCraftyInfra } from "../../infra/crafty-docker";
import { open } from 'lmdb';
import { LmdbAdmineStorage } from "../../core/repositories/lmdb-admine-storage";
import { StaticCraftyAuthProvider } from "../../core/providers/crafty/auth/static-crafty-auth-provider";
import { ManagedCraftyAuthProvider } from "../../core/providers/crafty/auth/managed-crafty-auth-provider";
import type { CraftyAuthProvider } from "../../core/providers/crafty/auth/crafty-auth-provider";
import { getNextStdinArg } from "./args";
import { type CliCommand, cliHelp, createCommands } from "./commands";
import { loadCliConfig, withCraftyCreds } from "./config";

const config = await loadCliConfig();

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

const minecraftServerService = new MinecraftServerService(craftyProvider);

const commands = createCommands(minecraftServerService, config, command);

const handler = commands[command];

if (handler === undefined) {
  cliHelp();
  process.exit(1);
};

await handler();
