import type { CraftyConfig, ManagedCraftyConfig } from "../../core/minecraft/crafty/crafty-config"
import type { CraftyCredsFile } from "../../core/minecraft/crafty/dtos/crafty-creds-file";

export type CliConfig ={
	craftyConfig: CraftyConfig;
}

export async function loadDefaultCreds(
    craftyConfig: ManagedCraftyConfig,
    credsPath = ".data/crafty/config/default-creds.txt",
): Promise<void> {
  const file = Bun.file(credsPath);
	const creds = (await file.json()) as CraftyCredsFile;

  craftyConfig.username = creds.username;
  craftyConfig.password = creds.password;
}
