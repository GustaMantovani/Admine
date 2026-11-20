import type { GameProperties } from "./game-types/game-type";
import type { ServerPorts } from "./port-purpose";

export type ServerMemory = { minimum: number; maximum: number };

type BaseSpec = {
	name: string;
	ports: ServerPorts;
	orchestratorProperties: GameProperties;
};

export type CreateServerSpec =
	| (BaseSpec & { type: "java"; version: string; memory: ServerMemory })
	| (BaseSpec & { type: "bedrock" })
	| (BaseSpec & { type: "hytale"; memory: ServerMemory });
