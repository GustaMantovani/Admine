import { $ } from "bun";

export type CraftyInfraConfig = {
	infraSetupCommand: string;
	infraCheckDelay: number;
	infraCheckCommand: string;
};

export async function ensureCraftyInfra(
	config: CraftyInfraConfig,
): Promise<void> {
	const firstInfraCheckExitCode = (
		await $`${{ raw: config.infraCheckCommand }}`.nothrow().quiet()
	).exitCode;
	if (firstInfraCheckExitCode === 0) return;

	await $`${{ raw: config.infraSetupCommand }}`.quiet();

	await Bun.sleep(config.infraCheckDelay);

	await $`${{ raw: config.infraCheckCommand }}`.quiet();
}
