use crate::persistence::key_value_storage_factory::StoreType;
use crate::pub_sub::pub_sub_factory::PubSubType;
use crate::vpn::vpn_factory::VpnType;
use getset::{Getters, Setters};
use serde::Deserialize;
use std::{env, time::Duration};

#[derive(Debug, Clone, Deserialize, Getters, Setters)]
#[getset(get = "pub", set = "pub")]
#[serde(default)]
pub struct ApiConfig {
    host: String,
    port: u16,
}

#[derive(Debug, Clone, Deserialize, Getters, Setters)]
#[getset(get = "pub", set = "pub")]
#[serde(default)]
pub struct PubSubConfig {
    url: String,
    pub_sub_type: PubSubType,
}

#[derive(Debug, Clone, Deserialize, Getters, Setters)]
#[getset(get = "pub", set = "pub")]
#[serde(default)]
pub struct VpnConfig {
    api_url: String,
    api_key: String,
    network_id: String,
    vpn_type: VpnType,
}

#[derive(Debug, Clone, Deserialize, Getters, Setters)]
#[getset(get = "pub", set = "pub")]
#[serde(default)]
pub struct DbConfig {
    path: String,
    store_type: StoreType,
}

#[derive(Debug, Clone, Deserialize, Getters, Setters)]
#[getset(get = "pub", set = "pub")]
#[serde(default)]
pub struct AdmineChannelsMap {
    server_channel: String,
    command_channel: String,
    vpn_channel: String,
}

#[derive(Debug, Clone, Deserialize, Getters, Setters)]
#[getset(get = "pub", set = "pub")]
#[serde(default)]
pub struct RetryConfig {
    attempts: usize,
    delay: Duration,
}

#[derive(Debug, Clone, Deserialize, Getters, Setters)]
#[getset(get = "pub", set = "pub")]
#[allow(dead_code)]
#[serde(default)]
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

impl Default for ApiConfig {
    fn default() -> Self {
        Self {
            host: "localhost".to_string(),
            port: 9000,
        }
    }
}

impl Default for PubSubConfig {
    fn default() -> Self {
        Self {
            url: "redis://localhost:6379".to_string(),
            pub_sub_type: PubSubType::Redis,
        }
    }
}

impl Default for VpnConfig {
    fn default() -> Self {
        Self {
            api_url: "https://api.zerotier.com/api/v1".to_string(),
            api_key: "".to_string(),
            network_id: "".to_string(),
            vpn_type: VpnType::Zerotier,
        }
    }
}

impl Default for DbConfig {
    fn default() -> Self {
        Self {
            path: "./etc/sled/vpn_store.db".to_string(),
            store_type: StoreType::Sled,
        }
    }
}

impl Default for AdmineChannelsMap {
    fn default() -> Self {
        Self {
            server_channel: "server_channel".to_string(),
            command_channel: "command_channel".to_string(),
            vpn_channel: "vpn_channel".to_string(),
        }
    }
}

impl Default for RetryConfig {
    fn default() -> Self {
        Self {
            attempts: 5,
            delay: Duration::from_secs(3),
        }
    }
}

impl Default for Config {
    fn default() -> Self {
        Self {
            self_origin_name: "vpn".to_string(),
            api_config: ApiConfig::default(),
            pub_sub_config: PubSubConfig::default(),
            vpn_config: VpnConfig::default(),
            db_config: DbConfig::default(),
            admine_channels_map: AdmineChannelsMap::default(),
            retry_config: RetryConfig::default(),
        }
    }
}
