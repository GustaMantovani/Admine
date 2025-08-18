class MinecraftServerInfo:
    def __init__(
        self,
        minecraft_version: str,
        java_version: str,
        mod_engine: str,
        max_players: int,
        seed: str,
    ):
        self.minecraft_version = minecraft_version
        self.java_version = java_version
        self.mod_engine = mod_engine
        self.max_players = max_players
        self.seed = seed

    @classmethod
    def from_json(cls, json_data: dict):
        return cls(
            minecraft_version=json_data.get("minecraftVersion"),
            java_version=json_data.get("javaVersion"),
            mod_engine=json_data.get("modEngine"),
            max_players=json_data.get("maxPlayers"),
            seed=json_data.get("seed"),
        )

    def to_json(self) -> dict:
        return {
            "minecraftVersion": self.minecraft_version,
            "javaVersion": self.java_version,
            "modEngine": self.mod_engine,
            "maxPlayers": self.max_players,
            "seed": self.seed,
        }

    def __str__(self):
        return (
            f"Minecraft Version: {self.minecraft_version}\n"
            f"Java Version: {self.java_version}\n"
            f"Mod Engine: {self.mod_engine}\n"
            f"Max Players: {self.max_players}\n"
            f"Seed: {self.seed}"
        )
