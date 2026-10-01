export type CraftyInfraConfig = {
	infraSetupCommand: string;
	infraCheckDelay: number;
	infraCheckCommand: string;
};

export type SelfHostedCraftyConfig = {
	deployType: "self-hosted";
	apiBaseUrl: string;
	apiToken: string;
};

export type ManagedCraftyConfig = CraftyInfraConfig & {
	deployType: "managed";
	apiBaseUrl: string;
	username: string;
	password: string;
};

export type CraftyConfig = SelfHostedCraftyConfig | ManagedCraftyConfig;
