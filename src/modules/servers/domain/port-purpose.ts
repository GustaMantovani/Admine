export type PortPurpose = "game" | "gameV6" | "rcon";

export type PortProtocol = "tcp" | "udp";

export type PortRange = { min: number; max: number };

export type ServerPorts = { game: number } & Partial<
	Record<Exclude<PortPurpose, "game">, number>
>;
