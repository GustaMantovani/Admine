import type { MinecraftServerWebhookUrls } from "../../types/minecraft-server-webhook-trigger";

export type CraftyWebhooksConfig = {
	enabled: boolean;
	urls: MinecraftServerWebhookUrls;
};

export type SelfHostedCraftyConfig = {
	deployType: "self-hosted";
	apiBaseUrl: string;
	apiToken: string;
	webhooks: CraftyWebhooksConfig;
};

export type ManagedCraftyConfig = {
	deployType: "managed";
	apiBaseUrl: string;
	username: string;
	password: string;
	webhooks: CraftyWebhooksConfig;
};

export type CraftyConfig = SelfHostedCraftyConfig | ManagedCraftyConfig;
