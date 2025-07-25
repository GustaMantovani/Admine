use crate::persistence::key_value_storage_factory::StoreType;
use crate::pub_sub::pub_sub_factory::PubSubType;
use crate::vpn::vpn_factory::VpnType;
use serde::Deserialize;
use std::time::Duration;

#[derive(Debug, Clone, Deserialize)]
pub struct ApiConfig {
    pub host: String,
    pub port: u16,
}

#[derive(Debug, Clone, Deserialize)]
pub struct PubSubConfig {
    pub url: String,
    pub pub_sub_type: PubSubType,
}

#[derive(Debug, Clone, Deserialize)]
pub struct VpnConfig {
    pub api_url: String,
    pub api_key: String,
    pub network_id: String,
    pub vpn_type: VpnType,
}

#[derive(Debug, Clone, Deserialize)]
pub struct DbConfig {
    pub path: String,
    pub store_type: StoreType,
}

#[derive(Debug, Clone, Deserialize)]
pub struct AdmineChannelsMap {
    pub server_channel: String,
    pub command_channel: String,
    pub vpn_channel: String,
}

#[derive(Debug, Clone, Deserialize)]
pub struct RetryConfig {
    pub attempts: usize,
    pub delay: Duration,
}

#[derive(Debug, Clone, Deserialize)]
pub struct Config {
    pub self_origin_name: String,
    pub api_config: ApiConfig,
    pub pub_sub_config: PubSubConfig,
    pub vpn_config: VpnConfig,
    pub db_config: DbConfig,
    pub admine_channels_map: AdmineChannelsMap,
    pub retry_config: RetryConfig,
}

impl Config {
    pub fn new() -> Result<Self, Box<dyn std::error::Error>> {
        let content = std::fs::read_to_string("./vpn_handler_config.toml")?;
        Ok(toml::from_str(&content)?)
    }
}
