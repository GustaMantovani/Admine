use dotenvy::dotenv;
use log::{error, info};
use serde::{Deserialize, Serialize};
use std::env;
use std::fs;
use std::path::Path;
use std::str::FromStr;
use std::time::Duration;

use crate::persistence::factories::StoreType;
use crate::pub_sub::factories::PubSubType;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PubSubConfig {
    pub url: String,
    pub pubsub_type: PubSubType,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct VpnConfig {
    pub api_url: String,
    pub api_key: String,
    pub network_id: String,
    pub retry_attempts: usize,
    pub retry_delay_ms: u64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ChannelsConfig {
    pub server_channel: String,
    pub command_channel: String,
    pub vpn_channel: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct StoreConfig {
    pub path: String,
    pub store_type: StoreType,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Config {
    pub pubsub: PubSubConfig,
    pub vpn: VpnConfig,
    pub channels: ChannelsConfig,
    pub store: StoreConfig,
}

impl Config {
    pub fn load() -> Result<Self, Box<dyn std::error::Error>> {
        // Try to load from config file first
        let home_dir = match dirs::home_dir() {
            Some(path) => path,
            None => {
                error!("Could not determine user's home directory");
                return Self::load_from_env();
            }
        };

        let config_path = home_dir.join(".config/vpn_handler/config.json");

        if config_path.exists() {
            match Self::load_from_file(&config_path) {
                Ok(config) => {
                    info!("Configuration loaded from file: {:?}", config_path);
                    return Ok(config);
                }
                Err(e) => {
                    error!("Error loading configuration from file: {}", e);
                    // If it fails, try to load from environment
                }
            }
        }

        // If not found, load from environment
        info!("Configuration file not found, loading from environment");
        Self::load_from_env()
    }

    pub fn load_from_env() -> Result<Self, Box<dyn std::error::Error>> {
        // Load variables from .env file if it exists
        dotenv().ok();

        // Helper to get environment variables
        fn get_env_var(name: &str) -> Result<String, Box<dyn std::error::Error>> {
            env::var(name).map_err(|_| {
                let message = format!("Environment variable not found: {}", name);
                error!("{}", message);
                message.into()
            })
        }

        // PubSub configuration
        let pubsub = PubSubConfig {
            url: get_env_var("PUBSUB_URL")?,
            pubsub_type: PubSubType::from_str(&get_env_var("PUBSUB_TYPE")?)
                .map_err(|_| "Unsupported PubSub type")?,
        };

        // VPN configuration
        let vpn = VpnConfig {
            api_url: get_env_var("VPN_API_URL")?,
            api_key: get_env_var("VPN_API_KEY")?,
            network_id: get_env_var("VPN_NETWORK_ID")?,
            retry_attempts: get_env_var("VPN_RETRY_ATTEMPTS")?.parse()?,
            retry_delay_ms: get_env_var("VPN_RETRY_DELAY_MS")?.parse()?,
        };

        // Channels configuration
        let channels = ChannelsConfig {
            server_channel: get_env_var("SERVER_CHANNEL")?,
            command_channel: get_env_var("COMMAND_CHANNEL")?,
            vpn_channel: get_env_var("VPN_CHANNEL")?,
        };

        // Store configuration
        let store = StoreConfig {
            path: get_env_var("DB_PATH")?,
            store_type: StoreType::from_str(&get_env_var("STORE_TYPE")?)
                .map_err(|_| "Unsupported store type")?,
        };

        Ok(Config {
            pubsub,
            vpn,
            channels,
            store,
        })
    }

    pub fn load_from_file<P: AsRef<Path>>(path: P) -> Result<Self, Box<dyn std::error::Error>> {
        let config_content = fs::read_to_string(path)?;
        let config: Config = serde_json::from_str(&config_content)?;
        Ok(config)
    }

    pub fn save_to_file<P: AsRef<Path>>(&self, path: P) -> Result<(), Box<dyn std::error::Error>> {
        // Ensure the directory exists
        if let Some(parent) = path.as_ref().parent() {
            fs::create_dir_all(parent)?;
        }

        let config_json = serde_json::to_string_pretty(self)?;
        fs::write(path, config_json)?;
        Ok(())
    }

    pub fn retry_config(&self) -> RetryConfig {
        RetryConfig {
            attempts: self.vpn.retry_attempts,
            delay: Duration::from_millis(self.vpn.retry_delay_ms),
        }
    }
}

#[derive(Debug, Clone)]
pub struct RetryConfig {
    pub attempts: usize,
    pub delay: Duration,
}
