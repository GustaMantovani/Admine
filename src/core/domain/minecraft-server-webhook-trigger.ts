export enum MinecraftServerWebhookTrigger {
	Started = "started",
	Stopped = "stopped",
	Crashed = "crashed",
	BackupCompleted = "backupCompleted",
	ExecutableUpdated = "executableUpdated",
	CommandReceived = "commandReceived",
	Killed = "killed",
}

// Where to send the webhook of each trigger. Triggers left out send nothing
export type MinecraftServerWebhookUrls = Partial<
	Record<MinecraftServerWebhookTrigger, string>
>;
