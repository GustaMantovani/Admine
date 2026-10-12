import { Rcon } from "rcon-client";

export async function sendCommand(
	host: string = "localhost",
	port: number = 25575,
	password: string,
	cmd: string,
): Promise<string> {
	const conn = await Rcon.connect({
		host: host,
		port: port,
		password: password,
	});
	const res = await conn.send(cmd);
	await conn.end();
	return res;
}
