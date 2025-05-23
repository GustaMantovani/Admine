mod errors;
mod handle;
mod models;
mod persistence;
mod pub_sub;
mod vpn;
mod config;

use handle::Handle;
use log::{error, info};
use log4rs;
use std::env;
use std::path::PathBuf;
use config::Config;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Initialize logger using configuration file.
    log4rs::init_file("./config/log4rs.yaml", Default::default()).map_err(|e| {
        error!("Error initializing logger: {}", e);
        e
    })?;

    info!("Iniciando aplicação.");
    
    // Verificar se foi solicitada exportação da configuração
    let args: Vec<String> = env::args().collect();
    if args.len() > 1 && args[1] == "--export-config" {
        let config = Config::load_from_env()?;
        
        let home_dir = dirs::home_dir()
            .ok_or("Não foi possível determinar o diretório home do usuário")?;
        let config_path = home_dir.join(".config/vpn_handler/config.json");
        
        config.save_to_file(&config_path)?;
        info!("Configuração exportada para: {:?}", config_path);
        return Ok(());
    }

    // Create and run the handle.
    let handle = Handle::new()?;
    info!("Handle criado com sucesso.");
    handle.run().await?;

    Ok(())
}
