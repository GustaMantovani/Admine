import type { ServerPorts } from "../port-purpose";
import type { GameProperties, GameType } from "./game-type";
import { mergeIntoPropertiesFile } from "./properties-file";

export class JavaGame implements GameType {
	readonly id = "java";
	readonly ports = [
		{
			purpose: "game",
			protocol: "tcp",
			mode: "allocated",
			defaultRange: { min: 25565, max: 25999 },
		},
	] as const;
	readonly accepts = { version: true, memory: true };
	readonly requiresEulaAcceptance = true;
	readonly configFilePath = "server.properties";

	portProperties(ports: ServerPorts): GameProperties {
		return { "server-port": ports.game };
	}

	mergeConfig(
		existingContent: string | undefined,
		config: GameProperties,
	): string {
		return mergeIntoPropertiesFile(existingContent, config);
	}
}
