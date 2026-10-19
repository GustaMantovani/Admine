export type CraftyDirectoryEntry = {
	path: string;
	dir: boolean;
	excluded: boolean;
	permissions: {
		can_read: boolean;
		can_write: boolean;
		can_execute: boolean;
	};
	modified: string;
	// Only present for files
	can_open?: boolean;
	mime?: string;
	size?: string;
};

// Crafty puts the entries next to root_path, keyed by file name
export type CraftyDirectory = {
	root_path: {
		local_path: string;
		path: string;
		top: boolean;
		modified: number;
	};
} & {
	[name: string]: CraftyDirectoryEntry;
};
