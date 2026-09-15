export type CreateMinecraftServerRequest = {
    name: string;
    readonly type: "java" | "bedrock" | "hytale";
    version: string;
    minimumMemory?: number;
    maximumMemory?: number;
    port?: number;
};
