export class MinecraftServer {
	constructor(
		public readonly id: string,
		public readonly name: string,
		public readonly type: "java" | "bedrock" | "hytale",
		public readonly status: "on" | "off" | "unknown",
		public readonly version?: string,
		public admineIpAddress?: string,
		public networkIpAddresses?: string[],
		public readonly port?: number,
		public readonly minimumMemory?: number,
		public readonly maximumMemory?: number,
	) {}
}
