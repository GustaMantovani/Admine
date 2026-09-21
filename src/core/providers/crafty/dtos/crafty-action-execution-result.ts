export type CraftyActionExecutionResult = {
	status: "ok" | "error";
	error?: string;
	error_data?: string;
};
