import { mkdir } from "node:fs/promises";
import path from "node:path";
import {
	adjectives,
	animals,
	uniqueNamesGenerator,
} from "unique-names-generator";
import { z } from "zod";
import { LibRconClient } from "../lib/rcon";
import {
	BedrockGame,
	CraftyApiClient,
	type CraftyAuthProvider,
	CraftyAuthTokenRepository,
	CraftyRconRepository,
	CraftyServerOrchestrator,
	ensureCraftyInfra,
	GameTypeRegistry,
	HytaleGame,
	JavaGame,
	ManagedCraftyAuthProvider,
	NetPortProbe,
	PortAllocator,
	resolveCraftyCredentials,
	ServerRegistry,
	ServersConfigSchema,
	ServersService,
	SqlitePortReservationRepository,
	SqliteServerRepository,
	StaticCraftyAuthProvider,
	serversMigrations,
	toTypeDefaults,
} from "../modules/servers";
import { loadConfigFile } from "../shared/config";
import { openDatabase } from "../shared/db/database";
import { runMigrations } from "../shared/db/migration-runner";

export const DEFAULT_CONFIG_PATH = "./config/admine.toml";

const AppConfigSchema = z.object({
	dataDir: z.string().default(".data"),
	servers: ServersConfigSchema.prefault({}),
});

export type App = {
	servers: ServersService;
};

export async function buildApp(configPath = DEFAULT_CONFIG_PATH): Promise<App> {
	const config = AppConfigSchema.parse(await loadConfigFile(configPath));

	await mkdir(config.dataDir, { recursive: true });
	const db = openDatabase(path.join(config.dataDir, "admine.sqlite"));
	runMigrations(db, [serversMigrations]);

	const crafty = config.servers.orchestrator;
	let auth: CraftyAuthProvider;
	if (crafty.deployment.type === "managed") {
		await ensureCraftyInfra(crafty.deployment);
		auth = new ManagedCraftyAuthProvider(
			crafty.apiBaseUrl,
			await resolveCraftyCredentials(
				crafty.deployment,
				path.join(config.dataDir, "crafty/config/default-creds.txt"),
			),
			new CraftyAuthTokenRepository(db),
		);
	} else {
		auth = new StaticCraftyAuthProvider(crafty.deployment);
	}
	const orchestrator = new CraftyServerOrchestrator(
		new CraftyApiClient(crafty.apiBaseUrl, auth),
		crafty,
		new CraftyRconRepository(db),
		new LibRconClient(),
	);

	const gameTypes = new GameTypeRegistry([
		new JavaGame(),
		new BedrockGame(),
		new HytaleGame(),
	]);
	const portAllocator = new PortAllocator(
		new SqlitePortReservationRepository(db),
		config.servers.ports.probe ? new NetPortProbe() : undefined,
		gameTypes.sharedPorts(),
	);
	const { probe: _, ...rangeOverrides } = config.servers.ports;

	const servers = new ServersService(
		orchestrator,
		new ServerRegistry(new SqliteServerRepository(db)),
		portAllocator,
		gameTypes,
		{
			...toTypeDefaults(config.servers.defaults),
			type: config.servers.defaults.type,
		},
		rangeOverrides,
		() =>
			uniqueNamesGenerator({
				dictionaries: [adjectives, animals],
				separator: "-",
				length: 2,
			}),
	);

	return { servers };
}
