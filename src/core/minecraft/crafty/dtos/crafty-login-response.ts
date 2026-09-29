export type CraftyLoginResponse = {
	status: "ok" | "error";
	data: {
		token: string;
	};
};
