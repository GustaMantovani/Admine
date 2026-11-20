import type { Database } from "bun:sqlite";
import type { PortPurpose } from "../../domain/port-purpose";
import type {
	PortReservation,
	PortReservationRepository,
} from "../../ports/port-reservation-repository";

type Row = { port: number; server_id: string; purpose: string };

export class SqlitePortReservationRepository
	implements PortReservationRepository
{
	constructor(private readonly db: Database) {}

	async reservedPorts(): Promise<Set<number>> {
		const rows = this.db
			.query("SELECT port FROM servers_port_reservations")
			.all() as {
			port: number;
		}[];
		return new Set(rows.map((row) => row.port));
	}

	async tryReserve(reservation: PortReservation): Promise<boolean> {
		const result = this.db
			.query(
				"INSERT OR IGNORE INTO servers_port_reservations (port, server_id, purpose) VALUES (?, ?, ?)",
			)
			.run(reservation.port, reservation.serverId, reservation.purpose);
		return result.changes === 1;
	}

	async findByServer(serverId: string): Promise<PortReservation[]> {
		const rows = this.db
			.query(
				"SELECT * FROM servers_port_reservations WHERE server_id = ? ORDER BY port",
			)
			.all(serverId) as Row[];
		return rows.map((row) => ({
			port: row.port,
			serverId: row.server_id,
			purpose: row.purpose as PortPurpose,
		}));
	}

	async releaseServer(serverId: string): Promise<void> {
		this.db
			.query("DELETE FROM servers_port_reservations WHERE server_id = ?")
			.run(serverId);
	}
}
