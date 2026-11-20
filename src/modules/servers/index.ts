import type { ModuleMigrations } from "../../shared/db/migration-runner";
import recordsMigration from "./migrations/001_records.sql" with {
	type: "text",
};
import craftyMigration from "./migrations/002_crafty.sql" with { type: "text" };

export const serversMigrations: ModuleMigrations = {
	module: "servers",
	migrations: [recordsMigration, craftyMigration],
};

export {
	CraftyApiClient,
	CraftyApiError,
} from "./adapters/crafty/crafty-api-client";
export {
	type CraftyAuthProvider,
	CraftyAuthTokenRepository,
	ManagedCraftyAuthProvider,
	resolveCraftyCredentials,
	StaticCraftyAuthProvider,
} from "./adapters/crafty/crafty-auth";
export { ensureCraftyInfra } from "./adapters/crafty/crafty-infra";
export { CraftyRconRepository } from "./adapters/crafty/crafty-rcon-repository";
export { CraftyServerOrchestrator } from "./adapters/crafty/crafty-server-orchestrator";
export { NetPortProbe } from "./adapters/net/net-port-probe";
export { SqlitePortReservationRepository } from "./adapters/sqlite/sqlite-port-reservation-repository";
export { SqliteServerRepository } from "./adapters/sqlite/sqlite-server-repository";

export { GameTypeRegistry } from "./application/game-type-registry";
export { PortAllocator } from "./application/port-allocator";
export { ServerRegistry } from "./application/server-registry";
export {
	type CreateServerInput,
	type PortRangeOverrides,
	ServersService,
} from "./application/servers-service";

export {
	type ServersConfig,
	ServersConfigSchema,
	toTypeDefaults,
} from "./config";

export * from "./domain/errors";
export type * from "./domain/game-server";
export { BedrockGame } from "./domain/game-types/bedrock-game";
export {
	GAME_TYPE_IDS,
	type GameType,
	type GameTypeId,
} from "./domain/game-types/game-type";
export { HytaleGame } from "./domain/game-types/hytale-game";
export { JavaGame } from "./domain/game-types/java-game";
