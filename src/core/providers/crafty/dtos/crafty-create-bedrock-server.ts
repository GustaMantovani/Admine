export type CraftyCreateBedrockServer = {
    name: string;
    monitoring_type: "minecraft_bedrock";
    minecraft_bedrock_monitoring_data: {
        host: string;
        port: number;
    };
    create_type: "minecraft_bedrock";
    minecraft_bedrock_create_data: {
        create_type: "download_exe";
        download_exe_create_data: {
            agree_to_eula: true;
        };
    };
};
