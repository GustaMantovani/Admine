export type GameServerFile = {
	path: string;
	content: string;
	// Sent back on write so the provider can reject it if the file changed since it was read
	modifiedEpoch?: number;
};
