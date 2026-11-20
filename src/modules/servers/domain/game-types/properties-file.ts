import { PropertiesEditor } from "properties-file/editor";
import type { GameProperties } from "./game-type";

export function mergeIntoPropertiesFile(
	existingContent: string | undefined,
	config: GameProperties,
): string {
	const editor = new PropertiesEditor(existingContent ?? "");
	for (const [key, value] of Object.entries(config))
		editor.insert(key, String(value));

	return editor.format({
		removeComments: true,
		removeBlankLines: true,
		removeLeadingWhitespace: true,
		deduplicateKeys: true,
		separatorChar: "=",
		separatorLeading: "",
		separatorTrailing: "",
	});
}
