export type CraftyFilePatchRequest = {
	path: string;
	contents: string;
	modified_epoch?: number;
	overwrite?: boolean;
};
