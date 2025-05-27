use crate::pub_sub::factories::PubSubType;
use crate::persistence::factories::StoreType;
use crate::vpn::factories::VpnType;
use dotenvy::dotenv;
use log::error;
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
    pub pub_sub_config: PubSubConfig,
    pub vpn_config: VpnConfig,
    pub db_config: DbConfig,
    pub admine_channels_map: AdmineChannelsMap,
    pub retry_config: RetryConfig,
}

impl Config {
    /// Carrega todas as configurações das variáveis de ambiente
    pub fn load() -> Result<Self, Box<dyn std::error::Error>> {
        dotenv().ok();

        fn fetch_env_var(var_name: &str) -> Result<String, Box<dyn std::error::Error>> {
            env::var(var_name).map_err(|_| {
                let msg = format!("Missing environment variable: {}", var_name);
                error!("{}", msg);
                msg.into()
            })
        }

        // Carrega todas as variáveis de ambiente
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

        // Parse dos tipos de enum
        let pub_sub_type = PubSubType::from_str(&pubsub_type).map_err(|_| {
            error!("Unsupported PubSub type: {}", pubsub_type);
            "Unsupported PubSub type"
        })?;

        let store_type_enum = StoreType::from_str(&store_type).map_err(|_| {
            error!("Unsupported Store type: {}", store_type);
            "Unsupported Store type"
        })?;

        // Criar as estruturas de configuração
        let pub_sub_config = PubSubConfig {
            url: pubsub_url,
            pub_sub_type,
        };

        let vpn_config = VpnConfig {
            api_url,
            api_key,
            network_id,
            vpn_type: VpnType::Zerotier, // Por enquanto fixo como Zerotier
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

        Ok(Config {
            pub_sub_config,
            vpn_config,
            db_config,
            admine_channels_map,
            retry_config,
        })
    }
}