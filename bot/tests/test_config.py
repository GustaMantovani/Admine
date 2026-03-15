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
    if os.path.exists(temp_file):
        os.remove(temp_file)


class TestConfigInstantiation:
    """Tests for Config instantiation behavior."""

    def test_config_creates_independent_instances(self, temp_config_file):
        """Verifies that each Config() call creates an independent instance."""
        config1 = Config(temp_config_file)
        config2 = Config(temp_config_file)

        assert config1 is not config2

    def test_config_loads_values_from_file(self, temp_config_file):
        """Verifies that config values are loaded from the given file."""
        config = Config(temp_config_file)

        assert config.get("logging.level") == "DEBUG"
        assert config.get("discord.token") == "test_token"


class TestConfigLoading:
    """Tests for configuration loading."""

    def test_load_from_json_success(self, temp_config_file):
        """Tests successful loading of a valid JSON file."""
        config = Config(temp_config_file)

        value = config.get("providers.messaging")
        assert value is not None

    def test_load_from_nonexistent_file(self):
        """Tests behavior when file does not exist."""
        with pytest.raises(ConfigError, match="Missing required"):
            Config("nonexistent_file.json")

    def test_load_invalid_json(self):
        """Tests error handling when JSON is invalid."""
        with tempfile.NamedTemporaryFile(mode="w", delete=False, suffix=".json") as f:
            f.write("{ invalid json content")
            temp_file = f.name

        try:
            with pytest.raises(ConfigFileError, match="Invalid JSON format"):
                Config(temp_file)
        finally:
            os.remove(temp_file)

    def test_load_with_io_error(self):
        """Tests file reading error handling."""
        with patch("builtins.open", side_effect=IOError("Permission denied")):
            with patch("os.path.exists", return_value=True):
                with pytest.raises(ConfigFileError, match="Error reading file"):
                    Config("some_file.json")


class TestConfigGet:
    """Tests for Config get method."""

    def test_get_existing_value(self, temp_config_file):
        """Tests retrieval of an existing value."""
        config = Config(temp_config_file)

        result = config.get("logging.level")
        assert result is not None

    def test_get_nested_value(self, temp_config_file):
        """Tests retrieval of nested values using dot notation."""
        config = Config(temp_config_file)

        result = config.get("providers.messaging")
        assert result is not None

    def test_get_nonexistent_key_with_default(self, temp_config_file):
        """Tests default value return when key does not exist."""
        config = Config(temp_config_file)

        assert config.get("nonexistent.key", "default_value") == "default_value"

    def test_get_nonexistent_key_without_default(self, temp_config_file):
        """Tests None return when key does not exist and there is no default value."""
        config = Config(temp_config_file)

        assert config.get("nonexistent.key") is None

    def test_get_with_invalid_path(self, temp_config_file):
        """Tests behavior when path traverses a non-dict value."""
        config = Config(temp_config_file)

        # "logging.level" is a string, so "logging.level.something" should return default
        assert config.get("logging.level.something", "default") == "default"

    def test_get_list_value(self, temp_config_file):
        """Tests retrieval of values that are lists."""
        config = Config(temp_config_file)

        administrators = config.get("discord.administrators")
        assert isinstance(administrators, list)

    def test_get_root_level_value(self, temp_config_file):
        """Tests retrieval of root level values."""
        config = Config(temp_config_file)

        providers = config.get("providers")
        assert isinstance(providers, dict)
