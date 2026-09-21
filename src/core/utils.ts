export function requiredEnv(name: string): string {
	const value = Bun.env[name];

	if (value === undefined)
		throw new Error(`Environment variable "${name}" is required`);

	return value;
}
