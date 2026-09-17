export type MinecraftServerStats = {
    running: boolean;
    crashed: boolean;
    startedAt?: string;
    cpu: number;
    memoryBytes: number;
    memoryPercent: number;
    playersOnline: number;
    maxPlayers: number;
    version?: string;
    worldName: string;
    worldSize: string;
};
