export type CraftyFile = {
	status: string;
	data: {
		content: string;
		attributes: {
			mime: string;
			modified: string;
			size: string;
			modified_epoch: number;
		};
	};
};
