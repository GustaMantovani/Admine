import type { PortPurpose } from "../domain/port-purpose";

export type PortReservation = {
	port: number;
	serverId: string;
	purpose: PortPurpose;
};

export interface PortReservationRepository {
	reservedPorts(): Promise<Set<number>>;
	tryReserve(reservation: PortReservation): Promise<boolean>;
	findByServer(serverId: string): Promise<PortReservation[]>;
	releaseServer(serverId: string): Promise<void>;
}
