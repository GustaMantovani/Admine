export type CraftyCreateJavaServer = {
	name: string;
	monitoring_type: "minecraft_java";
	minecraft_java_monitoring_data: {
		host: string;
		port: number;
	};
	create_type: "minecraft_java";
	minecraft_java_create_data: {
		create_type: "download_jar";
		download_jar_create_data: {
			category: string;
			type: string;
			version: string;
			mem_min: number;
			mem_max: number;
			server_properties_port: number;
		};
	};
};
