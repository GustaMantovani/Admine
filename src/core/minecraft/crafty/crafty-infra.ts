import { $ } from "bun";
import type { CraftyInfraConfig } from "./crafty-config";

export async function ensureCraftyInfra(
	config: CraftyInfraConfig,
): Promise<void> {
	const firstInfraCheckExitCode = (
		await $`${{ raw: config.infraCheckCommand }}`.nothrow().quiet()
	).exitCode;
	if (firstInfraCheckExitCode === 0) return;

	await $`${{ raw: config.infraSetupCommand }}`;

	Bun.sleepSync(config.infraCheckDelay);

	await $`${{ raw: config.infraCheckCommand }}`;
}
