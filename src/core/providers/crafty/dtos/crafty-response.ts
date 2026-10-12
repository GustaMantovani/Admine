import type { CraftyActionExecutionResponse } from "./crafty-action-execution-response";

export type CraftyResponse<T> = CraftyActionExecutionResponse & {
	data: T;
};
