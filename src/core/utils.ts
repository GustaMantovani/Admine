export function randomOutside(
	min: number,
	max: number,
	excluded: Iterable<number>,
): number | null {
	const sorted = [...new Set(excluded)]
		.filter((n) => n >= min && n <= max)
		.sort((a, b) => a - b);

	const available = max - min + 1 - sorted.length;
	if (available <= 0) return null;

	let value = min + Math.floor(Math.random() * available);
	for (const n of sorted) {
		if (value >= n) value++;
	}
	return value;
}
