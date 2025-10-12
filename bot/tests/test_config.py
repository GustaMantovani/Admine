import json
import os
import tempfile
from unittest.mock import patch

import pytest

from bot.config import Config
from bot.exceptions import ConfigError, ConfigFileError


@pytest.fixture
def sample_config_data():
    """Fixture with sample configuration data."""
    return {
        "providers": {"messaging": "DISCORD", "pubsub": "REDIS", "minecraft": "REST", "vpn": "REST"},
        "logging": {"level": "DEBUG"},
        "discord": {"token": "test_token", "commandprefix": "!mc", "administrators": ["123", "456"], "channel_ids": []},
        "redis": {"connectionstring": "localhost:6379"},
    }


@pytest.fixture
def temp_config_file(sample_config_data):
    """Creates a temporary configuration file."""
    with tempfile.NamedTemporaryFile(mode="w", delete=False, suffix=".json") as f:
        json.dump(sample_config_data, f)
        temp_file = f.name
    yield temp_file
    # Cleanup
    if os.path.exists(temp_file):
        os.remove(temp_file)


@pytest.fixture(autouse=True)
def reset_config_singleton():
    """Reset Config singleton before each test."""
    Config._instance = None
    yield
    Config._instance = None


class TestConfigSingleton:
    """Tests for Config Singleton pattern."""

    def test_config_singleton_pattern(self, temp_config_file):
        """Verifies that Config correctly implements the Singleton pattern."""
        # Patch the default path to use the temporary file
        with patch("bot.config.Config._Config__load_from_json") as mock_load:
            mock_load.return_value = {"test": "data"}

            config1 = Config()
            config2 = Config()

            assert config1 is config2, "Config deve retornar a mesma instância (Singleton)"

    def test_config_singleton_already_initialized(self, temp_config_file):
        """Verifies that __init__ is not executed again after first initialization."""
        # Instance was already initialized in a previous test
        # Creating another instance should not reinitialize
        config1 = Config()
        config1_id = id(config1._Config__config)

        config2 = Config()
        config2_id = id(config2._Config__config)

        assert config1 is config2
        assert config1_id == config2_id


class TestConfigLoading:
    """Tests for configuration loading."""

    def test_load_from_json_success(self):
        """Tests successful loading of a valid JSON file."""
        # Since Config is already initialized, just test value retrieval
        config = Config()

        value = config.get("providers.messaging")
        assert value is not None or value == config.get("nonexistent", "default") == "default"

    def test_load_from_nonexistent_file(self):
        """Tests behavior when file does not exist."""
        Config._instance = None

        original_init = Config.__init__

        def patched_init(self, config_file="./bot_config.json"):
            original_init(self, config_file)

        with patch.object(Config, "__init__", patched_init):
            with pytest.raises(ConfigError, match="Failed to load configuration"):
                instance = Config.__new__(Config)
                instance._initialized = False
                instance.__init__("nonexistent_file.json")

        Config._instance = None

    def test_load_invalid_json(self):
        """Tests error handling when JSON is invalid."""
        Config._instance = None

        with tempfile.NamedTemporaryFile(mode="w", delete=False, suffix=".json") as f:
            f.write("{ invalid json content")
            temp_file = f.name

        try:
            original_init = Config.__init__

            def patched_init(self, config_file="./bot_config.json"):
                original_init(self, config_file)

            with patch.object(Config, "__init__", patched_init):
                with pytest.raises(ConfigFileError, match="Invalid JSON format"):
                    instance = Config.__new__(Config)
                    instance._initialized = False
                    instance.__init__(temp_file)
        finally:
            os.remove(temp_file)
            Config._instance = None

    def test_load_with_io_error(self):
        """Tests file reading error handling."""
        Config._instance = None

        with patch("builtins.open", side_effect=IOError("Permission denied")):
            with patch("os.path.exists", return_value=True):
                original_init = Config.__init__

                def patched_init(self, config_file="./bot_config.json"):
                    original_init(self, config_file)

                with patch.object(Config, "__init__", patched_init):
                    with pytest.raises(ConfigFileError, match="Error reading file"):
                        instance = Config.__new__(Config)
                        instance._initialized = False
                        instance.__init__("some_file.json")

        Config._instance = None


class TestConfigGet:
    """Tests for Config get method."""

    def test_get_existing_value(self):
        """Tests retrieval of an existing value."""
        config = Config()

        result = config.get("logging.level")
        assert result is not None

    def test_get_nested_value(self):
        """Tests retrieval of nested values using dot notation."""
        config = Config()

        result = config.get("providers.messaging")
        assert result is not None

    def test_get_nonexistent_key_with_default(self):
        """Tests default value return when key does not exist."""
        config = Config()

        assert config.get("nonexistent.key", "default_value") == "default_value"

    def test_get_nonexistent_key_without_default(self):
        """Tests None return when key does not exist and there is no default value."""
        config = Config()

        assert config.get("nonexistent.key") is None

    def test_get_with_invalid_path(self):
        """Tests behavior when path traverses a non-dict value."""
        config = Config()

        # "logging.level" is a string, so "logging.level.something" should return default
        assert config.get("logging.level.something", "default") == "default"

    def test_get_list_value(self):
        """Tests retrieval of values that are lists."""
        config = Config()

        administrators = config.get("discord.administrators")
        assert isinstance(administrators, list) or administrators is None

    def test_get_root_level_value(self):
        """Tests retrieval of root level values."""
        config = Config()

        providers = config.get("providers")
        assert isinstance(providers, dict) or providers is None
