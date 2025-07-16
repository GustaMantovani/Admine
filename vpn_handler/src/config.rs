use crate::persistence::key_value_storage_factory::StoreType;
use crate::pub_sub::pub_sub_factory::PubSubType;
use crate::vpn::vpn_factory::VpnType;
use anyhow::Result;
use derive_builder::Builder;
use dotenvy::dotenv;
use log::info;
use serde::Deserialize;
use serde_with::{serde_as, DurationMilliSeconds};
use std::fmt;
use std::time::Duration;

#[serde_as]
#[derive(Debug, Clone, Deserialize, Builder)]
#[builder(setter(into))]
pub struct ApiConfig {
    #[serde(rename = "API_HOST")]
    pub host: String,
    #[serde(rename = "API_PORT")]
    pub port: u16,
}

#[derive(Debug, Clone, Deserialize, Builder)]
#[builder(setter(into))]
pub struct PubSubConfig {
    #[serde(rename = "PUBSUB_URL")]
    pub url: String,
    #[serde(rename = "PUBSUB_TYPE")]
    pub pub_sub_type: PubSubType,
}

#[derive(Debug, Clone, Deserialize, Builder)]
#[builder(setter(into))]
pub struct VpnConfig {
    #[serde(rename = "VPN_API_URL")]
    pub api_url: String,
    #[serde(rename = "VPN_API_KEY")]
    pub api_key: String,
    #[serde(rename = "VPN_NETWORK_ID")]
    pub network_id: String,
    #[serde(rename = "VPN_TYPE")]
    pub vpn_type: VpnType,
}

#[derive(Debug, Clone, Deserialize, Builder)]
#[builder(setter(into))]
pub struct DbConfig {
    #[serde(rename = "DB_PATH")]
    pub path: String,
    #[serde(rename = "STORE_TYPE")]
    pub store_type: StoreType,
}

#[derive(Debug, Clone, Deserialize, Builder)]
#[builder(setter(into))]
pub struct AdmineChannelsMap {
    #[serde(rename = "SERVER_CHANNEL")]
    pub server_channel: String,
    #[serde(rename = "COMMAND_CHANNEL")]
    pub command_channel: String,
    #[serde(rename = "VPN_CHANNEL")]
    pub vpn_channel: String,
}

impl fmt::Display for AdmineChannelsMap {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(
            f,
            "AdmineChannelsMap {{ server_channel: {}, command_channel: {}, vpn_channel: {} }}",
            self.server_channel, self.command_channel, self.vpn_channel
        )
    }
}

#[serde_as]
#[derive(Debug, Clone, Deserialize, Builder)]
#[builder(setter(into))]
pub struct RetryConfig {
    #[serde(rename = "VPN_RETRY_ATTEMPTS")]
    pub attempts: usize,
    #[serde(rename = "VPN_RETRY_DELAY_MS")]
    #[serde_as(as = "DurationMilliSeconds<u64>")]
    pub delay: Duration,
}

#[derive(Debug, Clone, Deserialize, Builder)]
#[builder(setter(into))]
pub struct Config {
    #[serde(rename = "SELF_ORIGIN_NAME")]
    pub self_origin_name: String,
    #[serde(flatten)]
    pub api_config: ApiConfig,
    #[serde(flatten)]
    pub pub_sub_config: PubSubConfig,
    #[serde(flatten)]
    pub vpn_config: VpnConfig,
    #[serde(flatten)]
    pub db_config: DbConfig,
    #[serde(flatten)]
    pub admine_channels_map: AdmineChannelsMap,
    #[serde(flatten)]
    pub retry_config: RetryConfig,
}

impl Config {
    /// Loads all configuration from environment variables
    pub fn new() -> Result<Self> {
        dotenv().ok();

        let config: Config = envy::from_env()?;

        info!("Configuration loaded successfully: {:#?}", config);

        Ok(config)
    }
}
