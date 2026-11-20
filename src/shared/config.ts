import { TOML } from "bun";

export async function loadConfigFile(
	path: string,
): Promise<Record<string, unknown>> {
	const file = Bun.file(path);
	if (!(await file.exists())) return {};
	return TOML.parse(await file.text()) as Record<string, unknown>;
}
