import { InvalidServerPorts, PortRangeExhausted } from "../domain/errors";
import type {
	PortProtocol,
	PortPurpose,
	PortRange,
	ServerPorts,
} from "../domain/port-purpose";
import type { PortProbe } from "../ports/port-probe";
import type {
	PortReservation,
	PortReservationRepository,
} from "../ports/port-reservation-repository";

export type PortRequest =
	| {
			purpose: PortPurpose;
			protocol: PortProtocol;
			mode: "allocated";
			range: PortRange;
	  }
	| {
			purpose: PortPurpose;
			protocol: PortProtocol;
			mode: "shared";
			port: number;
	  };

export class PortAllocator {
	constructor(
		private readonly reservations: PortReservationRepository,
		private readonly probe: PortProbe | undefined,
		private readonly sharedPorts: ReadonlySet<number>,
	) {}

	async reserve(
		serverId: string,
		requests: readonly PortRequest[],
	): Promise<ServerPorts> {
		const ports: Partial<Record<PortPurpose, number>> = {};
		try {
			for (const request of requests) {
				ports[request.purpose] =
					request.mode === "shared"
						? request.port
						: await this.reserveOne(serverId, request);
			}
		} catch (error) {
			await this.reservations.releaseServer(serverId);
			throw error;
		}

		if (ports.game === undefined)
			throw new InvalidServerPorts("Every server needs a game port");
		return ports as ServerPorts;
	}

	async release(serverId: string): Promise<void> {
		await this.reservations.releaseServer(serverId);
	}

	async reservationsOf(serverId: string): Promise<PortReservation[]> {
		return this.reservations.findByServer(serverId);
	}

	private async reserveOne(
		serverId: string,
		request: Extract<PortRequest, { mode: "allocated" }>,
	): Promise<number> {
		const reserved = await this.reservations.reservedPorts();

		for (let port = request.range.min; port <= request.range.max; port++) {
			if (reserved.has(port) || this.sharedPorts.has(port)) continue;
			if (this.probe && !(await this.probe.isFree(port, request.protocol)))
				continue;
			if (
				await this.reservations.tryReserve({
					port,
					serverId,
					purpose: request.purpose,
				})
			)
				return port;
		}

		throw new PortRangeExhausted(
			`No free ${request.protocol.toUpperCase()} port for ${request.purpose} between ${request.range.min} and ${request.range.max}`,
		);
	}
}
