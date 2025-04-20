class MinecraftServerInfo:
    def __init__(self, minecraft_version: str, java_version: str, mod_engine: str, max_players: int, seed: str):
        self.minecraft_version = minecraft_version
        self.java_version = java_version
        self.mod_engine = mod_engine
        self.max_players = max_players
        self.seed = seed

    @classmethod
    def from_json(cls, json_data: dict):
        """
        Create a MinecraftServerInfo instance from a JSON dictionary
        with camelCase keys (for external compatibility).
        """
        return cls(
            minecraft_version=json_data.get("minecraftVersion"),
            java_version=json_data.get("javaVersion"),
            mod_engine=json_data.get("modEngine"),
            max_players=json_data.get("maxPlayers"),
            seed=json_data.get("seed")
        )

    def to_json(self) -> dict:
        """
        Convert the MinecraftServerInfo instance to a JSON dictionary
        with camelCase keys (for external compatibility).
        """
        return {
            "minecraftVersion": self.minecraft_version,
            "javaVersion": self.java_version,
            "modEngine": self.mod_engine,
            "maxPlayers": self.max_players,
            "seed": self.seed
        }