export type CraftyWebhookTrigger =
	| "start_server"
	| "stop_server"
	| "crash_detected"
	| "backup_server"
	| "jar_update"
	| "send_command"
	| "kill";
