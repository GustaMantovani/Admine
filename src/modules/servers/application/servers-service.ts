import type { BunFile } from "bun";
import type { TypeDefaults } from "../config";
import type { CreateServerSpec } from "../domain/create-server-spec";
import {
	CapabilityNotSupported,
	InvalidCreateOption,
	InvalidFile,
	InvalidServerName,
	NotApplicable,
	ServerNameTaken,
	SharedPortBusy,
	UnsupportedGameType,
} from "../domain/errors";
import type {
	CommandOutput,
	GameServerFile,
	GameServerFileDownload,
	GameServerFileEntry,
	GameServerStats,
	ServerDetails,
	ServerOverview,
	ServerRecord,
	ServerStatus,
} from "../domain/game-server";
import type {
	GameProperties,
	GameTypeId,
} from "../domain/game-types/game-type";
import type { PortPurpose, PortRange } from "../domain/port-purpose";
import {
	type FileCapability,
	hasBackup,
	hasFiles,
	hasUpdate,
	type OrchestratorTarget,
	type ServerOrchestrator,
} from "../ports/server-orchestrator";
import type { GameTypeRegistry } from "./game-type-registry";
import type { PortAllocator, PortRequest } from "./port-allocator";
import type { ServerRegistry } from "./server-registry";

export type CreateServerInput = {
	name?: string;
	type?: GameTypeId;
	version?: string;
};

export type PortRangeOverrides = Partial<
	Record<GameTypeId, Partial<Record<PortPurpose, PortRange>>>
>;

const SERVER_NAME = /^[a-z0-9]+(-[a-z0-9]+)*$/;

export class ServersService {
	constructor(
		private readonly orchestrator: ServerOrchestrator,
		private readonly registry: ServerRegistry,
		private readonly portAllocator: PortAllocator,
		private readonly gameTypes: GameTypeRegistry,
		private readonly defaults: Record<GameTypeId, TypeDefaults> & {
			type: GameTypeId;
		},
		private readonly rangeOverrides: PortRangeOverrides,
		private readonly generateName: () => string,
	) {}

	async create(input: CreateServerInput): Promise<ServerRecord> {
		const type = input.type ?? this.defaults.type;
		const game = this.gameTypes.get(type);
		const defaults = this.defaults[type];

		if (!this.orchestrator.supports(type))
			throw new UnsupportedGameType(
				`The orchestrator can't run ${type} servers`,
			);
		if (input.version !== undefined && !game.accepts.version)
			throw new InvalidCreateOption(`${type} servers don't accept a version`);

		const version = game.accepts.version
			? (input.version ?? defaults.version)
			: undefined;
		if (game.accepts.version && version === undefined)
			throw new InvalidCreateOption(`${type} servers need a version`);
		const memory = game.accepts.memory ? defaults.memory : undefined;

		const name = input.name ?? this.generateName();
		if (!SERVER_NAME.test(name))
			throw new InvalidServerName(
				`Server names use lowercase letters, numbers and hyphens: ${name}`,
			);
		if (await this.registry.isNameTaken(name))
			throw new ServerNameTaken(`There is already a server named ${name}`);

		const id = Bun.randomUUIDv7();
		const ports = await this.portAllocator.reserve(
			id,
			this.portRequestsFor(type),
		);

		try {
			const orchestratorProperties = this.orchestrator.propertiesFor(
				type,
				ports,
			);
			const spec = {
				type,
				name,
				ports,
				orchestratorProperties,
				...(version !== undefined && { version }),
				...(memory !== undefined && { memory }),
			} as CreateServerSpec;

			const ref = await this.orchestrator.create(spec);
			const record: ServerRecord = {
				id,
				name,
				type,
				version,
				memory,
				orchestratorRef: ref,
				gamePort: ports.game,
				createdAt: new Date().toISOString(),
			};

			try {
				await this.registry.add(record);
				await this.applyConfig(record, {
					...defaults.properties,
					...orchestratorProperties,
					...game.portProperties(ports),
				});
				if (game.requiresEulaAcceptance && defaults.acceptEulaOnCreate)
					await this.orchestrator.acceptEula({ ref, type });
			} catch (error) {
				await this.registry.remove(id);
				await this.orchestrator.delete({ ref, type }, true);
				throw error;
			}

			return record;
		} catch (error) {
			await this.portAllocator.release(id);
			throw error;
		}
	}

	async list(): Promise<ServerOverview[]> {
		const records = await this.registry.list();
		return Promise.all(records.map((record) => this.overview(record)));
	}

	async info(nameOrId: string): Promise<ServerDetails> {
		const record = await this.registry.get(nameOrId);
		const reservations = await this.portAllocator.reservationsOf(record.id);
		const ports = reservations.map(({ port, purpose }) => ({ port, purpose }));
		if (!ports.some((entry) => entry.purpose === "game"))
			ports.unshift({ port: record.gamePort, purpose: "game" });

		return { ...(await this.overview(record)), ports };
	}

	async stats(nameOrId: string): Promise<GameServerStats> {
		return this.orchestrator.stats(await this.target(nameOrId));
	}

	async logs(nameOrId: string, lastLines?: number): Promise<string[]> {
		const lines = await this.orchestrator.logs(await this.target(nameOrId));
		return lastLines === undefined ? lines : lines.slice(-lastLines);
	}

	async start(nameOrId: string): Promise<void> {
		const record = await this.registry.get(nameOrId);
		await this.ensureSharedPortFree(record);
		await this.orchestrator.start(toTarget(record));
	}

	async stop(nameOrId: string): Promise<void> {
		await this.orchestrator.stop(await this.target(nameOrId));
	}

	async restart(nameOrId: string): Promise<void> {
		const record = await this.registry.get(nameOrId);
		await this.ensureSharedPortFree(record);
		await this.orchestrator.restart(toTarget(record));
	}

	async kill(nameOrId: string): Promise<void> {
		await this.orchestrator.kill(await this.target(nameOrId));
	}

	async acceptEula(nameOrId: string): Promise<void> {
		const record = await this.registry.get(nameOrId);
		if (!this.gameTypes.get(record.type).requiresEulaAcceptance)
			throw new NotApplicable(
				`${record.type} servers don't need to accept the EULA`,
			);
		await this.orchestrator.acceptEula(toTarget(record));
	}

	async sendCommand(nameOrId: string, command: string): Promise<CommandOutput> {
		return this.orchestrator.sendCommand(await this.target(nameOrId), command);
	}

	async backup(nameOrId: string): Promise<void> {
		const orchestrator = this.orchestrator;
		if (!hasBackup(orchestrator)) throw notSupported("backup");
		await orchestrator.backup(await this.target(nameOrId));
	}

	async update(nameOrId: string): Promise<void> {
		const orchestrator = this.orchestrator;
		if (!hasUpdate(orchestrator)) throw notSupported("update");
		await orchestrator.update(await this.target(nameOrId));
	}

	async delete(nameOrId: string, removeFiles: boolean): Promise<ServerRecord> {
		const record = await this.registry.get(nameOrId);
		await this.orchestrator.delete(toTarget(record), removeFiles);
		await this.registry.remove(record.id);
		await this.portAllocator.release(record.id);
		return record;
	}

	async deleteAll(removeFiles: boolean): Promise<ServerRecord[]> {
		const deleted: ServerRecord[] = [];
		for (const record of await this.registry.list())
			deleted.push(await this.delete(record.id, removeFiles));
		return deleted;
	}

	async stopAll(): Promise<number> {
		return this.orchestrator.stopAll();
	}

	async listFiles(
		nameOrId: string,
		path: string,
	): Promise<GameServerFileEntry[]> {
		return this.files().listFiles(await this.target(nameOrId), path);
	}

	async readFile(nameOrId: string, path: string): Promise<GameServerFile> {
		return this.files().readFile(await this.target(nameOrId), path);
	}

	async writeFile(nameOrId: string, file: GameServerFile): Promise<void> {
		return this.files().writeFile(await this.target(nameOrId), file);
	}

	async createFile(nameOrId: string, path: string): Promise<void> {
		return this.files().createFile(await this.target(nameOrId), path);
	}

	async createDirectory(nameOrId: string, path: string): Promise<void> {
		return this.files().createDirectory(await this.target(nameOrId), path);
	}

	async renameFile(
		nameOrId: string,
		path: string,
		newName: string,
	): Promise<void> {
		return this.files().renameFile(await this.target(nameOrId), path, newName);
	}

	async moveFiles(
		nameOrId: string,
		sourcePaths: string[],
		targetDirectory: string,
	): Promise<void> {
		return this.files().moveFiles(
			await this.target(nameOrId),
			sourcePaths,
			targetDirectory,
		);
	}

	async copyFiles(
		nameOrId: string,
		sourcePaths: string[],
		targetDirectory: string,
	): Promise<void> {
		return this.files().copyFiles(
			await this.target(nameOrId),
			sourcePaths,
			targetDirectory,
		);
	}

	async deleteFiles(nameOrId: string, paths: string[]): Promise<void> {
		return this.files().deleteFiles(await this.target(nameOrId), paths);
	}

	async downloadFile(
		nameOrId: string,
		path: string,
	): Promise<GameServerFileDownload> {
		return this.files().downloadFile(await this.target(nameOrId), path);
	}

	async uploadFile(
		nameOrId: string,
		path: string,
		file: BunFile,
	): Promise<void> {
		if (!(await file.exists()))
			throw new InvalidFile(`File ${file.name} doesn't exist`);
		if (file.size === 0) throw new InvalidFile(`File ${file.name} is empty`);
		return this.files().uploadFile(await this.target(nameOrId), path, file);
	}

	private async applyConfig(
		record: ServerRecord,
		config: GameProperties,
	): Promise<void> {
		if (Object.keys(config).length === 0) return;

		const game = this.gameTypes.get(record.type);
		const target = toTarget(record);
		const existing = await this.orchestrator.readConfigFile(
			target,
			game.configFilePath,
		);
		await this.orchestrator.writeConfigFile(
			target,
			game.configFilePath,
			game.mergeConfig(existing, config),
		);
	}

	private portRequestsFor(type: GameTypeId): PortRequest[] {
		const overrides = this.rangeOverrides[type] ?? {};
		const gamePorts = this.gameTypes.get(type).ports.map(
			(spec): PortRequest =>
				spec.mode === "shared"
					? spec
					: {
							purpose: spec.purpose,
							protocol: spec.protocol,
							mode: "allocated",
							range: overrides[spec.purpose] ?? spec.defaultRange,
						},
		);
		const orchestratorPorts = this.orchestrator.portsFor(type).map(
			(request): PortRequest => ({
				purpose: request.purpose,
				protocol: request.protocol,
				mode: "allocated",
				range: overrides[request.purpose] ?? request.defaultRange,
			}),
		);
		return [...gamePorts, ...orchestratorPorts];
	}

	private async ensureSharedPortFree(record: ServerRecord): Promise<void> {
		const shared = this.gameTypes
			.get(record.type)
			.ports.some((spec) => spec.purpose === "game" && spec.mode === "shared");
		if (!shared) return;

		const others = (await this.registry.list()).filter(
			(other) => other.id !== record.id && other.gamePort === record.gamePort,
		);
		for (const other of others) {
			if ((await this.statusOf(other)).status === "on")
				throw new SharedPortBusy(
					`${other.name} is already running on port ${record.gamePort}; stop it before starting ${record.name}`,
				);
		}
	}

	private async overview(record: ServerRecord): Promise<ServerOverview> {
		return { ...record, ...(await this.statusOf(record)) };
	}

	private async statusOf(
		record: ServerRecord,
	): Promise<{ status: ServerStatus; runningVersion?: string }> {
		try {
			const stats = await this.orchestrator.stats(toTarget(record));
			return {
				status: stats.running ? "on" : "off",
				runningVersion: stats.version,
			};
		} catch {
			return { status: "unknown" };
		}
	}

	private async target(nameOrId: string): Promise<OrchestratorTarget> {
		return toTarget(await this.registry.get(nameOrId));
	}

	private files(): ServerOrchestrator & FileCapability {
		const orchestrator = this.orchestrator;
		if (!hasFiles(orchestrator)) throw notSupported("files");
		return orchestrator;
	}
}

function toTarget(record: ServerRecord): OrchestratorTarget {
	return { ref: record.orchestratorRef, type: record.type };
}

function notSupported(capability: string): CapabilityNotSupported {
	return new CapabilityNotSupported(
		`The orchestrator doesn't support ${capability} operations`,
	);
}
