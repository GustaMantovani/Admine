import { z } from "zod";
import { CraftyConfigSchema } from "./adapters/crafty/crafty-config";
import type { ServerMemory } from "./domain/create-server-spec";
import {
	GAME_TYPE_IDS,
	type GameProperties,
	type GameTypeId,
} from "./domain/game-types/game-type";
import type { PortPurpose, PortRange } from "./domain/port-purpose";

const PropertiesSchema = z.record(
	z.string(),
	z.union([z.string(), z.number(), z.boolean()]),
);

const MemorySchema = {
	minimumMemory: z.number().positive().default(2),
	maximumMemory: z.number().positive().default(4),
};

const DefaultsSchema = z
	.object({
		type: z.enum(GAME_TYPE_IDS).default("java"),
		java: z
			.object({
				version: z.string().default("26.3"),
				...MemorySchema,
				acceptEulaOnCreate: z.boolean().default(true),
				properties: PropertiesSchema.default({
					"white-list": false,
					"online-mode": false,
				}),
			})
			.prefault({}),
		bedrock: z
			.object({ properties: PropertiesSchema.default({}) })
			.prefault({}),
		hytale: z
			.object({ ...MemorySchema, properties: PropertiesSchema.default({}) })
			.prefault({}),
	})
	.prefault({});

const PortSchema = z.number().int().min(1).max(65535);

const RangeSchema = z
	.object({ min: PortSchema, max: PortSchema })
	.refine((range) => range.min <= range.max, "min must be <= max");

const TypeRangesSchema = z.partialRecord(
	z.enum(["game", "gameV6", "rcon"] satisfies PortPurpose[]),
	RangeSchema,
);

const PortsSchema = z
	.object({
		probe: z.boolean().default(true),
		java: TypeRangesSchema.optional(),
		bedrock: TypeRangesSchema.optional(),
		hytale: TypeRangesSchema.optional(),
	})
	.prefault({});

export const ServersConfigSchema = z.object({
	defaults: DefaultsSchema,
	ports: PortsSchema,
	orchestrator: CraftyConfigSchema.prefault({ kind: "crafty" }),
});

export type ServersConfig = z.infer<typeof ServersConfigSchema>;

export type TypeDefaults = {
	version?: string;
	memory?: ServerMemory;
	acceptEulaOnCreate: boolean;
	properties: GameProperties;
};

export function toTypeDefaults(
	defaults: ServersConfig["defaults"],
): Record<GameTypeId, TypeDefaults> {
	return {
		java: {
			version: defaults.java.version,
			memory: {
				minimum: defaults.java.minimumMemory,
				maximum: defaults.java.maximumMemory,
			},
			acceptEulaOnCreate: defaults.java.acceptEulaOnCreate,
			properties: defaults.java.properties,
		},
		bedrock: {
			acceptEulaOnCreate: false,
			properties: defaults.bedrock.properties,
		},
		hytale: {
			memory: {
				minimum: defaults.hytale.minimumMemory,
				maximum: defaults.hytale.maximumMemory,
			},
			acceptEulaOnCreate: false,
			properties: defaults.hytale.properties,
		},
	};
}

export function configuredRange(
	ports: ServersConfig["ports"],
	type: GameTypeId,
	purpose: PortPurpose,
): PortRange | undefined {
	return ports[type]?.[purpose];
}
