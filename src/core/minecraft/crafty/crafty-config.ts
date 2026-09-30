type BaseConfig = {
	apiBaseUrl: string;
	isSelfHosted: boolean;
};
export type CraftyConfig =
	| (BaseConfig & {
			// Self hosted options
			apiToken: string;
	  })
	| {
			// Managed by Admin options
			userName: string;
			password: string;
			defaultCredsFilePath: string;
			infraSetupCommand: string;
			infraCheckDelay: number;
			infraCheckCommand: string;
	  };
