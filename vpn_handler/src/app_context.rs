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

static APP_CONTEXT: OnceLock<AppContext> = OnceLock::new();

impl AppContext {
    fn new() -> Result<Self, String> {
        // Ordered component initialization
        let config = Config::new().map_err(|e| format!("Failed to load config: {}", e))?;
        let storage = StoreFactory::create_store_instance(
            config.db_config().store_type().clone(),
            config.db_config().path(),
        )?;

        // Create VPN client based on configuration
        let vpn_client = VpnFactory::create_vpn(
            config.vpn_config().vpn_type().clone(),
            config.vpn_config().api_url().clone(),
            config.vpn_config().api_key().clone(),
            config.vpn_config().network_id().clone(),
        )
        .map_err(|e| format!("Failed to create VPN client: {:?}", e))?;

        // Create PubSub client based on configuration
        let mut pub_sub = PubSubFactory::create_pubsub_instance(
            config.pub_sub_config().pub_sub_type().clone(),
            config.pub_sub_config().url(),
        )
        .map_err(|e| format!("Failed to create PubSub client: {:?}", e))?;

        // Subscribe to channels
        pub_sub
            .subscribe(vec![
                config.admine_channels_map().server_channel().clone(),
                config.admine_channels_map().command_channel().clone(),
            ])
            .map_err(|e| format!("Failed to subscribe to channels: {:?}", e))?;

        Ok(Self {
            config,
            storage,
            vpn_client,
            pub_sub: Mutex::new(pub_sub),
        })
    }

    pub fn instance() -> &'static AppContext {
        APP_CONTEXT.get_or_init(|| Self::new().expect("Failed to initialize application context"))
    }
}
