import type { CraftyActionExecutionResult } from "./crafty-action-execution-result";

export type CraftyResponse<T> = CraftyActionExecutionResult & {
	data: T;
};
