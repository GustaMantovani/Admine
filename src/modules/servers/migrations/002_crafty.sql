CREATE TABLE servers_crafty_auth_token (
	id INTEGER PRIMARY KEY CHECK (id = 1),
	token TEXT NOT NULL
);

CREATE TABLE servers_crafty_rcon (
	orchestrator_ref TEXT PRIMARY KEY,
	rcon_port INTEGER NOT NULL CHECK (rcon_port BETWEEN 1 AND 65535),
	rcon_password TEXT NOT NULL
);
