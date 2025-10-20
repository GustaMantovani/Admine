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

    /// Creates a new Config instance with injected parameters (useful for testing and programmatic configuration)
    #[cfg(test)]
    pub fn new_with_params(
        self_origin_name: String,
        api_config: ApiConfig,
        pub_sub_config: PubSubConfig,
        vpn_config: VpnConfig,
        db_config: DbConfig,
        admine_channels_map: AdmineChannelsMap,
        retry_config: RetryConfig,
    ) -> Self {
        Self {
            self_origin_name,
            api_config,
            pub_sub_config,
            vpn_config,
            db_config,
            admine_channels_map,
            retry_config,
        }
    }

    /// Creates a default Config instance for testing purposes
    #[cfg(test)]
    pub fn new_test_config() -> Self {
        let api_config = ApiConfig {
            host: "localhost".to_string(),
            port: 8080,
        };

        let pub_sub_config = PubSubConfig {
            url: "redis://localhost:6379".to_string(),
            pub_sub_type: PubSubType::Redis,
        };

        let vpn_config = VpnConfig {
            api_url: "https://api.zerotier.com".to_string(),
            api_key: "test_api_key".to_string(),
            network_id: "test_network_id".to_string(),
            vpn_type: VpnType::Zerotier,
        };

        let db_config = DbConfig {
            path: "./test_data".to_string(),
            store_type: StoreType::Sled,
        };

        let admine_channels_map = AdmineChannelsMap {
            server_channel: "server_channel".to_string(),
            command_channel: "command_channel".to_string(),
            vpn_channel: "vpn_channel".to_string(),
        };

        let retry_config = RetryConfig {
            attempts: 3,
            delay: Duration::from_millis(100),
        };

        Self::new_with_params(
            "test_vpn_handler".to_string(),
            api_config,
            pub_sub_config,
            vpn_config,
            db_config,
            admine_channels_map,
            retry_config,
        )
    }
}

impl ApiConfig {
    #[cfg(test)]
    pub fn new(host: String, port: u16) -> Self {
        Self { host, port }
    }
}

impl PubSubConfig {
    #[cfg(test)]
    pub fn new(url: String, pub_sub_type: PubSubType) -> Self {
        Self { url, pub_sub_type }
    }
}

impl VpnConfig {
    #[cfg(test)]
    pub fn new(api_url: String, api_key: String, network_id: String, vpn_type: VpnType) -> Self {
        Self {
            api_url,
            api_key,
            network_id,
            vpn_type,
        }
    }
}

impl DbConfig {
    #[cfg(test)]
    pub fn new(path: String, store_type: StoreType) -> Self {
        Self { path, store_type }
    }
}

impl AdmineChannelsMap {
    #[cfg(test)]
    pub fn new(server_channel: String, command_channel: String, vpn_channel: String) -> Self {
        Self {
            server_channel,
            command_channel,
            vpn_channel,
        }
    }
}

impl RetryConfig {
    #[cfg(test)]
    pub fn new(attempts: usize, delay: Duration) -> Self {
        Self { attempts, delay }
    }
}
