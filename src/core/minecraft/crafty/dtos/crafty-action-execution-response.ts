export type CraftyActionExecutionResponse = {
	status: "ok" | "error";
	error?: string;
	error_data?: string;
};
