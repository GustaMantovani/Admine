import type { MinecraftServerWebhookUrls } from "../../types/minecraft-server-webhook-trigger";

export type CraftyInfraConfig = {
	infraSetupCommand: string;
	infraCheckDelay: number;
	infraCheckCommand: string;
};

export type CraftyWebhooksConfig = {
	webhooks: MinecraftServerWebhookUrls;
};

export type SelfHostedCraftyConfig = CraftyWebhooksConfig & {
	deployType: "self-hosted";
	apiBaseUrl: string;
	apiToken: string;
};

export type ManagedCraftyConfig = CraftyInfraConfig &
	CraftyWebhooksConfig & {
		deployType: "managed";
		apiBaseUrl: string;
		username: string;
		password: string;
	};

export type CraftyConfig = SelfHostedCraftyConfig | ManagedCraftyConfig;
