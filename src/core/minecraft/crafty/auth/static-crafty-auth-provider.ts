import type { SelfHostedCraftyConfig } from "../crafty-config";
import type { CraftyAuthProvider } from "./crafty-auth-provider";

export class StaticCraftyAuthProvider implements CraftyAuthProvider {
	constructor(private readonly config: SelfHostedCraftyConfig) {}

	async getToken(): Promise<string> {
		return this.config.apiToken;
	}

	async getNewToken(): Promise<string> {
		return this.config.apiToken;
	}
}
