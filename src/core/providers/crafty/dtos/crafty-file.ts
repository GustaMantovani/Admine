export type CraftyFile = {
	content: string;
	attributes: {
		mime: string;
		modified: string;
		size: string;
		modified_epoch: number;
	};
};
