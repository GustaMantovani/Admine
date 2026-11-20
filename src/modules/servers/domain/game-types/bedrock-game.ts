import { InvalidServerPorts } from "../errors";
import type { ServerPorts } from "../port-purpose";
import type { GameProperties, GameType } from "./game-type";
import { mergeIntoPropertiesFile } from "./properties-file";

export class BedrockGame implements GameType {
	readonly id = "bedrock";
	readonly ports = [
		{
			purpose: "game",
			protocol: "udp",
			mode: "allocated",
			defaultRange: { min: 19132, max: 19999 },
		},
		{
			purpose: "gameV6",
			protocol: "udp",
			mode: "allocated",
			defaultRange: { min: 19132, max: 19999 },
		},
	] as const;
	readonly accepts = { version: false, memory: false };
	readonly requiresEulaAcceptance = false;
	readonly configFilePath = "server.properties";

	portProperties(ports: ServerPorts): GameProperties {
		if (ports.gameV6 === undefined)
			throw new InvalidServerPorts("Bedrock servers need an IPv6 port");
		return { "server-port": ports.game, "server-portv6": ports.gameV6 };
	}

	mergeConfig(
		existingContent: string | undefined,
		config: GameProperties,
	): string {
		return mergeIntoPropertiesFile(existingContent, config);
	}
}
