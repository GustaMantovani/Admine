use crate::persistence::key_value_storage_factory::StoreType;
use crate::pub_sub::pub_sub_factory::PubSubType;
use crate::vpn::vpn_factory::VpnType;
use dotenvy::dotenv;
use log::error;
use log::info;
use std::env;
use std::fmt;
use std::str::FromStr;
use std::time::Duration;

#[derive(Debug, Clone)]
pub struct PubSubConfig {
    pub url: String,
    pub pub_sub_type: PubSubType,
}

#[derive(Debug, Clone)]
pub struct VpnConfig {
    pub api_url: String,
    pub api_key: String,
    pub network_id: String,
    pub vpn_type: VpnType,
}

#[derive(Debug, Clone)]
pub struct DbConfig {
    pub path: String,
    pub store_type: StoreType,
}

#[derive(Debug, Clone)]
pub struct AdmineChannelsMap {
    pub server_channel: String,
    pub command_channel: String,
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

#[derive(Debug, Clone)]
pub struct RetryConfig {
    pub attempts: usize,
    pub delay: Duration,
}

#[derive(Debug, Clone)]
pub struct Config {
    pub self_origin_name: String,
    pub pub_sub_config: PubSubConfig,
    pub vpn_config: VpnConfig,
    pub db_config: DbConfig,
    pub admine_channels_map: AdmineChannelsMap,
    pub retry_config: RetryConfig,
}

impl fmt::Display for Config {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(
            f,
            "Config {{\n\
            \x20\x20pub_sub_config: PubSubConfig {{ url: {}, type: {:?} }},\n\
            \x20\x20vpn_config: VpnConfig {{ api_url: {}, network_id: {}, type: {:?} }},\n\
            \x20\x20db_config: DbConfig {{ path: {}, type: {:?} }},\n\
            \x20\x20admine_channels_map: {},\n\
            \x20\x20retry_config: RetryConfig {{ attempts: {}, delay: {:?} }}\n\
            }}",
            self.pub_sub_config.url,
            self.pub_sub_config.pub_sub_type,
            self.vpn_config.api_url,
            self.vpn_config.network_id,
            self.vpn_config.vpn_type,
            self.db_config.path,
            self.db_config.store_type,
            self.admine_channels_map,
            self.retry_config.attempts,
            self.retry_config.delay
        )
    }
}

impl Config {
    /// Loads all configuration from environment variables
    pub fn load() -> Result<Self, Box<dyn std::error::Error>> {
        dotenv().ok();

        fn fetch_env_var(var_name: &str) -> Result<String, Box<dyn std::error::Error>> {
            env::var(var_name).map_err(|_| {
                let msg = format!("Missing environment variable: {}", var_name);
                error!("{}", msg);
                msg.into()
            })
        }

        // Load all environment variables
        let self_origin_name = fetch_env_var("SELF_ORIGIN_NAME")?;
        let pubsub_url = fetch_env_var("PUBSUB_URL")?;
        let pubsub_type = fetch_env_var("PUBSUB_TYPE")?;
        let api_url = fetch_env_var("VPN_API_URL")?;
        let api_key = fetch_env_var("VPN_API_KEY")?;
        let network_id = fetch_env_var("VPN_NETWORK_ID")?;
        let server_channel = fetch_env_var("SERVER_CHANNEL")?;
        let command_channel = fetch_env_var("COMMAND_CHANNEL")?;
        let vpn_channel = fetch_env_var("VPN_CHANNEL")?;
        let db_path = fetch_env_var("DB_PATH")?;
        let store_type = fetch_env_var("STORE_TYPE")?;
        let retry_attempts = fetch_env_var("VPN_RETRY_ATTEMPTS")?;
        let retry_delay_ms = fetch_env_var("VPN_RETRY_DELAY_MS")?;

        // Parse enum types
        let pub_sub_type = PubSubType::from_str(&pubsub_type).map_err(|_| {
            error!("Unsupported PubSub type: {}", pubsub_type);
            "Unsupported PubSub type"
        })?;

        let store_type_enum = StoreType::from_str(&store_type).map_err(|_| {
            error!("Unsupported Store type: {}", store_type);
            "Unsupported Store type"
        })?;

        // Create configuration structures
        let pub_sub_config = PubSubConfig {
            url: pubsub_url,
            pub_sub_type,
        };

        let vpn_config = VpnConfig {
            api_url,
            api_key,
            network_id,
            vpn_type: VpnType::Zerotier, // Currently fixed as Zerotier
        };

        let db_config = DbConfig {
            path: db_path,
            store_type: store_type_enum,
        };

        let admine_channels_map = AdmineChannelsMap {
            server_channel,
            command_channel,
            vpn_channel,
        };

        let retry_config = RetryConfig {
            attempts: retry_attempts.parse()?,
            delay: Duration::from_millis(retry_delay_ms.parse()?),
        };

        let final_config = Config {
            self_origin_name,
            pub_sub_config,
            vpn_config,
            db_config,
            admine_channels_map,
            retry_config,
        };

        info!("Configuration loaded successfully: {}", final_config);

        Ok(final_config)
    }
}
