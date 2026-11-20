import { $ } from "bun";
import type { ManagedCraftyDeployment } from "./crafty-config";

export async function ensureCraftyInfra(
	deployment: ManagedCraftyDeployment,
): Promise<void> {
	const check = await $`${{ raw: deployment.infraCheckCommand }}`
		.nothrow()
		.quiet();
	if (check.exitCode === 0) return;

	await $`${{ raw: deployment.infraSetupCommand }}`.quiet();
	await Bun.sleep(deployment.infraCheckDelay);
	await $`${{ raw: deployment.infraCheckCommand }}`.quiet();
}
