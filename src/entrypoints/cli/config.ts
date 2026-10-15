import { randomBytes } from "node:crypto";
import { TOML } from "bun";
import { get, has, set } from "lodash";
import { z } from "zod";
import { GameServerType } from "../../core/domain/game-server-type";
import { GameServerWebhookTrigger } from "../../core/domain/game-server-webhook-trigger";
import type {
	ManagedCraftyConfig,
	SelfHostedCraftyConfig,
} from "../../core/providers/crafty/crafty-config";
import type { CraftyCredsFile } from "../../core/providers/crafty/dtos/crafty-creds-file";
import type { CraftyInfraConfig } from "../../infra/crafty-docker";
import {
	DEFAULT_CLI_CONFIG_FILE_PATH,
	DEFAULT_CRAFTY_API_BASE_URL,
	DEFAULT_CRAFTY_CREDS_FILE_PATH,
	DEFAULT_CRAFTY_DEPLOY_TYPE,
	DEFAULT_CRAFTY_INFRA_CHECK_COMMAND,
	DEFAULT_CRAFTY_INFRA_CHECK_DELAY_MS,
	DEFAULT_CRAFTY_INFRA_SETUP_COMMAND,
	DEFAULT_CRAFTY_WEBHOOK_URLS,
	DEFAULT_CRAFTY_WEBHOOKS_ENABLED,
	DEFAULT_MEINCRAFT_SERVER_PROPERTIES,
	DEFAULT_MINECRAFT_SERVER_ACCEPT_EULA_ON_CREATE,
	DEFAULT_MINECRAFT_SERVER_MAXIMUM_MEMORY,
	DEFAULT_MINECRAFT_SERVER_MINIMUM_MEMORY,
	DEFAULT_MINECRAFT_SERVER_RCON_ENABLED,
	DEFAULT_MINECRAFT_SERVER_TYPE,
	DEFAULT_MINECRAFT_SERVER_VERSION,
	DEFAULT_SERVERS_HOST_FOR_ADMINE,
} from "./defaults";

const CraftyWebhooksSchema = z
	.object({
		enabled: z.boolean().default(DEFAULT_CRAFTY_WEBHOOKS_ENABLED),
		urls: z
			.partialRecord(z.enum(GameServerWebhookTrigger), z.url())
			.default(DEFAULT_CRAFTY_WEBHOOK_URLS),
	})
	.prefault({});

const ManagedCraftySchema = z.object({
	deployType: z.literal("managed"),
	webhooks: CraftyWebhooksSchema,
	apiBaseUrl: z.string().default(DEFAULT_CRAFTY_API_BASE_URL),
	infraSetupCommand: z.string().default(DEFAULT_CRAFTY_INFRA_SETUP_COMMAND),
	infraCheckCommand: z.string().default(DEFAULT_CRAFTY_INFRA_CHECK_COMMAND),
	infraCheckDelay: z.number().default(DEFAULT_CRAFTY_INFRA_CHECK_DELAY_MS),
	// Read from the Crafty default creds file when not set
	username: z.string().optional(),
	password: z.string().optional(),
});

const SelfHostedCraftySchema = z.object({
	deployType: z.literal("self-hosted"),
	webhooks: CraftyWebhooksSchema,
	apiBaseUrl: z.string(),
	apiToken: z.string(),
});

const RconSchema = z
	.object({
		enabled: z.boolean().default(DEFAULT_MINECRAFT_SERVER_RCON_ENABLED),
		password: z.string().default(randomBytes(32).toString("base64url")),
	})
	.prefault({});

const MinecraftServerDefaultsSchema = z.object({
	rcon: RconSchema,
	type: z.enum(GameServerType).default(DEFAULT_MINECRAFT_SERVER_TYPE),
	version: z.string().default(DEFAULT_MINECRAFT_SERVER_VERSION),
	minimumMemory: z.number().default(DEFAULT_MINECRAFT_SERVER_MINIMUM_MEMORY),
	maximumMemory: z.number().default(DEFAULT_MINECRAFT_SERVER_MAXIMUM_MEMORY),
	properties: z
		.record(z.string(), z.union([z.string(), z.number(), z.boolean()]))
		.default(DEFAULT_MEINCRAFT_SERVER_PROPERTIES),
	acceptEulaOnCreate: z
		.boolean()
		.default(DEFAULT_MINECRAFT_SERVER_ACCEPT_EULA_ON_CREATE),
});

const CliConfigSchema = z.object({
	serversHost: z.string().default(DEFAULT_SERVERS_HOST_FOR_ADMINE),
	minecraftServer: MinecraftServerDefaultsSchema.prefault({}),
	craftyConfig: z.preprocess(
		(value) => ({ deployType: DEFAULT_CRAFTY_DEPLOY_TYPE, ...(value ?? {}) }),
		z.discriminatedUnion("deployType", [
			ManagedCraftySchema,
			SelfHostedCraftySchema,
		]),
	),
});

// Credentials are only known after the Crafty infra is up, see withCraftyCreds
export type CliManagedCraftyConfig = Omit<
	ManagedCraftyConfig,
	"username" | "password"
> &
	CraftyInfraConfig & { username?: string; password?: string };

export type CliConfig = {
	serversHost: string;
	minecraftServer: z.infer<typeof MinecraftServerDefaultsSchema>;
	craftyConfig: SelfHostedCraftyConfig | CliManagedCraftyConfig;
};

export async function loadCliConfig(
	configFilePath = DEFAULT_CLI_CONFIG_FILE_PATH,
): Promise<CliConfig> {
	const file = Bun.file(configFilePath);
	// A missing file means the user wants every default
	const raw = (await file.exists()) ? TOML.parse(await file.text()) : {};
	const parsedConfig = CliConfigSchema.parse(raw);

	// The generated password must survive restarts, so persist it once
	const rconPasswordPath = "minecraftServer.rcon.password";
	if (!has(raw, rconPasswordPath)) {
		await writeSinglePropertyToTomlFile(
			parsedConfig,
			configFilePath,
			rconPasswordPath,
		);
	}

	return parsedConfig;
}

// Must run after the Crafty infra is up: the creds file is created by Crafty itself
export async function withCraftyCreds(
	craftyConfig: CliManagedCraftyConfig,
	credsPath = DEFAULT_CRAFTY_CREDS_FILE_PATH,
): Promise<ManagedCraftyConfig> {
	const file = Bun.file(credsPath);
	const creds = (await file.json()) as CraftyCredsFile;
	return {
		...craftyConfig,
		username: craftyConfig.username ?? creds.username,
		password: craftyConfig.password ?? creds.password,
	};
}

// Copies one property into the TOML file, keeping whatever else it already holds
async function writeSinglePropertyToTomlFile(
	inputObj: Record<string, unknown>,
	filePath: string,
	keyPath: string,
): Promise<void> {
	if (!has(inputObj, keyPath)) {
		throw new Error(`Property "${keyPath}" was not found`);
	}

	const file = Bun.file(filePath);
	const outputObj: Record<string, unknown> = (await file.exists())
		? (TOML.parse(await file.text()) as Record<string, unknown>)
		: {};
	set(outputObj, keyPath, get(inputObj, keyPath));

	const toml = TOML.stringify(outputObj);
	if (toml === undefined) {
		throw new Error("Could not serialize the config to TOML");
	}

	await Bun.write(filePath, toml);
}
