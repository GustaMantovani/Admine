export type CraftyFileOperationRequest = {
	file_system_objects: {
		source_path: string;
		// A directory: the item keeps its name inside it
		target_path: string;
	}[];
};
