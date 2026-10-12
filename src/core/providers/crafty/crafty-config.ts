import type { GameServerWebhookUrls } from "../../domain/game-server-webhook-trigger";

export type CraftyWebhooksConfig = {
	enabled: boolean;
	urls: GameServerWebhookUrls;
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
