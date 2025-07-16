use crate::config::Config;
use crate::persistence::key_value_storage_factory::{DynKeyValueStore, StoreFactory};
use crate::vpn::vpn::TVpnClient;
use crate::vpn::vpn_factory::VpnFactory;
use anyhow::Result;
use std::sync::{Mutex, OnceLock};

pub struct AppContext {
    config: Config,
    storage: Mutex<DynKeyValueStore>,
    vpn_client: Box<dyn TVpnClient + Send + Sync>,
}

static APP_CONTEXT: OnceLock<AppContext> = OnceLock::new();

impl AppContext {
    fn new() -> Result<Self> {
        // Ordered component initialization
        let config = Config::new()?;
        let storage = StoreFactory::create_store_instance(
            config.db_config.store_type.clone(),
            &config.db_config.path,
        )
        .map_err(|e| anyhow::anyhow!("Failed to create storage: {}", e))?;

        // Create VPN client based on configuration
        let vpn_client = VpnFactory::create_vpn(
            config.vpn_config.vpn_type.clone(),
            config.vpn_config.api_url.clone(),
            config.vpn_config.api_key.clone(),
            config.vpn_config.network_id.clone(),
        )
        .map_err(|e| anyhow::anyhow!("Failed to create VPN client: {:?}", e))?;

        Ok(Self {
            config,
            storage: Mutex::new(storage),
            vpn_client,
        })
    }

    pub fn instance() -> &'static AppContext {
        APP_CONTEXT.get_or_init(|| Self::new().expect("Failed to initialize application context"))
    }

    // Component accessors
    pub fn config(&self) -> &Config {
        &self.config
    }

    pub fn storage(&self) -> &Mutex<DynKeyValueStore> {
        &self.storage
    }

    pub fn vpn_client(&self) -> &Box<dyn TVpnClient + Send + Sync> {
        &self.vpn_client
    }

    // Storage methods
    pub fn set_storage(&self, key: String, value: String) -> Result<(), String> {
        let guard = self.storage.lock().unwrap();
        guard.set(key, value)
    }

    pub fn get_storage(&self, key: &str) -> Result<Option<String>, String> {
        let guard = self.storage.lock().unwrap();
        guard.get(key)
    }
}
