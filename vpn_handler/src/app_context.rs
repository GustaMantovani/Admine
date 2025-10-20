use getset::{Getters, Setters};

use crate::config::Config;
use crate::persistence::key_value_storage::DynKeyValueStore;
use crate::persistence::key_value_storage_factory::StoreFactory;
use crate::pub_sub::pub_sub::DynPubSub;
use crate::pub_sub::pub_sub_factory::PubSubFactory;
use crate::vpn::vpn::DynVpn;
use crate::vpn::vpn_factory::VpnFactory;
use std::sync::{Mutex, OnceLock};

#[derive(Getters, Setters)]
#[getset(get = "pub", set = "pub")]
pub struct AppContext {
    config: Config,
    storage: DynKeyValueStore,
    vpn_client: DynVpn,
    pub_sub: Mutex<DynPubSub>,
}

/// Builder for creating AppContext instances flexibly.
pub struct AppContextBuilder {
    config: Option<Config>,
    storage: Option<DynKeyValueStore>,
    vpn_client: Option<DynVpn>,
    pub_sub: Option<DynPubSub>,
    subscribe_to_channels: bool, // To control subscription to channels
}

static APP_CONTEXT: OnceLock<AppContext> = OnceLock::new();

impl AppContextBuilder {
    /// Creates a new empty builder.
    pub fn new() -> Self {
        Self {
            config: None,
            storage: None,
            vpn_client: None,
            pub_sub: None,
            subscribe_to_channels: true, // By default, the real application subscribes.
        }
    }

    /// Sets a custom configuration.
    pub fn with_config(mut self, config: Config) -> Self {
        self.config = Some(config);
        self
    }

    /// Sets a custom storage (useful for injecting mocks).
    pub fn with_storage(mut self, storage: DynKeyValueStore) -> Self {
        self.storage = Some(storage);
        self
    }

    /// Sets a custom VPN client.
    pub fn with_vpn_client(mut self, vpn_client: DynVpn) -> Self {
        self.vpn_client = Some(vpn_client);
        self
    }

    /// Sets a custom PubSub.
    pub fn with_pub_sub(mut self, pub_sub: DynPubSub) -> Self {
        self.pub_sub = Some(pub_sub);
        self
    }

    /// Sets whether PubSub should subscribe to channels when being built.
    pub fn subscribe_to_channels(mut self, subscribe: bool) -> Self {
        self.subscribe_to_channels = subscribe;
        self
    }

    /// Builds the AppContext with the provided components or creates defaults.
    pub fn build(self) -> Result<AppContext, String> {
        // 1. Config: uses provided or loads from file.
        let config = self.config.map_or_else(
            || Config::new().map_err(|e| format!("Failed to load config: {}", e)),
            Ok,
        )?;

        // 2. Storage: uses provided or creates based on config.
        let storage = self.storage.map_or_else(
            || {
                StoreFactory::create_store_instance(
                    config.db_config().store_type().clone(),
                    config.db_config().path(),
                )
            },
            Ok,
        )?;

        // 3. VPN Client: uses provided or creates based on config.
        let vpn_client = self.vpn_client.map_or_else(
            || {
                VpnFactory::create_vpn(
                    config.vpn_config().vpn_type().clone(),
                    config.vpn_config().api_url().clone(),
                    config.vpn_config().api_key().clone(),
                    config.vpn_config().network_id().clone(),
                )
                .map_err(|e| format!("Failed to create VPN client: {}", e))
            },
            Ok,
        )?;

        // 4. PubSub: uses provided or creates based on config.
        let mut pub_sub = self.pub_sub.map_or_else(
            || {
                PubSubFactory::create_pubsub_instance(
                    config.pub_sub_config().pub_sub_type().clone(),
                    config.pub_sub_config().url(),
                )
                .map_err(|e| format!("Failed to create PubSub client: {}", e))
            },
            Ok,
        )?;

        // 5. Channel subscription logic (only if enabled).
        if self.subscribe_to_channels {
            let channels = vec![
                config.admine_channels_map().server_channel().clone(),
                config.admine_channels_map().command_channel().clone(),
                config.admine_channels_map().vpn_channel().clone(),
            ];

            if let Err(e) = pub_sub.subscribe(channels.clone()) {
                log::warn!("Failed to subscribe to channels {:?}: {}", channels, e);
            }
        }

        Ok(AppContext {
            config,
            storage,
            vpn_client,
            pub_sub: Mutex::new(pub_sub),
        })
    }
}

// Default implementation to facilitate builder creation
impl Default for AppContextBuilder {
    fn default() -> Self {
        Self::new()
    }
}

impl AppContext {
    /// The original `new` method now uses the builder with default configuration.
    fn new() -> Result<Self, String> {
        AppContextBuilder::new().build()
    }

    /// The singleton instance for the main application.
    pub fn instance() -> &'static AppContext {
        APP_CONTEXT.get_or_init(|| {
            Self::new().unwrap_or_else(|e| {
                log::error!("Failed to initialize application context: {}", e);
                panic!("Failed to initialize application context: {}", e);
            })
        })
    }

    /// Entry point for creating a custom context. This is what we'll use in tests!
    pub fn builder() -> AppContextBuilder {
        AppContextBuilder::new()
    }
}
