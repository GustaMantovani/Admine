import type { SelfHostedCraftyDeployment } from "../crafty-config";
import type { CraftyAuthProvider } from "./crafty-auth-provider";

export class StaticCraftyAuthProvider implements CraftyAuthProvider {
	constructor(private readonly deployment: SelfHostedCraftyDeployment) {}

	async getToken(): Promise<string> {
		return this.deployment.apiToken;
	}

	async getNewToken(): Promise<string> {
		return this.deployment.apiToken;
	}
}
