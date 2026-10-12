import type { CraftyWebhookTrigger } from "./crafty-webhook-trigger";

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
