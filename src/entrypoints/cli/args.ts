let argCounter = 2;

export function getNextStdinArg(): string | undefined {
	return Bun.argv[argCounter++];
}

export function getRemainingArgs(): string[] {
	return Bun.argv.slice(argCounter);
}
