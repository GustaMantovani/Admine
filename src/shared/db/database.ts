import { Database } from "bun:sqlite";

export function openDatabase(path: string): Database {
	const db = new Database(path, { create: true, strict: true });
	db.run("PRAGMA journal_mode = WAL");
	db.run("PRAGMA foreign_keys = ON");
	db.run("PRAGMA busy_timeout = 5000");
	return db;
}
