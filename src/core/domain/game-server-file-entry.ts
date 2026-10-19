export type GameServerFileEntry = {
	name: string;
	// Relative to the server's root directory
	path: string;
	type: "file" | "directory";
};
