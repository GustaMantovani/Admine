import { z } from "zod";

export const CRAFTY_WEBHOOK_TRIGGERS = [
	"started",
	"stopped",
	"crashed",
	"backupCompleted",
	"executableUpdated",
	"commandReceived",
	"killed",
] as const;
export type CraftyWebhookEvent = (typeof CRAFTY_WEBHOOK_TRIGGERS)[number];

const ManagedDeploymentSchema = z.object({
	type: z.literal("managed"),
	infraSetupCommand: z
		.string()
		.default(
			'docker start crafty_container 2>/dev/null || docker run --name crafty_container --detach --net=host -e TZ=Etc/UTC -v "/$(pwd)/.data/crafty/backups:/crafty/backups" -v "/$(pwd)/.data/crafty/logs:/crafty/logs" -v "/$(pwd)/.data/crafty/servers:/crafty/servers" -v "/$(pwd)/.data/crafty/config:/crafty/app/config" -v "/$(pwd)/.data/crafty/import:/crafty/import" arcadiatechnology/crafty-4:latest',
		),
	infraCheckCommand: z
		.string()
		.default(
			"docker inspect -f '{{.State.Running}}' crafty_container | grep -q true",
		),
	infraCheckDelay: z.number().default(30000),
	username: z.string().optional(),
	password: z.string().optional(),
});

const SelfHostedDeploymentSchema = z.object({
	type: z.literal("self-hosted"),
	apiToken: z.string(),
});

export const CraftyConfigSchema = z.object({
	kind: z.literal("crafty"),
	apiBaseUrl: z.string().default("https://localhost:8443/api/v2"),
	serversHost: z.string().default("127.0.0.1"),
	maxConcurrentRequests: z.number().int().positive().default(4),
	webhooks: z
		.object({
			enabled: z.boolean().default(true),
			urls: z
				.partialRecord(z.enum(CRAFTY_WEBHOOK_TRIGGERS), z.url())
				.default({}),
		})
		.prefault({}),
	deployment: z
		.discriminatedUnion("type", [
			ManagedDeploymentSchema,
			SelfHostedDeploymentSchema,
		])
		.prefault({ type: "managed" }),
});

export type CraftyConfig = z.infer<typeof CraftyConfigSchema>;
export type ManagedCraftyDeployment = z.infer<typeof ManagedDeploymentSchema>;
export type SelfHostedCraftyDeployment = z.infer<
	typeof SelfHostedDeploymentSchema
>;
