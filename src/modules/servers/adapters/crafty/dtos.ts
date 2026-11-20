export type CraftyActionExecutionResponse = {
	status: "ok" | "error";
	error?: string;
	error_data?: string;
};

export type CraftyResponse<T> = CraftyActionExecutionResponse & {
	data: T;
};

export type CraftyAction =
	| "clone_server"
	| "start_server"
	| "stop_server"
	| "restart_server"
	| "kill_server"
	| "backup_server"
	| "update_executable"
	| "eula";

export type CraftyCreateJavaServer = {
	name: string;
	monitoring_type: "minecraft_java";
	minecraft_java_monitoring_data: {
		host: string;
		port: number;
	};
	create_type: "minecraft_java";
	minecraft_java_create_data: {
		create_type: "download_jar";
		download_jar_create_data: {
			category: string;
			type: string;
			version: string;
			mem_min: number;
			mem_max: number;
			server_properties_port: number;
		};
	};
};

export type CraftyCreateBedrockServer = {
	name: string;
	monitoring_type: "minecraft_bedrock";
	minecraft_bedrock_monitoring_data: {
		host: string;
		port: number;
	};
	create_type: "minecraft_bedrock";
	minecraft_bedrock_create_data: {
		create_type: "download_exe";
		download_exe_create_data: {
			agree_to_eula: true;
		};
	};
};

export type CraftyCreateHytaleServer = {
	name: string;
	monitoring_type: "hytale";
	hytale_monitoring_data: {
		host: string;
		port: number;
	};
	create_type: "hytale";
	hytale_create_data: {
		create_type: "download_exe";
		download_exe_create_data: {
			agree_to_eula: true;
			mem_min: number;
			mem_max: number;
		};
	};
};

export type CraftyCreateServer =
	| CraftyCreateJavaServer
	| CraftyCreateBedrockServer
	| CraftyCreateHytaleServer;

export type CraftyCreateServerResponse = {
	new_server_id: string;
	new_server_uuid: string;
};

export type CraftyWebhookTrigger =
	| "start_server"
	| "stop_server"
	| "crash_detected"
	| "backup_server"
	| "jar_update"
	| "send_command"
	| "kill";

export type CraftyCreateWebhook = {
	webhook_type: "Discord" | "Mattermost" | "Slack" | "Teams";
	name: string;
	url: string;
	bot_name: string;
	trigger: CraftyWebhookTrigger[];
	body: string;
	color: string;
	enabled: boolean;
};

export type CraftyCreateWebhookResponse = {
	webhook_id: number;
};

export type CraftyCredsFile = {
	username: string;
	password: string;
	info?: string;
};

export type CraftyLoginResponse = {
	status: "ok" | "error";
	data: {
		token: string;
	};
};

export type CraftyServer = {
	server_id: string;
	created: string;
	server_uuid: string;
	server_name: string;
	path: string;
	backup_path: string;
	executable: string;
	log_path: string;
	execution_command: string;
	auto_start: boolean;
	auto_start_delay: number;
	crash_detection: boolean;
	stop_command: string;
	executable_update_url: string;
	server_ip: string;
	server_port: number;
	logs_delete_after: number;
	type: string;
};

export type CraftyServerDetail = CraftyServer & {
	status: {
		update_available: boolean;
		updating: boolean;
		backing_up: boolean;
		last_backup: boolean;
	};
};

export type CraftyServerStats = {
	created: string;
	started: string;
	running: boolean;
	cpu: number;
	mem: number | string;
	mem_percent: number;
	world_name: string;
	world_size: string;
	server_port: number;
	online: number;
	max: number;
	players: string;
	desc: string;
	version: string;
	updating: boolean;
	waiting_start: boolean;
	crashed: boolean;
	importing: boolean;
};

export type CraftyStopAllResponse = {
	message: string;
	count: number;
};

export type CraftyFile = {
	content: string;
	attributes: {
		mime: string;
		modified: string;
		size: string;
		modified_epoch: number;
	};
};

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
	can_open?: boolean;
	mime?: string;
	size?: string;
};

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

export type CraftyFileCreateRequest = {
	parent: string;
	name: string;
	directory: boolean;
};

export type CraftyFileDeleteRequest = {
	file_system_objects: {
		filename: string;
	}[];
};

export type CraftyFileOperation = "move" | "copy";

export type CraftyFileOperationRequest = {
	file_system_objects: {
		source_path: string;
		target_path: string;
	}[];
};

export type CraftyFilePatchRequest = {
	path: string;
	contents: string;
	modified_epoch?: number;
	overwrite?: boolean;
};

export type CraftyFileRenameRequest = {
	path: string;
	new_name: string;
};
