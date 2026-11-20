import type { Database } from "bun:sqlite";

export type ModuleMigrations = {
	module: string;
	migrations: readonly string[];
};

export function runMigrations(
	db: Database,
	modules: readonly ModuleMigrations[],
): void {
	db.run(
		"CREATE TABLE IF NOT EXISTS schema_versions (module TEXT PRIMARY KEY, version INTEGER NOT NULL)",
	);

	for (const { module, migrations } of modules) {
		const migrate = db.transaction(() => {
			const row = db
				.query("SELECT version FROM schema_versions WHERE module = ?")
				.get(module) as { version: number } | null;
			const current = row?.version ?? 0;
			if (current >= migrations.length) return;

			for (const sql of migrations.slice(current)) db.run(sql);

			db.query(
				"INSERT INTO schema_versions (module, version) VALUES (?, ?) ON CONFLICT (module) DO UPDATE SET version = excluded.version",
			).run(module, migrations.length);
		});
		migrate.immediate();
	}
}
