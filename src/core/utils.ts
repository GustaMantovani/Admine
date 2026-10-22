import type { BunFile } from "bun";

export function randomOutside(
	min: number,
	max: number,
	excluded: Iterable<number>,
): number | undefined {
	const sorted = [...new Set(excluded)]
		.filter((n) => n >= min && n <= max)
		.sort((a, b) => a - b);

	const available = max - min + 1 - sorted.length;
	if (available <= 0) return undefined;

	let value = min + Math.floor(Math.random() * available);
	for (const n of sorted) {
		if (value >= n) value++;
	}
	return value;
}

export function splitFile(chunkSize: number, file: BunFile): Blob[] {
	const chunks: Blob[] = [];

	for (let offset = 0; offset < file.size; offset += chunkSize) {
		chunks.push(file.slice(offset, offset + chunkSize));
	}

	return chunks;
}

export async function hashChunk(chunk: Blob): Promise<string> {
	const buffer = await chunk.arrayBuffer();
	return new Bun.CryptoHasher("sha256").update(buffer).digest("hex");
}
