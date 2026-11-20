import { ServerNotFound } from "../domain/errors";
import type { ServerRecord } from "../domain/game-server";
import type { ServerRepository } from "../ports/server-repository";

export class ServerRegistry {
	constructor(private readonly repository: ServerRepository) {}

	async get(nameOrId: string): Promise<ServerRecord> {
		const record =
			(await this.repository.findByName(nameOrId)) ??
			(await this.repository.findById(nameOrId));
		if (record === undefined)
			throw new ServerNotFound(`Server ${nameOrId} not found`);
		return record;
	}

	async isNameTaken(name: string): Promise<boolean> {
		return (await this.repository.findByName(name)) !== undefined;
	}

	async list(): Promise<ServerRecord[]> {
		return this.repository.list();
	}

	async add(record: ServerRecord): Promise<void> {
		await this.repository.insert(record);
	}

	async remove(id: string): Promise<void> {
		await this.repository.delete(id);
	}
}
