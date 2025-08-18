class ConfigError(Exception):
    def __init__(self, message: str = "Configuration error"):
        self.message = message
        super().__init__(self.message)


class ConfigFileError(ConfigError):
    def __init__(self, file_path: str, message: str = "Could not load configuration file"):
        self.file_path = file_path
        super().__init__(f"{message}: {file_path}")


class MessageServiceFactoryError(Exception):
    def __init__(self, provider_type, message="Failed to create message service provider"):
        self.provider_type = provider_type
        super().__init__(f"{message}: {provider_type}")


class MessageServiceError(Exception):
    def __init__(self, message: str = "Message service error"):
        self.message = message
        super().__init__(self.message)


class DiscordTokenException(MessageServiceError):
    def __init__(self, message: str = "Invalid Discord token"):
        self.message = message
        super().__init__(self.message)


class DiscordCommandPrefixException(MessageServiceError):
    def __init__(self, message: str = "Invalid Discord command prefix"):
        self.message = message
        super().__init__(self.message)


class PubSubServiceFactoryException(Exception):
    def __init__(self, provider_type, message="Failed to create PubSub service provider"):
        self.provider_type = provider_type
        super().__init__(f"{message}: {provider_type}")


class MinecraftInfoServiceFactoryException(Exception):
    def __init__(self, provider_type, message="Failed to create Minecraft info service provider"):
        self.provider_type = provider_type
        super().__init__(f"{message}: {provider_type}")


class VpnServiceFactoryException(Exception):
    def __init__(self, provider_type, message="Failed to create Vpn info service provider"):
        self.provider_type = provider_type
        super().__init__(f"{message}: {provider_type}")
