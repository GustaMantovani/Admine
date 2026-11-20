import type { GameProperties, GameType } from "./game-type";

export class HytaleGame implements GameType {
	readonly id = "hytale";
	readonly ports = [
		{ purpose: "game", protocol: "udp", mode: "shared", port: 5520 },
	] as const;
	readonly accepts = { version: false, memory: true };
	readonly requiresEulaAcceptance = false;
	readonly configFilePath = "config.json";

	portProperties(): GameProperties {
		return {};
	}

	mergeConfig(
		existingContent: string | undefined,
		config: GameProperties,
	): string {
		const merged = existingContent
			? (JSON.parse(existingContent) as Record<string, unknown>)
			: {};

		for (const [key, value] of Object.entries(config)) {
			const path = key.split(".");
			const last = path.pop() as string;
			let node = merged;
			for (const part of path) {
				const child = node[part];
				node[part] = child !== null && typeof child === "object" ? child : {};
				node = node[part] as Record<string, unknown>;
			}
			node[last] = value;
		}

		return `${JSON.stringify(merged, null, 2)}\n`;
	}
}
