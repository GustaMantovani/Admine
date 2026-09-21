import type { CraftyServer } from "./crafty-server";

export type CraftyServerDetail = CraftyServer & {
	status: {
		update_available: boolean;
		updating: boolean;
		backing_up: boolean;
		last_backup: boolean;
	};
};
