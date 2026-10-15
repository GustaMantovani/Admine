import { Database } from "bun:sqlite";
import schema from "./schema.sql" with { type: "text" };

export function getDatabase(path: string): Database {
	const db = new Database(path);
	db.run("PRAGMA foreign_keys = ON");
	db.run(schema);
	return db;
}
