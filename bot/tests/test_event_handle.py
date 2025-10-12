from unittest.mock import AsyncMock, MagicMock

import pytest

from bot.handles.event_handle import EventHandle
from bot.models.admine_message import AdmineMessage


@pytest.fixture
def mock_message_services():
    """Creates mocks for message services."""
    service1 = MagicMock()
    service1.send_message = AsyncMock()

    service2 = MagicMock()
    service2.send_message = AsyncMock()

    return [service1, service2]


@pytest.fixture
def event_handle(mock_message_services):
    """Creates an EventHandle instance with mocked services."""
    return EventHandle(mock_message_services)


@pytest.fixture
def event_handle_no_services():
    """Creates an EventHandle instance without message services."""
    return EventHandle(None)


class TestEventHandleInitialization:
    """EventHandle initialization tests."""

    def test_initialization_with_services(self, mock_message_services):
        """Verifies that initialization with services works correctly."""
        event_handle = EventHandle(mock_message_services)

        assert event_handle._EventHandle__message_services == mock_message_services

    def test_initialization_without_services(self, event_handle_no_services):
        """Verifies that initialization without services creates an empty list."""
        assert event_handle_no_services._EventHandle__message_services == []


class TestServerEvents:
    """Tests for server related events."""

    @pytest.mark.asyncio
    async def test_server_on_event(self, event_handle, mock_message_services):
        """Tests server_on event processing."""
        event = AdmineMessage("ServerHandler", ["server_on"], "Server started successfully")

        await event_handle.handle_event(event)

        # Verifies that all services were notified
        for service in mock_message_services:
            service.send_message.assert_called_once()
            call_args = service.send_message.call_args[0][0]
            assert "Server has started with message: Server started successfully" in call_args

    @pytest.mark.asyncio
    async def test_server_off_event(self, event_handle, mock_message_services):
        """Tests server_off event processing."""
        event = AdmineMessage("ServerHandler", ["server_off"], "Server stopped gracefully")

        await event_handle.handle_event(event)

        # Verifies that all services were notified
        for service in mock_message_services:
            service.send_message.assert_called_once()
            call_args = service.send_message.call_args[0][0]
            assert "Server has stopped with message: Server stopped gracefully" in call_args

    @pytest.mark.asyncio
    async def test_new_server_ips_event(self, event_handle, mock_message_services):
        """Tests new_server_ips event processing."""
        event = AdmineMessage("VPN", ["new_server_ips"], "192.168.1.1, 192.168.1.2")

        await event_handle.handle_event(event)

        # Verifies that all services were notified
        for service in mock_message_services:
            service.send_message.assert_called_once()
            call_args = service.send_message.call_args[0][0]
            assert "Received new server IPs: 192.168.1.1, 192.168.1.2" in call_args

    @pytest.mark.asyncio
    async def test_notification_event(self, event_handle, mock_message_services):
        """Tests notification event processing."""
        event = AdmineMessage("System", ["notification"], "Important system message")

        await event_handle.handle_event(event)

        # Verifies that all services were notified with the exact message
        for service in mock_message_services:
            service.send_message.assert_called_once()
            call_args = service.send_message.call_args[0][0]
            assert call_args == "Important system message"


class TestUnknownEvents:
    """Tests for unknown events."""

    @pytest.mark.asyncio
    async def test_unknown_tag_handling(self, event_handle, mock_message_services):
        """Verifies that unknown event tags are logged without causing failures."""
        event = AdmineMessage("System", ["unknown_tag"], "Some message")

        # Should not raise an exception
        await event_handle.handle_event(event)

        # No service should be called for unknown tags
        for service in mock_message_services:
            service.send_message.assert_not_called()

    @pytest.mark.asyncio
    async def test_multiple_unknown_tags(self, event_handle, mock_message_services):
        """Verifies behavior with multiple unknown tags."""
        event = AdmineMessage("System", ["unknown1", "unknown2", "unknown3"], "Test message")

        await event_handle.handle_event(event)

        # No service should be called
        for service in mock_message_services:
            service.send_message.assert_not_called()


class TestMultipleTags:
    """Tests for events with multiple tags."""

    @pytest.mark.asyncio
    async def test_event_with_multiple_tags(self, event_handle, mock_message_services):
        """Tests processing of events with multiple tags."""
        event = AdmineMessage("System", ["server_on", "notification"], "Server is now online")

        await event_handle.handle_event(event)

        # Each service should have been called twice (once for each tag)
        for service in mock_message_services:
            assert service.send_message.call_count == 2

    @pytest.mark.asyncio
    async def test_event_with_mixed_known_unknown_tags(self, event_handle, mock_message_services):
        """Tests events with mixed known and unknown tags."""
        event = AdmineMessage("System", ["server_on", "unknown_tag", "notification"], "Mixed tags")

        await event_handle.handle_event(event)

        # Should only be called for known tags (server_on and notification)
        for service in mock_message_services:
            assert service.send_message.call_count == 2


class TestNotifyAll:
    """Tests for notify all services functionality."""

    @pytest.mark.asyncio
    async def test_notify_all_services(self, event_handle, mock_message_services):
        """Verifies that all message services receive the notification."""
        event = AdmineMessage("Test", ["notification"], "Test notification")

        await event_handle.handle_event(event)

        # All services should have been called
        for service in mock_message_services:
            service.send_message.assert_called_once_with("Test notification")

    @pytest.mark.asyncio
    async def test_notify_with_no_services(self, event_handle_no_services):
        """Verifies that there is no error when there are no message services."""
        event = AdmineMessage("Test", ["notification"], "Test message")

        # Should not cause an error
        await event_handle_no_services.handle_event(event)


class TestEventHandleEdgeCases:
    """Edge case tests."""

    @pytest.mark.asyncio
    async def test_empty_tags_list(self, event_handle, mock_message_services):
        """Tests behavior with empty tags list."""
        event = AdmineMessage("System", [], "Message without tags")

        await event_handle.handle_event(event)

        # No service should be called if there are no tags
        for service in mock_message_services:
            service.send_message.assert_not_called()

    @pytest.mark.asyncio
    async def test_empty_message(self, event_handle, mock_message_services):
        """Tests processing of events with empty message."""
        event = AdmineMessage("System", ["notification"], "")

        await event_handle.handle_event(event)

        # Should process normally, even with an empty message
        for service in mock_message_services:
            service.send_message.assert_called_once_with("")

    @pytest.mark.asyncio
    async def test_service_error_handling(self, event_handle, mock_message_services):
        """Tests behavior when a service throws an exception."""
        # Configure a service to raise an exception
        mock_message_services[0].send_message.side_effect = Exception("Service error")

        event = AdmineMessage("System", ["notification"], "Test message")

        # The exception should propagate (EventHandle does not handle service exceptions)
        with pytest.raises(Exception, match="Service error"):
            await event_handle.handle_event(event)
