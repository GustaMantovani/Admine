from unittest.mock import AsyncMock, MagicMock

import pytest

from bot.handles.command_handle import CommandHandle
from bot.models.admine_message import AdmineMessage


@pytest.fixture
def mock_services():
    """Creates mocks for the services used by CommandHandle."""
    pubsub_service = MagicMock()
    pubsub_service.send_message = MagicMock()

    minecraft_service = MagicMock()
    minecraft_service.command = AsyncMock(return_value={"exit_code": 0, "output": "Command executed"})
    minecraft_service.get_info = AsyncMock(return_value={"version": "1.20", "players": 5})
    minecraft_service.get_status = AsyncMock(return_value={"status": "running", "health": "healthy"})

    vpn_service = MagicMock()
    vpn_service.auth_member = AsyncMock(return_value="Member authorized")
    vpn_service.get_vpn_id = AsyncMock(return_value="vpn-123456")
    vpn_service.get_server_ips = AsyncMock(return_value=["192.168.1.1", "192.168.1.2"])

    return {"pubsub": pubsub_service, "minecraft": minecraft_service, "vpn": vpn_service}


@pytest.fixture
def command_handle(mock_services):
    """Creates an instance of CommandHandle with mocked services."""
    return CommandHandle(mock_services["pubsub"], mock_services["minecraft"], mock_services["vpn"])


class TestCommandHandleBasics:
    """Basic tests for command processing."""

    @pytest.mark.asyncio
    async def test_unknown_command(self, command_handle):
        """Verifies that unknown commands are handled correctly."""
        result = await command_handle.process_command("unknown_command")

        assert result == "Unknown command"

    @pytest.mark.asyncio
    async def test_command_with_none_args(self, command_handle):
        """Verifies that args=None is handled correctly."""
        result = await command_handle.process_command("info", None)

        # Should not cause an error, should process normally
        assert result is not None


class TestAdminCommands:
    """Tests for administrative commands."""

    @pytest.mark.asyncio
    async def test_admin_command_with_authorized_user(self, command_handle, mock_services):
        """Verifies that an authorized user can execute administrative commands."""
        await command_handle.process_command("on", [], user_id="admin_user", administrators=["admin_user"])

        # Verifies that pubsub was called
        mock_services["pubsub"].send_message.assert_called_once()

        # Verifies that the message has the correct structure
        call_args = mock_services["pubsub"].send_message.call_args[0][0]
        assert isinstance(call_args, AdmineMessage)
        assert call_args.origin == "Bot"
        assert "server_on" in call_args.tags

    @pytest.mark.asyncio
    async def test_admin_command_with_unauthorized_user(self, command_handle, mock_services):
        """Verifies that an unauthorized user is blocked from using admin commands."""
        result = await command_handle.process_command("on", [], user_id="regular_user", administrators=["admin_user"])

        assert result == "Unauthorized command usage"
        # Verifies that the service was not called
        mock_services["pubsub"].send_message.assert_not_called()

    @pytest.mark.asyncio
    async def test_admin_command_without_user_id(self, command_handle, mock_services):
        """Verifies that an admin command without a user_id is blocked."""
        result = await command_handle.process_command("on", [], user_id=None, administrators=["admin_user"])

        assert result == "Unauthorized command usage"
        mock_services["pubsub"].send_message.assert_not_called()

    @pytest.mark.asyncio
    async def test_admin_command_without_administrators_list(self, command_handle, mock_services):
        """Verifies that an admin command without an administrators list is blocked."""
        result = await command_handle.process_command("on", [], user_id="user_123", administrators=None)

        assert result == "Unauthorized command usage"
        mock_services["pubsub"].send_message.assert_not_called()


class TestServerControlCommands:
    """Tests for server control commands (on/off/restart)."""

    @pytest.mark.asyncio
    async def test_server_on_command(self, command_handle, mock_services):
        """Tests the server on command."""
        await command_handle.process_command("on", [], user_id="admin", administrators=["admin"])

        mock_services["pubsub"].send_message.assert_called_once()
        message = mock_services["pubsub"].send_message.call_args[0][0]
        assert "server_on" in message.tags

    @pytest.mark.asyncio
    async def test_server_off_command(self, command_handle, mock_services):
        """Tests the server off command."""
        await command_handle.process_command("off", [], user_id="admin", administrators=["admin"])

        mock_services["pubsub"].send_message.assert_called_once()
        message = mock_services["pubsub"].send_message.call_args[0][0]
        assert "server_off" in message.tags

    @pytest.mark.asyncio
    async def test_restart_command(self, command_handle, mock_services):
        """Tests the server restart command."""
        await command_handle.process_command("restart", [], user_id="admin", administrators=["admin"])

        mock_services["pubsub"].send_message.assert_called_once()
        message = mock_services["pubsub"].send_message.call_args[0][0]
        assert "restart" in message.tags


class TestMinecraftCommands:
    """Tests for Minecraft related commands."""

    @pytest.mark.asyncio
    async def test_minecraft_command_execution(self, command_handle, mock_services):
        """Verifies that Minecraft commands are forwarded correctly."""
        result = await command_handle.process_command(
            "command", ["say", "Hello", "World"], user_id="admin", administrators=["admin"]
        )

        mock_services["minecraft"].command.assert_called_once_with("say Hello World")
        assert result == {"exit_code": 0, "output": "Command executed"}

    @pytest.mark.asyncio
    async def test_minecraft_command_with_error(self, command_handle, mock_services):
        """Verifies that Minecraft service errors are handled appropriately."""
        mock_services["minecraft"].command.side_effect = Exception("Connection error")

        result = await command_handle.process_command("command", ["test"], user_id="admin", administrators=["admin"])

        assert result == {"error": "Error executing command"}

    @pytest.mark.asyncio
    async def test_info_command(self, command_handle, mock_services):
        """Tests the server info command."""
        result = await command_handle.process_command("info", [])

        mock_services["minecraft"].get_info.assert_called_once()
        assert result == {"version": "1.20", "players": 5}

    @pytest.mark.asyncio
    async def test_info_command_with_error(self, command_handle, mock_services):
        """Tests error handling in info command."""
        mock_services["minecraft"].get_info.side_effect = Exception("Service unavailable")

        result = await command_handle.process_command("info", [])

        assert result == {"error": "Error getting server info"}

    @pytest.mark.asyncio
    async def test_status_command(self, command_handle, mock_services):
        """Tests the server status command."""
        result = await command_handle.process_command("status", [])

        mock_services["minecraft"].get_status.assert_called_once()
        assert result == {"status": "running", "health": "healthy"}

    @pytest.mark.asyncio
    async def test_status_command_with_error(self, command_handle, mock_services):
        """Tests error handling in status command."""
        mock_services["minecraft"].get_status.side_effect = Exception("Connection timeout")

        result = await command_handle.process_command("status", [])

        assert result == {"error": "Error getting server status"}


class TestVPNCommands:
    """Tests for VPN related commands."""

    @pytest.mark.asyncio
    async def test_auth_member_command(self, command_handle, mock_services):
        """Tests the member authorization command."""
        result = await command_handle.process_command("auth", ["member_123"])

        mock_services["vpn"].auth_member.assert_called_once_with("member_123")
        assert result == "Member authorized"

    @pytest.mark.asyncio
    async def test_auth_member_with_multiple_args(self, command_handle, mock_services):
        """Tests authorization with multiple arguments."""
        await command_handle.process_command("auth", ["member", "123", "test"])

        mock_services["vpn"].auth_member.assert_called_once_with("member 123 test")

    @pytest.mark.asyncio
    async def test_auth_member_with_error(self, command_handle, mock_services):
        """Tests error handling in authorization."""
        mock_services["vpn"].auth_member.side_effect = Exception("Auth service error")

        result = await command_handle.process_command("auth", ["member_123"])

        assert result == "Error authorizing member ID: member_123"

    @pytest.mark.asyncio
    async def test_vpn_id_command(self, command_handle, mock_services):
        """Tests the get VPN ID command."""
        result = await command_handle.process_command("vpn_id", [])

        mock_services["vpn"].get_vpn_id.assert_called_once()
        assert result == "vpn-123456"

    @pytest.mark.asyncio
    async def test_vpn_id_with_error(self, command_handle, mock_services):
        """Tests error handling when getting VPN ID."""
        mock_services["vpn"].get_vpn_id.side_effect = Exception("Service error")

        result = await command_handle.process_command("vpn_id", [])

        assert result == "Error getting vpn id"

    @pytest.mark.asyncio
    async def test_server_ips_command(self, command_handle, mock_services):
        """Tests the get server IPs command."""
        result = await command_handle.process_command("server_ips", [])

        mock_services["vpn"].get_server_ips.assert_called_once()
        assert result == ["192.168.1.1", "192.168.1.2"]

    @pytest.mark.asyncio
    async def test_server_ips_with_error(self, command_handle, mock_services):
        """Tests error handling when getting IPs."""
        mock_services["vpn"].get_server_ips.side_effect = Exception("Network error")

        await command_handle.process_command("server_ips", [])

        # Verifies that error was handled (should not throw exception)
        # Method should return an error message
        pass  # Test passes if no exception is thrown
