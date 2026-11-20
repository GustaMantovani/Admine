import type {
	PortProtocol,
	PortPurpose,
	PortRange,
	ServerPorts,
} from "../port-purpose";

export const GAME_TYPE_IDS = ["java", "bedrock", "hytale"] as const;
export type GameTypeId = (typeof GAME_TYPE_IDS)[number];

export type GameProperties = Record<string, string | number | boolean>;

export type GamePortSpec =
	| {
			purpose: PortPurpose;
			protocol: PortProtocol;
			mode: "allocated";
			defaultRange: PortRange;
	  }
	| {
			purpose: PortPurpose;
			protocol: PortProtocol;
			mode: "shared";
			port: number;
	  };

export interface GameType {
	readonly id: GameTypeId;
	readonly ports: readonly GamePortSpec[];
	readonly accepts: { readonly version: boolean; readonly memory: boolean };
	readonly requiresEulaAcceptance: boolean;
	readonly configFilePath: string;
	portProperties(ports: ServerPorts): GameProperties;
	mergeConfig(
		existingContent: string | undefined,
		config: GameProperties,
	): string;
}
