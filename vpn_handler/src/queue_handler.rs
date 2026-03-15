use std::sync::Arc;

use crate::config::{Config, RetryConfig};
use crate::models::admine_message::AdmineMessage;
use crate::persistence::key_value_storage::DynKeyValueStore;
use crate::pub_sub::pub_sub::DynPubSub;
use crate::vpn::vpn::DynVpn;
use log::{error, info, warn};
use tokio::sync::watch;
use tokio::time::sleep;

/// Main handle structure.
pub struct Handle {
    pub_sub: DynPubSub,
    vpn_client: Arc<DynVpn>,
    storage: Arc<DynKeyValueStore>,
    config: Arc<Config>,
    shutdown: watch::Receiver<bool>,
}

impl Handle {
    pub fn new(
        pub_sub: DynPubSub,
        vpn_client: Arc<DynVpn>,
        storage: Arc<DynKeyValueStore>,
        config: Arc<Config>,
        shutdown: watch::Receiver<bool>,
    ) -> Self {
        Self {
            pub_sub,
            vpn_client,
            storage,
            config,
            shutdown,
        }
    }

    /// Helper function to update the server member ID in the database.
    async fn update_server_id(
        new_id: &str,
        storage: &DynKeyValueStore,
    ) -> Result<(), Box<dyn std::error::Error>> {
        storage
            .set("server_member_id".to_string(), new_id.to_string())
            .map_err(|e| {
                error!("Error saving new server member id: {}", e);
                e.into()
            })
    }

    /// Process server_up messages with retry logic and IP fetching
    async fn process_server_up(
        member_id: String,
        vpn_client: &DynVpn,
        storage: &DynKeyValueStore,
        pub_sub: &mut DynPubSub,
        retry_config: &RetryConfig,
        vpn_channel: &str,
        origin: &str,
    ) {
        // Retry logic to authenticate member and fetch IPs until available
        let mut attempts = *retry_config.attempts();
        let member_ips = loop {
            // First, try to authenticate the member
            if let Err(e) = vpn_client.auth_member(member_id.clone(), None).await {
                if attempts == 0 {
                    error!(
                        "Exceeded retry attempts to authenticate member {}: {}",
                        member_id, e
                    );
                    return;
                }
                attempts -= 1;
                error!(
                    "Error authenticating member {}: {}. Retrying in {:?}... (attempts left: {})",
                    member_id,
                    e,
                    retry_config.delay(),
                    attempts
                );
                sleep(*retry_config.delay()).await;
                continue;
            }

            info!("Member {} authenticated successfully.", member_id);

            // Then try to fetch the IPs
            match vpn_client.get_member_ips_in_vpn(member_id.clone()).await {
                Ok(ips) => {
                    if ips.is_empty() {
                        if attempts == 0 {
                            error!(
                                "Exceeded retry attempts to fetch IPs for member {}",
                                member_id
                            );
                            return;
                        }
                        attempts -= 1;
                        info!(
                            "No IPs available yet for member {}. Retrying in {:?}... (attempts left: {})",
                            member_id,
                            retry_config.delay(),
                            attempts
                        );
                        sleep(*retry_config.delay()).await;
                        continue; // Continue the loop to retry
                    }

                    break ips; // Only break when we have IPs
                }
                Err(e) => {
                    if attempts == 0 {
                        error!(
                            "Exceeded retry attempts to fetch IPs for member {}: {}",
                            member_id, e
                        );
                        return;
                    }
                    attempts -= 1;
                    info!(
                        "IPs not available yet for member {}. Retrying in {:?}... (attempts left: {})",
                        member_id,
                        retry_config.delay(),
                        attempts
                    );
                    sleep(*retry_config.delay()).await;
                }
            }
        };

        // Publish new server IPs
        let new_message = AdmineMessage::new(
            origin.to_string(),
            vec!["new_server_ips".to_string()],
            member_ips
                .iter()
                .map(|ip| ip.to_string())
                .collect::<Vec<String>>()
                .join(","),
        );

        let serialized_message = match serde_json::to_string(&new_message) {
            Ok(json) => json,
            Err(e) => {
                error!("Failed to serialize message: {}", e);
                return;
            }
        };

        if let Err(e) = pub_sub
            .publish(vpn_channel.to_string(), serialized_message)
            .await
        {
            error!("Failed to publish message: {}", e);
        } else {
            info!("New server up message published successfully.");
        }

        // Handle old server cleanup
        let old_member_id = storage.get("server_member_id").unwrap_or_default();

        if !old_member_id.is_empty() && old_member_id != member_id {
            if let Err(e) = vpn_client.delete_member(old_member_id.clone()).await {
                error!("Error deleting old member {}: {}", old_member_id, e);
            }
        }

        // Save the new server member ID
        if let Err(e) = Self::update_server_id(&member_id, storage).await {
            error!("Failed to update server member id: {}", e);
        }
    }

    /// Process auth_member command
    async fn process_auth_member(
        member_id: String,
        vpn_client: &DynVpn,
        pub_sub: &mut DynPubSub,
        vpn_channel: &str,
        origin: &str,
    ) {
        if let Err(e) = vpn_client.auth_member(member_id.clone(), None).await {
            error!("Error authenticating member {}: {}", member_id, e);
            return;
        }
        info!("Member {} authenticated successfully.", member_id);

        // Publish success message
        let success_message = AdmineMessage::new(
            origin.to_string(),
            vec!["auth_member_success".to_string()],
            member_id,
        );

        let serialized_message = match serde_json::to_string(&success_message) {
            Ok(json) => json,
            Err(e) => {
                error!("Failed to serialize success message: {}", e);
                return;
            }
        };

        if let Err(e) = pub_sub
            .publish(vpn_channel.to_string(), serialized_message)
            .await
        {
            error!("Failed to publish success message: {}", e);
        } else {
            info!("Auth member success message published successfully.");
        }
    }

    /// Process incoming messages based on channel and tags (with injected dependencies)
    async fn process_message_with_deps(
        admine_message: AdmineMessage,
        vpn_client: &DynVpn,
        storage: &DynKeyValueStore,
        pub_sub: &mut DynPubSub,
        retry_config: &RetryConfig,
        vpn_channel: &str,
        origin: &str,
    ) {
        info!(
            "Dispatching message: origin={}, message_len={}",
            admine_message.origin(),
            admine_message.message().len()
        );
        match admine_message.origin() {
            // Server channel - handle server_up messages
            org if org == "server" => {
                if admine_message.has_tag("server_on") && !admine_message.message().is_empty() {
                    let member_id = admine_message.message().clone();
                    info!("server_up received: member_id={}", member_id);
                    Self::process_server_up(
                        member_id,
                        vpn_client,
                        storage,
                        pub_sub,
                        retry_config,
                        vpn_channel,
                        origin,
                    )
                    .await;
                } else {
                    warn!("Ignored server message...");
                }
            }
            // Command channel - handle auth_member commands
            org if org == "bot" => {
                if admine_message.has_tag("auth_member") && !admine_message.message().is_empty() {
                    let member_id = admine_message.message().clone();
                    info!("auth_member command received: member_id={}", member_id);
                    Self::process_auth_member(member_id, vpn_client, pub_sub, vpn_channel, origin)
                        .await;
                } else {
                    warn!("Ignored bot command...");
                }
            }
            other => {
                warn!("Unsupported channel: {}", other);
            }
        }
    }

    /// Main run loop.
    pub async fn run(mut self) {
        info!("Handle run started.");

        let retry_config = self.config.retry_config().clone();
        let vpn_channel = self.config.admine_channels_map().vpn_channel().clone();
        let origin = self.config.self_origin_name().clone();

        loop {
            info!("Waiting for a new message...");

            let raw_message = tokio::select! {
                biased;

                _ = self.shutdown.changed() => {
                    if *self.shutdown.borrow() {
                        info!("Shutdown signal received, stopping queue handler.");
                        break;
                    }
                    continue;
                }

                result = self.pub_sub.listen_until_receive_message() => result,
            };

            let (payload, channel) = match raw_message {
                Ok(msg) => {
                    info!("Message received on channel {}: {:?}", msg.1, msg.0);
                    msg
                }
                Err(e) => {
                    error!("Error receiving message: {}", e);
                    continue;
                }
            };

            let admine_message = match serde_json::from_str::<AdmineMessage>(&payload) {
                Ok(msg) => msg,
                Err(e) => {
                    error!("Error deserializing message: {}", e);
                    continue;
                }
            };

            info!(
                "Processing message received on channel {}: {:?}",
                channel, admine_message
            );

            Self::process_message_with_deps(
                admine_message,
                &self.vpn_client,
                &self.storage,
                &mut self.pub_sub,
                &retry_config,
                &vpn_channel,
                &origin,
            )
            .await;
        }

        info!("Queue handler stopped.");
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::config::RetryConfig;
    use crate::errors::{PubSubError, VpnError};
    use crate::models::admine_message::AdmineMessage;
    use crate::persistence::key_value_storage::KeyValueStore;
    use crate::pub_sub::pub_sub::{TPublisher, TSubscriber};
    use crate::vpn::vpn::TVpnClient;
    use async_trait::async_trait;
    use mockall::{mock, predicate::*};
    use std::net::IpAddr;
    use std::sync::Arc;
    use std::time::Duration;
    use tokio::sync::watch;

    // Mock implementations
    mock! {
        TestVpn {}

        #[async_trait]
        impl TVpnClient for TestVpn {
            async fn auth_member(&self, member_id: String, member_token: Option<String>) -> Result<(), VpnError>;
            async fn delete_member(&self, member_id: String) -> Result<(), VpnError>;
            async fn get_member_ips_in_vpn(&self, member_id: String) -> Result<Vec<IpAddr>, VpnError>;
        }
    }

    mock! {
        TestStorage {}

        impl KeyValueStore for TestStorage {
            fn set(&self, key: String, value: String) -> Result<(), Box<dyn std::error::Error>>;
            fn get(&self, key: &str) -> Option<String>;
        }
    }

    mock! {
        TestPubSub {}

        #[async_trait]
        impl TSubscriber for TestPubSub {
            fn subscribe(&mut self, topics: Vec<String>) -> Result<(), PubSubError>;
            async fn listen_until_receive_message(&mut self) -> Result<(String, String), PubSubError>;
        }

        #[async_trait]
        impl TPublisher for TestPubSub {
            async fn publish(&mut self, topic: String, message: String) -> Result<(), PubSubError>;
        }
    }

    fn create_test_retry_config() -> RetryConfig {
        RetryConfig::new(3, Duration::from_millis(10))
    }

    #[tokio::test]
    async fn test_update_server_id_success() {
        let mut mock_storage = MockTestStorage::new();

        mock_storage
            .expect_set()
            .with(
                eq("server_member_id".to_string()),
                eq("test_member_123".to_string()),
            )
            .times(1)
            .returning(|_, _| Ok(()));

        let storage: Box<dyn KeyValueStore + Send + Sync> = Box::new(mock_storage);

        let result = Handle::update_server_id("test_member_123", &storage).await;

        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_update_server_id_failure() {
        let mut mock_storage = MockTestStorage::new();

        mock_storage
            .expect_set()
            .with(
                eq("server_member_id".to_string()),
                eq("test_member_123".to_string()),
            )
            .times(1)
            .returning(|_, _| Err("Database error".into()));

        let storage: Box<dyn KeyValueStore + Send + Sync> = Box::new(mock_storage);

        let result = Handle::update_server_id("test_member_123", &storage).await;

        assert!(result.is_err());
    }

    #[tokio::test]
    async fn test_process_auth_member_success() {
        let mut mock_vpn = MockTestVpn::new();
        let mut mock_pubsub = MockTestPubSub::new();

        mock_vpn
            .expect_auth_member()
            .with(eq("test_member_123".to_string()), eq(None))
            .times(1)
            .returning(|_, _| Ok(()));

        mock_pubsub
            .expect_publish()
            .with(eq("vpn_channel".to_string()), always())
            .times(1)
            .returning(|_, _| Ok(()));

        let vpn: Box<dyn TVpnClient + Send + Sync> = Box::new(mock_vpn);
        let mut pubsub: DynPubSub = Box::new(mock_pubsub);

        Handle::process_auth_member(
            "test_member_123".to_string(),
            &vpn,
            &mut pubsub,
            "vpn_channel",
            "vpn",
        )
        .await;
    }

    #[tokio::test]
    async fn test_process_auth_member_auth_failure() {
        let mut mock_vpn = MockTestVpn::new();
        let mut mock_pubsub = MockTestPubSub::new();

        mock_vpn
            .expect_auth_member()
            .with(eq("test_member_123".to_string()), eq(None))
            .times(1)
            .returning(|_, _| Err(VpnError::InternalError("Auth failed".to_string())));

        // Should not publish when auth fails
        mock_pubsub.expect_publish().times(0);

        let vpn: Box<dyn TVpnClient + Send + Sync> = Box::new(mock_vpn);
        let mut pubsub: DynPubSub = Box::new(mock_pubsub);

        Handle::process_auth_member(
            "test_member_123".to_string(),
            &vpn,
            &mut pubsub,
            "vpn_channel",
            "vpn",
        )
        .await;
    }

    #[tokio::test]
    async fn test_process_server_up_success() {
        let mut mock_vpn = MockTestVpn::new();
        let mut mock_storage = MockTestStorage::new();
        let mut mock_pubsub = MockTestPubSub::new();

        let test_ips = vec![
            "192.168.1.100".parse().unwrap(),
            "10.0.0.50".parse().unwrap(),
        ];

        mock_vpn
            .expect_auth_member()
            .with(eq("test_member_123".to_string()), eq(None))
            .times(1)
            .returning(|_, _| Ok(()));

        mock_vpn
            .expect_get_member_ips_in_vpn()
            .with(eq("test_member_123".to_string()))
            .times(1)
            .returning(move |_| Ok(test_ips.clone()));

        mock_storage
            .expect_get()
            .with(eq("server_member_id"))
            .times(1)
            .returning(|_| None);

        mock_storage
            .expect_set()
            .with(
                eq("server_member_id".to_string()),
                eq("test_member_123".to_string()),
            )
            .times(1)
            .returning(|_, _| Ok(()));

        mock_pubsub
            .expect_publish()
            .with(eq("vpn_channel".to_string()), always())
            .times(1)
            .returning(|_, _| Ok(()));

        let vpn: Box<dyn TVpnClient + Send + Sync> = Box::new(mock_vpn);
        let storage: Box<dyn KeyValueStore + Send + Sync> = Box::new(mock_storage);
        let mut pubsub: DynPubSub = Box::new(mock_pubsub);
        let retry_config = create_test_retry_config();

        Handle::process_server_up(
            "test_member_123".to_string(),
            &vpn,
            &storage,
            &mut pubsub,
            &retry_config,
            "vpn_channel",
            "vpn",
        )
        .await;
    }

    #[tokio::test]
    async fn test_process_server_up_with_old_member_cleanup() {
        let mut mock_vpn = MockTestVpn::new();
        let mut mock_storage = MockTestStorage::new();
        let mut mock_pubsub = MockTestPubSub::new();

        let test_ips = vec!["192.168.1.100".parse().unwrap()];

        mock_vpn.expect_auth_member().returning(|_, _| Ok(()));

        mock_vpn
            .expect_get_member_ips_in_vpn()
            .returning(move |_| Ok(test_ips.clone()));

        mock_vpn
            .expect_delete_member()
            .with(eq("old_member_456".to_string()))
            .times(1)
            .returning(|_| Ok(()));

        mock_storage
            .expect_get()
            .with(eq("server_member_id"))
            .times(1)
            .returning(|_| Some("old_member_456".to_string()));

        mock_storage.expect_set().returning(|_, _| Ok(()));

        mock_pubsub.expect_publish().returning(|_, _| Ok(()));

        let vpn: Box<dyn TVpnClient + Send + Sync> = Box::new(mock_vpn);
        let storage: Box<dyn KeyValueStore + Send + Sync> = Box::new(mock_storage);
        let mut pubsub: DynPubSub = Box::new(mock_pubsub);
        let retry_config = create_test_retry_config();

        Handle::process_server_up(
            "test_member_123".to_string(),
            &vpn,
            &storage,
            &mut pubsub,
            &retry_config,
            "vpn_channel",
            "vpn",
        )
        .await;
    }

    #[tokio::test]
    async fn test_process_server_up_retry_on_no_ips() {
        let mut mock_vpn = MockTestVpn::new();
        let mut mock_storage = MockTestStorage::new();
        let mut mock_pubsub = MockTestPubSub::new();

        let empty_ips: Vec<IpAddr> = vec![];
        let final_ips = vec!["192.168.1.100".parse().unwrap()];

        mock_vpn.expect_auth_member().returning(|_, _| Ok(()));

        // First call returns empty IPs, second call returns actual IPs
        mock_vpn.expect_get_member_ips_in_vpn().times(2).returning({
            let mut call_count = 0;
            move |_| {
                call_count += 1;
                if call_count == 1 {
                    Ok(empty_ips.clone())
                } else {
                    Ok(final_ips.clone())
                }
            }
        });

        mock_storage.expect_get().returning(|_| None);
        mock_storage.expect_set().returning(|_, _| Ok(()));
        mock_pubsub.expect_publish().returning(|_, _| Ok(()));

        let vpn: Box<dyn TVpnClient + Send + Sync> = Box::new(mock_vpn);
        let storage: Box<dyn KeyValueStore + Send + Sync> = Box::new(mock_storage);
        let mut pubsub: DynPubSub = Box::new(mock_pubsub);
        let retry_config = create_test_retry_config();

        Handle::process_server_up(
            "test_member_123".to_string(),
            &vpn,
            &storage,
            &mut pubsub,
            &retry_config,
            "vpn_channel",
            "vpn",
        )
        .await;
    }

    #[tokio::test]
    async fn test_process_message_with_deps_server_on() {
        let mut mock_vpn = MockTestVpn::new();
        let mut mock_storage = MockTestStorage::new();
        let mut mock_pubsub = MockTestPubSub::new();

        let test_ips = vec!["192.168.1.100".parse().unwrap()];

        mock_vpn.expect_auth_member().returning(|_, _| Ok(()));
        mock_vpn
            .expect_get_member_ips_in_vpn()
            .returning(move |_| Ok(test_ips.clone()));
        mock_storage.expect_get().returning(|_| None);
        mock_storage.expect_set().returning(|_, _| Ok(()));
        mock_pubsub.expect_publish().returning(|_, _| Ok(()));

        let message = AdmineMessage::new(
            "server".to_string(),
            vec!["server_on".to_string()],
            "test_member_123".to_string(),
        );

        let vpn: Box<dyn TVpnClient + Send + Sync> = Box::new(mock_vpn);
        let storage: Box<dyn KeyValueStore + Send + Sync> = Box::new(mock_storage);
        let mut pubsub: DynPubSub = Box::new(mock_pubsub);
        let retry_config = create_test_retry_config();

        Handle::process_message_with_deps(
            message,
            &vpn,
            &storage,
            &mut pubsub,
            &retry_config,
            "vpn_channel",
            "vpn",
        )
        .await;
    }

    #[tokio::test]
    async fn test_process_message_with_deps_auth_member() {
        let mut mock_vpn = MockTestVpn::new();
        let mock_storage = MockTestStorage::new();
        let mut mock_pubsub = MockTestPubSub::new();

        mock_vpn
            .expect_auth_member()
            .with(eq("test_member_123".to_string()), eq(None))
            .times(1)
            .returning(|_, _| Ok(()));

        mock_pubsub
            .expect_publish()
            .times(1)
            .returning(|_, _| Ok(()));

        let message = AdmineMessage::new(
            "bot".to_string(),
            vec!["auth_member".to_string()],
            "test_member_123".to_string(),
        );

        let vpn: Box<dyn TVpnClient + Send + Sync> = Box::new(mock_vpn);
        let storage: Box<dyn KeyValueStore + Send + Sync> = Box::new(mock_storage);
        let mut pubsub: DynPubSub = Box::new(mock_pubsub);
        let retry_config = create_test_retry_config();

        Handle::process_message_with_deps(
            message,
            &vpn,
            &storage,
            &mut pubsub,
            &retry_config,
            "vpn_channel",
            "vpn",
        )
        .await;
    }

    #[tokio::test]
    async fn test_process_message_with_deps_unsupported_channel() {
        let mock_vpn = MockTestVpn::new();
        let mock_storage = MockTestStorage::new();
        let mock_pubsub = MockTestPubSub::new();

        let message = AdmineMessage::new(
            "unknown_channel".to_string(),
            vec!["some_tag".to_string()],
            "test_message".to_string(),
        );

        let vpn: Box<dyn TVpnClient + Send + Sync> = Box::new(mock_vpn);
        let storage: Box<dyn KeyValueStore + Send + Sync> = Box::new(mock_storage);
        let mut pubsub: DynPubSub = Box::new(mock_pubsub);
        let retry_config = create_test_retry_config();

        Handle::process_message_with_deps(
            message,
            &vpn,
            &storage,
            &mut pubsub,
            &retry_config,
            "vpn_channel",
            "vpn",
        )
        .await;
        // This test mainly checks that the function doesn't panic with unknown channels
    }

    #[tokio::test]
    async fn test_process_message_with_deps_empty_message() {
        let mock_vpn = MockTestVpn::new();
        let mock_storage = MockTestStorage::new();
        let mock_pubsub = MockTestPubSub::new();

        let message = AdmineMessage::new(
            "server".to_string(),
            vec!["server_on".to_string()],
            "".to_string(), // Empty message
        );

        let vpn: Box<dyn TVpnClient + Send + Sync> = Box::new(mock_vpn);
        let storage: Box<dyn KeyValueStore + Send + Sync> = Box::new(mock_storage);
        let mut pubsub: DynPubSub = Box::new(mock_pubsub);
        let retry_config = create_test_retry_config();

        Handle::process_message_with_deps(
            message,
            &vpn,
            &storage,
            &mut pubsub,
            &retry_config,
            "vpn_channel",
            "vpn",
        )
        .await;
        // This test checks that empty messages are ignored
    }

    #[tokio::test]
    async fn test_handle_shutdown_on_watch_signal() {
        // No expectations on pub_sub: with biased select!, shutdown is checked first.
        // Since the signal is sent before run() starts, listen_until_receive_message
        // is never polled.
        let mock_vpn = MockTestVpn::new();
        let mock_storage = MockTestStorage::new();
        let mock_pubsub = MockTestPubSub::new();

        let (shutdown_tx, shutdown_rx) = watch::channel(false);

        // Signal shutdown before run() is called — the watch receiver has not yet
        // observed this change, so changed() resolves on the very first poll.
        let _ = shutdown_tx.send(true);

        let config = Arc::new(crate::config::Config::default());
        let vpn_client = Arc::new(Box::new(mock_vpn) as Box<dyn TVpnClient + Send + Sync>);
        let storage = Arc::new(Box::new(mock_storage) as Box<dyn KeyValueStore + Send + Sync>);
        let pubsub: DynPubSub = Box::new(mock_pubsub);

        let handle = Handle::new(pubsub, vpn_client, storage, config, shutdown_rx);

        // run() should return immediately without touching pub_sub
        tokio::time::timeout(Duration::from_secs(1), handle.run())
            .await
            .expect("Handle::run() did not stop after shutdown signal");
    }
}
