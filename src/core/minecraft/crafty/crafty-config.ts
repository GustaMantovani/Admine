export type CraftyConfig = {
	// General
	apiBaseUrl: string;
	isSelfHosted: boolean;

	// Self hosted options
	apiToken: string;

	// Managed by Admin options
	userName: string;
	password: string;
	defaultCredsFilePath: string;
	infraSetupCommand: string;
	infraCheckDelay: number;
	infraCheckCommand: string;
};
