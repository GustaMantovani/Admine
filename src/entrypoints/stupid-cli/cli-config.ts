import { TOML } from "bun";
import { z } from "zod";
import type { CraftyConfig } from "../../core/minecraft/crafty/crafty-config";
import type { CraftyCredsFile } from "../../core/minecraft/crafty/dtos/crafty-creds-file";
import { MinecraftServerType } from "../../core/types/minecraft-server-type";
import {
	DEFAULT_CLI_CONFIG_FILE_PATH,
	DEFAULT_CRAFTY_API_BASE_URL,
	DEFAULT_CRAFTY_CREDS_FILE_PATH,
	DEFAULT_CRAFTY_INFRA_CHECK_COMMAND,
	DEFAULT_CRAFTY_INFRA_CHECK_DELAY_MS,
	DEFAULT_CRAFTY_INFRA_SETUP_COMMAND,
	DEFAULT_MINECRAFT_SERVER_MAXIMUM_MEMORY,
	DEFAULT_MINECRAFT_SERVER_MINIMUM_MEMORY,
	DEFAULT_MINECRAFT_SERVER_TYPE,
	DEFAULT_MINECRAFT_SERVER_VERSION,
	DEFAULT_SERVERS_HOST_FOR_ADMINE,
} from "./defaults";

const ManagedCraftySchema = z.object({
	deployType: z.literal("managed"),
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
	apiBaseUrl: z.string(),
	apiToken: z.string(),
});

const MinecraftServerDefaultsSchema = z.object({
	type: z.enum(MinecraftServerType).default(DEFAULT_MINECRAFT_SERVER_TYPE),
	version: z.string().default(DEFAULT_MINECRAFT_SERVER_VERSION),
	minimumMemory: z.number().default(DEFAULT_MINECRAFT_SERVER_MINIMUM_MEMORY),
	maximumMemory: z.number().default(DEFAULT_MINECRAFT_SERVER_MAXIMUM_MEMORY),
});

const CliConfigSchema = z.object({
	serversHost: z.string().default(DEFAULT_SERVERS_HOST_FOR_ADMINE),
	minecraftServer: MinecraftServerDefaultsSchema.prefault({}),
	craftyConfig: z
		.discriminatedUnion("deployType", [
			ManagedCraftySchema,
			SelfHostedCraftySchema,
		])
		.prefault({ deployType: "managed" }),
});

export type CliConfig = {
	serversHost: string;
	minecraftServer: z.infer<typeof MinecraftServerDefaultsSchema>;
	craftyConfig: CraftyConfig;
};

export async function loadCliConfig(
	configFilePath = DEFAULT_CLI_CONFIG_FILE_PATH,
): Promise<CliConfig> {
	const file = Bun.file(configFilePath);
	// A missing file means the user wants every default
	const raw = (await file.exists()) ? TOML.parse(await file.text()) : {};

	const { craftyConfig, ...rest } = CliConfigSchema.parse(raw);

	if (craftyConfig.deployType === "self-hosted") {
		return { ...rest, craftyConfig };
	}

	const needsCreds = !craftyConfig.username || !craftyConfig.password;
	const creds = needsCreds ? await loadDefaultCreds() : undefined;

	return {
		...rest,
		craftyConfig: {
			...craftyConfig,
			username: craftyConfig.username ?? creds?.username ?? "",
			password: craftyConfig.password ?? creds?.password ?? "",
		},
	};
}

async function loadDefaultCreds(
	credsPath = DEFAULT_CRAFTY_CREDS_FILE_PATH,
): Promise<CraftyCredsFile> {
	return (await Bun.file(credsPath).json()) as CraftyCredsFile;
}
