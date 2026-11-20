import { Rcon } from "rcon-client";
import {
	RconAuthError,
	RconConnectionError,
	type RconError,
	RconTimeoutError,
} from "./errors";

export type RconTarget = {
	host: string;
	port: number;
	password: string;
};

export type RconSendOptions = {
	timeoutMs?: number;
};

export interface RconClient {
	send(
		target: RconTarget,
		command: string,
		options?: RconSendOptions,
	): Promise<string>;
}

const DEFAULT_TIMEOUT_MS = 5000;

export class LibRconClient implements RconClient {
	private readonly queues = new Map<string, Promise<unknown>>();

	send(
		target: RconTarget,
		command: string,
		options?: RconSendOptions,
	): Promise<string> {
		const key = `${target.host}:${target.port}`;
		const previous = this.queues.get(key) ?? Promise.resolve();
		const result = previous.then(() =>
			this.sendOnce(target, command, options?.timeoutMs ?? DEFAULT_TIMEOUT_MS),
		);
		const tail = result.catch(() => undefined);
		this.queues.set(key, tail);
		tail.then(() => {
			if (this.queues.get(key) === tail) this.queues.delete(key);
		});
		return result;
	}

	private async sendOnce(
		target: RconTarget,
		command: string,
		timeoutMs: number,
	): Promise<string> {
		let connection: Rcon;
		try {
			connection = await Rcon.connect({
				host: target.host,
				port: target.port,
				password: target.password,
				timeout: timeoutMs,
			});
		} catch (error) {
			throw toRconError(error, target);
		}

		try {
			return await connection.send(command);
		} catch (error) {
			throw toRconError(error, target);
		} finally {
			await connection.end().catch(() => undefined);
		}
	}
}

function toRconError(error: unknown, target: RconTarget): RconError {
	const where = `${target.host}:${target.port}`;
	const message = error instanceof Error ? error.message : String(error);

	if (message === "Authentication failed")
		return new RconAuthError(`RCON authentication failed on ${where}`, {
			cause: error,
		});
	if (message.startsWith("Timeout for packet"))
		return new RconTimeoutError(`RCON timed out on ${where}`, { cause: error });
	return new RconConnectionError(
		`RCON connection failed on ${where}: ${message}`,
		{
			cause: error,
		},
	);
}
