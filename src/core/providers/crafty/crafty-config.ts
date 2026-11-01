import type { GameServerWebhookUrls } from "../../domain/game-server-webhook-trigger";

export type CraftyWebhooksConfig = {
	enabled: boolean;
	urls: GameServerWebhookUrls;
};

export type SelfHostedCraftyDeployment = {
	type: "self-hosted";
	apiToken: string;
};

export type ManagedCraftyDeployment = {
	type: "managed";
	username: string;
	password: string;
};

export type CraftyDeployment =
	| SelfHostedCraftyDeployment
	| ManagedCraftyDeployment;

export type CraftyConfig = {
	apiBaseUrl: string;
	webhooks: CraftyWebhooksConfig;
	maxConcurrentRequests: number;
	deployment: CraftyDeployment;
};
