import { UnsupportedGameType } from "../domain/errors";
import type { GameType, GameTypeId } from "../domain/game-types/game-type";

export class GameTypeRegistry {
	private readonly types: ReadonlyMap<GameTypeId, GameType>;

	constructor(types: readonly GameType[]) {
		this.types = new Map(types.map((type) => [type.id, type]));
	}

	get(id: GameTypeId): GameType {
		const type = this.types.get(id);
		if (type === undefined)
			throw new UnsupportedGameType(`Game type ${id} is not registered`);
		return type;
	}

	list(): GameType[] {
		return [...this.types.values()];
	}

	sharedPorts(): Set<number> {
		return new Set(
			this.list().flatMap((type) =>
				type.ports.flatMap((spec) =>
					spec.mode === "shared" ? [spec.port] : [],
				),
			),
		);
	}
}
