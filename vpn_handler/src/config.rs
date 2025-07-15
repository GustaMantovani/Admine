use crate::persistence::key_value_storage_factory::StoreType;
use crate::pub_sub::pub_sub_factory::PubSubType;
use crate::vpn::vpn_factory::VpnType;
use dotenvy::dotenv;
use log::error;
use std::env;
use std::fmt;
use std::str::FromStr;
use std::time::Duration;

#[derive(Debug, Clone)]
pub struct ApiConfig {
    host: String,
    port: u16,
}

impl ApiConfig {
    pub fn new(host: String, port: u16) -> Self {
        Self { host, port }
    }

    pub fn host(&self) -> &str {
        &self.host
    }

    pub fn port(&self) -> &u16 {
        &self.port
    }
}

#[derive(Debug, Clone)]
pub struct PubSubConfig {
    url: String,
    pub_sub_type: PubSubType,
}

impl PubSubConfig {
    pub fn new(url: String, pub_sub_type: PubSubType) -> Self {
        Self { url, pub_sub_type }
    }

    pub fn url(&self) -> &str {
        &self.url
    }

    pub fn pub_sub_type(&self) -> &PubSubType {
        &self.pub_sub_type
    }
}

#[derive(Debug, Clone)]
pub struct VpnConfig {
    api_url: String,
    api_key: String,
    network_id: String,
    vpn_type: VpnType,
}

impl VpnConfig {
    pub fn new(api_url: String, api_key: String, network_id: String, vpn_type: VpnType) -> Self {
        Self {
            api_url,
            api_key,
            network_id,
            vpn_type,
        }
    }

    pub fn api_url(&self) -> &str {
        &self.api_url
    }

    pub fn api_key(&self) -> &str {
        &self.api_key
    }

    pub fn network_id(&self) -> &str {
        &self.network_id
    }

    pub fn vpn_type(&self) -> &VpnType {
        &self.vpn_type
    }
}

#[derive(Debug, Clone)]
pub struct DbConfig {
    path: String,
    store_type: StoreType,
}

impl DbConfig {
    pub fn new(path: String, store_type: StoreType) -> Self {
        Self { path, store_type }
    }

    pub fn path(&self) -> &str {
        &self.path
    }

    pub fn store_type(&self) -> &StoreType {
        &self.store_type
    }
}

#[derive(Debug, Clone)]
pub struct AdmineChannelsMap {
    server_channel: String,
    command_channel: String,
    vpn_channel: String,
}

impl AdmineChannelsMap {
    pub fn new(server_channel: String, command_channel: String, vpn_channel: String) -> Self {
        Self {
            server_channel,
            command_channel,
            vpn_channel,
        }
    }

    pub fn server_channel(&self) -> &str {
        &self.server_channel
    }

    pub fn command_channel(&self) -> &str {
        &self.command_channel
    }

    pub fn vpn_channel(&self) -> &str {
        &self.vpn_channel
    }
}

impl fmt::Display for AdmineChannelsMap {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(
            f,
            "AdmineChannelsMap {{ server_channel: {}, command_channel: {}, vpn_channel: {} }}",
            self.server_channel(),
            self.command_channel(),
            self.vpn_channel()
        )
    }
}

#[derive(Debug, Clone)]
pub struct RetryConfig {
    attempts: usize,
    delay: Duration,
}

impl RetryConfig {
    pub fn new(attempts: usize, delay: Duration) -> Self {
        Self { attempts, delay }
    }

    pub fn attempts(&self) -> usize {
        self.attempts
    }

    pub fn delay(&self) -> Duration {
        self.delay
    }
}

#[derive(Debug, Clone)]
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
    pub fn self_origin_name(&self) -> &str {
        &self.self_origin_name
    }

    pub fn api_config(&self) -> &ApiConfig {
        &self.api_config
    }

    pub fn pub_sub_config(&self) -> &PubSubConfig {
        &self.pub_sub_config
    }

    pub fn vpn_config(&self) -> &VpnConfig {
        &self.vpn_config
    }

    pub fn db_config(&self) -> &DbConfig {
        &self.db_config
    }

    pub fn admine_channels_map(&self) -> &AdmineChannelsMap {
        &self.admine_channels_map
    }

    pub fn retry_config(&self) -> &RetryConfig {
        &self.retry_config
    }
}

impl fmt::Display for Config {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(
            f,
            "Config {{\n\
            \x20\x20self_origin_name: {},\n\
            \x20\x20pub_sub_config: PubSubConfig {{ url: {}, type: {:?} }},\n\
            \x20\x20vpn_config: VpnConfig {{ api_url: {}, network_id: {}, type: {:?} }},\n\
            \x20\x20db_config: DbConfig {{ path: {}, type: {:?} }},\n\
            \x20\x20admine_channels_map: {},\n\
            \x20\x20retry_config: RetryConfig {{ attempts: {}, delay: {:?} }}\n\
            }}",
            self.self_origin_name,
            self.pub_sub_config.url(),
            self.pub_sub_config.pub_sub_type(),
            self.vpn_config.api_url(),
            self.vpn_config.network_id(),
            self.vpn_config.vpn_type(),
            self.db_config.path(),
            self.db_config.store_type(),
            self.admine_channels_map,
            self.retry_config.attempts(),
            self.retry_config.delay()
        )
    }
}

impl Config {
    /// Loads all configuration from environment variables
    pub fn new() -> Result<Self, String> {
        dotenv().ok();

        fn fetch_env_var(var_name: &str) -> Result<String, String> {
            env::var(var_name).map_err(|_| {
                let msg = format!("Missing environment variable: {}", var_name);
                msg
            })
        }

        // Load all environment variables
        let self_origin_name = fetch_env_var("SELF_ORIGIN_NAME")?;
        let api_host = fetch_env_var("API_HOST")?;
        let api_port = fetch_env_var("API_PORT")?;
        let pubsub_url = fetch_env_var("PUBSUB_URL")?;
        let pubsub_type = fetch_env_var("PUBSUB_TYPE")?;
        let api_url = fetch_env_var("VPN_API_URL")?;
        let api_key = fetch_env_var("VPN_API_KEY")?;
        let network_id = fetch_env_var("VPN_NETWORK_ID")?;
        let server_channel = fetch_env_var("SERVER_CHANNEL")?;
        let command_channel = fetch_env_var("COMMAND_CHANNEL")?;
        let vpn_channel = fetch_env_var("VPN_CHANNEL")?;
        let vpn_type = fetch_env_var("VPN_TYPE")?;
        let db_path = fetch_env_var("DB_PATH")?;
        let store_type = fetch_env_var("STORE_TYPE")?;
        let retry_attempts = fetch_env_var("VPN_RETRY_ATTEMPTS")?;
        let retry_delay_ms = fetch_env_var("VPN_RETRY_DELAY_MS")?;

        let api_config = ApiConfig::new(
            api_host,
            u16::from_str(&api_port).map_err(|e| e.to_string())?,
        );

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
        let pub_sub_config = PubSubConfig::new(pubsub_url, pub_sub_type);

        let vpn_config = VpnConfig::new(api_url, api_key, network_id, VpnType::PublicIp);

        let db_config = DbConfig::new(db_path, store_type_enum);

        let admine_channels_map =
            AdmineChannelsMap::new(server_channel, command_channel, vpn_channel);

        let retry_config = RetryConfig::new(
            retry_attempts
                .parse()
                .map_err(|e: std::num::ParseIntError| e.to_string())?,
            Duration::from_millis(
                retry_delay_ms
                    .parse()
                    .map_err(|e: std::num::ParseIntError| e.to_string())?,
            ),
        );

        let final_config = Config {
            self_origin_name,
            api_config,
            pub_sub_config,
            vpn_config,
            db_config,
            admine_channels_map,
            retry_config,
        };

        // info!("Configuration loaded successfully: {}", final_config);

        Ok(final_config)
    }
}
