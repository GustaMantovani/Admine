export type GameServerFileDownload = {
	contentStream: ReadableStream<Uint8Array>;
	fileNameOnServer?: string;
};
