-- Single row: Crafty issues one token per Admine user
CREATE TABLE IF NOT EXISTS CraftyAuthToken (
	id INTEGER PRIMARY KEY CHECK (id = 1),
	token VARCHAR(255) NOT NULL
);

-- Crafty doesn't expose the RCON port, it only lives in server.properties
CREATE TABLE IF NOT EXISTS CraftyRconPort (
	serverId VARCHAR(36) PRIMARY KEY,
	rconPort INTEGER NOT NULL CHECK (rconPort BETWEEN 1 AND 65535)
);
