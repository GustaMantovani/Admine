export class RconError extends Error {
	constructor(message: string, options?: ErrorOptions) {
		super(message, options);
		this.name = new.target.name;
	}
}

export class RconConnectionError extends RconError {}
export class RconAuthError extends RconError {}
export class RconTimeoutError extends RconError {}
