CREATE TABLE servers_records (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	type TEXT NOT NULL,
	version TEXT,
	min_memory INTEGER,
	max_memory INTEGER,
	orchestrator_ref TEXT NOT NULL UNIQUE,
	game_port INTEGER NOT NULL,
	created_at TEXT NOT NULL
);

CREATE TABLE servers_port_reservations (
	port INTEGER PRIMARY KEY CHECK (port BETWEEN 1 AND 65535),
	server_id TEXT NOT NULL,
	purpose TEXT NOT NULL
);

CREATE INDEX servers_port_reservations_server_id ON servers_port_reservations (server_id);
