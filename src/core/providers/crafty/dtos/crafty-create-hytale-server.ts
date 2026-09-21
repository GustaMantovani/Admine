export type CraftyCreateHytaleServer = {
	name: string;
	monitoring_type: "hytale";
	hytale_monitoring_data: {
		host: string;
		port: number;
	};
	create_type: "hytale";
	hytale_create_data: {
		create_type: "download_exe";
		download_exe_create_data: {
			agree_to_eula: true;
			mem_min: number;
			mem_max: number;
		};
	};
};
