use crate::persistence::key_value_storage_factory::StoreType;
use crate::pub_sub::pub_sub_factory::PubSubType;
use crate::vpn::vpn_factory::VpnType;
use getset::{Getters, Setters};
use serde::Deserialize;
use std::{env, time::Duration};

#[derive(Debug, Clone, Deserialize, Getters, Setters)]
#[getset(get = "pub", set = "pub")]
pub struct ApiConfig {
    host: String,
    port: u16,
}

#[derive(Debug, Clone, Deserialize, Getters, Setters)]
#[getset(get = "pub", set = "pub")]
pub struct PubSubConfig {
    url: String,
    pub_sub_type: PubSubType,
}

#[derive(Debug, Clone, Deserialize, Getters, Setters)]
#[getset(get = "pub", set = "pub")]
pub struct VpnConfig {
    api_url: String,
    api_key: String,
    network_id: String,
    vpn_type: VpnType,
}

#[derive(Debug, Clone, Deserialize, Getters, Setters)]
#[getset(get = "pub", set = "pub")]
pub struct DbConfig {
    path: String,
    store_type: StoreType,
}

#[derive(Debug, Clone, Deserialize, Getters, Setters)]
#[getset(get = "pub", set = "pub")]
pub struct AdmineChannelsMap {
    server_channel: String,
    command_channel: String,
    vpn_channel: String,
}

#[derive(Debug, Clone, Deserialize, Getters, Setters)]
#[getset(get = "pub", set = "pub")]
pub struct RetryConfig {
    attempts: usize,
    delay: Duration,
}

#[derive(Debug, Clone, Deserialize, Getters, Setters)]
#[getset(get = "pub", set = "pub")]
#[allow(dead_code)]
pub struct Config {
    self_origin_name: String,
    api_config: ApiConfig,
    pub_sub_config: PubSubConfig,
    vpn_config: VpnConfig,
    db_config: DbConfig,
    admine_channels_map: AdmineChannelsMap,
    retry_config: RetryConfig,
}

impl Config {
    pub fn new() -> Result<Self, Box<dyn std::error::Error>> {
        let args: Vec<String> = env::args().collect();

        let config_path = if args.len() > 1 {
            args[1].clone()
        } else {
            String::from("./etc/vpn_handler_config.toml")
        };

        let content = std::fs::read_to_string(&config_path)?;
        Ok(toml::from_str(&content)?)
    }
}

impl RetryConfig {
    #[cfg(test)]
    pub fn new(attempts: usize, delay: Duration) -> Self {
        Self { attempts, delay }
    }
}
