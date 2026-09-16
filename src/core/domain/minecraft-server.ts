export class MinecraftServer {
    constructor(
        public readonly id: string,
        public name: string,
        public readonly type: "java" | "bedrock" | "hytale",
        public version?: string,
        public admineIpAddress?: string,
        public networkIpAddresses?: string[],
        public port?: number,
        public minimumMemory?: number,
        public maximumMemory?: number
    ) { }
}