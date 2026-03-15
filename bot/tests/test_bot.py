import json
import tempfile
from unittest.mock import MagicMock, patch

import pytest

from bot.bot import Bot
from bot.services.messaging.discord_message_service import MessageServiceProviderType
from bot.services.minecraft.server_handler_api_service import MinecraftServiceProviderType
from bot.services.pubsub.redis_pubsub_service import PubSubServiceProviderType
from bot.services.vpn.api_vpn_service import VpnServiceProviderType


@pytest.fixture
def sample_bot_config():
    """Creates a sample configuration for the bot."""
    return {
        "providers": {"messaging": "DISCORD", "pubsub": "REDIS", "minecraft": "REST", "vpn": "REST"},
        "logging": {"level": "DEBUG"},
        "discord": {"token": "test_token", "commandprefix": "!test", "administrators": ["123"], "channel_ids": []},
        "redis": {"connectionstring": "localhost:6379"},
        "minecraft": {"connectionstring": "http://localhost:3000/", "token": ""},
        "vpn": {"connectionstring": "http://localhost:9000", "token": ""},
    }


@pytest.fixture
def temp_bot_config_file(sample_bot_config):
    """Creates a temporary configuration file for the bot."""
    with tempfile.NamedTemporaryFile(mode="w", delete=False, suffix=".json") as f:
        json.dump(sample_bot_config, f)
        temp_file = f.name
    yield temp_file
    import os

    if os.path.exists(temp_file):
        os.remove(temp_file)


@pytest.fixture
def mock_config():
    """Creates a mock Config for the bot."""
    config = MagicMock()
    config.get.side_effect = lambda key, default=None: {
        "providers.pubsub": "REDIS",
        "providers.minecraft": "REST",
        "providers.vpn": "REST",
        "providers.messaging": "DISCORD",
    }.get(key, default)
    return config


@pytest.fixture
def mock_factories():
    """Creates mocks for all factories."""
    mock_pubsub = MagicMock()
    mock_pubsub.send_message = MagicMock()
    mock_pubsub.listen_message = MagicMock()

    mock_minecraft = MagicMock()
    mock_vpn = MagicMock()

    mock_message_service = MagicMock()
    mock_message_service.set_callback = MagicMock()
    mock_message_service.connect = MagicMock()

    return {"pubsub": mock_pubsub, "minecraft": mock_minecraft, "vpn": mock_vpn, "message": mock_message_service}


class TestBotInitialization:
    """Bot initialization tests."""

    @patch("bot.bot.PubSubServiceFactory.create")
    @patch("bot.bot.MinecraftServiceFactory.create")
    @patch("bot.bot.VpnServiceFactory.create")
    @patch("bot.bot.MessageServiceFactory.create")
    def test_bot_initialization_creates_all_services(
        self, mock_msg_factory, mock_vpn_factory, mock_mc_factory, mock_pubsub_factory, mock_config
    ):
        # Tests that initialization creates all necessary instances.
        mock_pubsub = MagicMock()
        mock_pubsub_factory.return_value = mock_pubsub

        mock_minecraft = MagicMock()
        mock_mc_factory.return_value = mock_minecraft

        mock_vpn = MagicMock()
        mock_vpn_factory.return_value = mock_vpn

        mock_message = MagicMock()
        mock_msg_factory.return_value = mock_message

        # Create the bot
        bot = Bot(mock_config)

        # Verify all factories were called with expected provider types
        mock_pubsub_factory.assert_called_once_with(PubSubServiceProviderType.REDIS, mock_config)
        mock_mc_factory.assert_called_once_with(MinecraftServiceProviderType.REST, mock_config)
        mock_vpn_factory.assert_called_once_with(VpnServiceProviderType.REST, mock_config)
        mock_msg_factory.assert_called_once_with(MessageServiceProviderType.DISCORD, mock_config)

        # Verify bot has the correct instances
        assert bot._Bot__pubsub_service == mock_pubsub
        assert bot._Bot__minecraft_info_service == mock_minecraft
        assert bot._Bot__vpn_service == mock_vpn
        assert mock_message in bot._Bot__message_services

    @patch("bot.bot.PubSubServiceFactory.create")
    @patch("bot.bot.MinecraftServiceFactory.create")
    @patch("bot.bot.VpnServiceFactory.create")
    @patch("bot.bot.MessageServiceFactory.create")
    def test_bot_initialization_creates_handles(
        self, mock_msg_factory, mock_vpn_factory, mock_mc_factory, mock_pubsub_factory, mock_config
    ):
        """Tests that handles are created during initialization."""
        mock_pubsub_factory.return_value = MagicMock()
        mock_mc_factory.return_value = MagicMock()
        mock_vpn_factory.return_value = MagicMock()
        mock_msg_factory.return_value = MagicMock()

        # Create the bot
        bot = Bot(mock_config)

        # Verify that handles were created
        assert bot._Bot__command_handle is not None
        assert bot._Bot__event_handle is not None


class TestBotProviderConfiguration:
    """Provider configuration tests."""

    @patch("bot.bot.PubSubServiceFactory.create")
    @patch("bot.bot.MinecraftServiceFactory.create")
    @patch("bot.bot.VpnServiceFactory.create")
    @patch("bot.bot.MessageServiceFactory.create")
    def test_bot_uses_config_provider_types(
        self, mock_msg_factory, mock_vpn_factory, mock_mc_factory, mock_pubsub_factory, mock_config
    ):
        # Tests that bot uses provider types from the configuration.
        mock_pubsub_factory.return_value = MagicMock()
        mock_mc_factory.return_value = MagicMock()
        mock_vpn_factory.return_value = MagicMock()
        mock_msg_factory.return_value = MagicMock()

        # Create the bot
        Bot(mock_config)

        # Verify correct provider types were used
        mock_pubsub_factory.assert_called_once()
        call_args = mock_pubsub_factory.call_args[0]
        assert call_args[0] == PubSubServiceProviderType.REDIS

        mock_mc_factory.assert_called_once()
        call_args = mock_mc_factory.call_args[0]
        assert call_args[0] == MinecraftServiceProviderType.REST

    @patch("bot.bot.PubSubServiceFactory.create")
    @patch("bot.bot.MinecraftServiceFactory.create")
    @patch("bot.bot.VpnServiceFactory.create")
    @patch("bot.bot.MessageServiceFactory.create")
    def test_bot_uses_default_provider_types(
        self, mock_msg_factory, mock_vpn_factory, mock_mc_factory, mock_pubsub_factory, mock_config
    ):
        # Tests that the bot uses default provider types when config doesn't specify.
        mock_pubsub_factory.return_value = MagicMock()
        mock_mc_factory.return_value = MagicMock()
        mock_vpn_factory.return_value = MagicMock()
        mock_msg_factory.return_value = MagicMock()

        # Create the bot
        Bot(mock_config)

        # Verify defaults were used
        mock_pubsub_factory.assert_called_once()
        call_args = mock_pubsub_factory.call_args[0]
        assert call_args[0] == PubSubServiceProviderType.REDIS

        mock_mc_factory.assert_called_once()
        call_args = mock_mc_factory.call_args[0]
        assert call_args[0] == MinecraftServiceProviderType.REST


class TestBotIntegration:
    """Integration tests for bot components."""

    @patch("bot.bot.PubSubServiceFactory.create")
    @patch("bot.bot.MinecraftServiceFactory.create")
    @patch("bot.bot.VpnServiceFactory.create")
    @patch("bot.bot.MessageServiceFactory.create")
    def test_command_handle_receives_correct_services(
        self, mock_msg_factory, mock_vpn_factory, mock_mc_factory, mock_pubsub_factory, mock_config
    ):
        # Tests that CommandHandle receives correct services.
        mock_pubsub = MagicMock()
        mock_minecraft = MagicMock()
        mock_vpn = MagicMock()

        mock_pubsub_factory.return_value = mock_pubsub
        mock_mc_factory.return_value = mock_minecraft
        mock_vpn_factory.return_value = mock_vpn
        mock_msg_factory.return_value = MagicMock()

        # Create the bot
        bot = Bot(mock_config)

        # Verify CommandHandle has correct services
        command_handle = bot._Bot__command_handle
        assert command_handle._CommandHandle__pubsub_service == mock_pubsub
        assert command_handle._CommandHandle__minecraft_info_service == mock_minecraft
        assert command_handle._CommandHandle__vpn_service == mock_vpn

    @patch("bot.bot.PubSubServiceFactory.create")
    @patch("bot.bot.MinecraftServiceFactory.create")
    @patch("bot.bot.VpnServiceFactory.create")
    @patch("bot.bot.MessageServiceFactory.create")
    def test_event_handle_receives_message_services(
        self, mock_msg_factory, mock_vpn_factory, mock_mc_factory, mock_pubsub_factory, mock_config
    ):
        # Tests that EventHandle receives the correct message services.
        mock_message_service = MagicMock()

        mock_pubsub_factory.return_value = MagicMock()
        mock_mc_factory.return_value = MagicMock()
        mock_vpn_factory.return_value = MagicMock()
        mock_msg_factory.return_value = mock_message_service

        # Create the bot
        bot = Bot(mock_config)

        # Verify EventHandle has the correct message services
        event_handle = bot._Bot__event_handle
        assert mock_message_service in event_handle._EventHandle__message_services
