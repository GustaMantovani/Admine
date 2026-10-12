import { Rcon } from "rcon-client";
import type { RconConfig } from "./rcon-config";

export class RconClient {
	constructor(private readonly config: RconConfig) {}

	
	async getConnection(): Promise<Rcon> {
		return Rcon.connect({
			host: "localhost",
			port: this.config.port,
			password: this.config.password,
		});
	}

	async sendCommand(cmd: string): Promise<string> {
		const conn = await this.getConnection();
		const res = await conn.send(cmd);
		await conn.end();
		return res;
	}
}
