use dotenvy::dotenv;
use log::{error, info};
use std::env;
use std::str::FromStr;
use std::time::Duration;
use serde::{Deserialize, Serialize};
use std::fs;
use std::path::Path;

use crate::pub_sub::factories::PubSubType;
use crate::persistence::factories::StoreType;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PubSubConfig {
    pub url: String,
    pub tipo: PubSubType,
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
    pub tipo: StoreType,
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
        // Tenta carregar do arquivo de configuração primeiro
        let home_dir = match dirs::home_dir() {
            Some(path) => path,
            None => {
                error!("Não foi possível determinar o diretório home do usuário");
                return Self::load_from_env();
            }
        };
        
        let config_path = home_dir.join(".config/vpn_handler/config.json");
        
        if config_path.exists() {
            match Self::load_from_file(&config_path) {
                Ok(config) => {
                    info!("Configuração carregada do arquivo: {:?}", config_path);
                    return Ok(config);
                }
                Err(e) => {
                    error!("Erro ao carregar configuração do arquivo: {}", e);
                    // Se falhar, tenta carregar do ambiente
                }
            }
        }
        
        // Se não conseguir carregar do arquivo, carrega do ambiente
        info!("Arquivo de configuração não encontrado, carregando do ambiente");
        Self::load_from_env()
    }

    pub fn load_from_env() -> Result<Self, Box<dyn std::error::Error>> {
        // Carrega variáveis do arquivo .env, se existir
        dotenv().ok();

        // Helper para obter variáveis de ambiente
        fn obter_var_env(nome: &str) -> Result<String, Box<dyn std::error::Error>> {
            env::var(nome).map_err(|_| {
                let mensagem = format!("Variável de ambiente não encontrada: {}", nome);
                error!("{}", mensagem);
                mensagem.into()
            })
        }

        // Configuração do PubSub
        let pubsub = PubSubConfig {
            url: obter_var_env("PUBSUB_URL")?,
            tipo: PubSubType::from_str(&obter_var_env("PUBSUB_TYPE")?)
                .map_err(|_| "Tipo de PubSub não suportado")?,
        };

        // Configuração da VPN
        let vpn = VpnConfig {
            api_url: obter_var_env("VPN_API_URL")?,
            api_key: obter_var_env("VPN_API_KEY")?,
            network_id: obter_var_env("VPN_NETWORK_ID")?,
            retry_attempts: obter_var_env("VPN_RETRY_ATTEMPTS")?.parse()?,
            retry_delay_ms: obter_var_env("VPN_RETRY_DELAY_MS")?.parse()?,
        };

        // Configuração dos canais
        let channels = ChannelsConfig {
            server_channel: obter_var_env("SERVER_CHANNEL")?,
            command_channel: obter_var_env("COMMAND_CHANNEL")?,
            vpn_channel: obter_var_env("VPN_CHANNEL")?,
        };

        // Configuração do armazenamento
        let store = StoreConfig {
            path: obter_var_env("DB_PATH")?,
            tipo: StoreType::from_str(&obter_var_env("STORE_TYPE")?)
                .map_err(|_| "Tipo de armazenamento não suportado")?,
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
        // Garantir que o diretório exista
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