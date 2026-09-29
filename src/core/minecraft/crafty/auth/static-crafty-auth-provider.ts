import type { CraftyConfig } from "../crafty-config";
import type { CraftyAuthProvider } from "./crafty-auth-provider";

export class StaticCraftyAuthProvider implements CraftyAuthProvider {
	constructor(private readonly config: CraftyConfig) {}

	async getToken(): Promise<string> {
		return this.config.apiToken;
	}
}
