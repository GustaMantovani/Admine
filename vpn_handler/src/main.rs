mod config;
mod errors;
mod handle;
mod models;
mod persistence;
mod pub_sub;
mod vpn;

use config::Config;
use handle::Handle;
use log::{error, info};
use log4rs;
use std::env;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Initialize logger using configuration file.
    log4rs::init_file("./config/log4rs.yaml", Default::default()).map_err(|e| {
        error!("Error initializing logger: {}", e);
        e
    })?;

    info!("Starting application.");

    // Check if configuration export was requested
    let args: Vec<String> = env::args().collect();
    if args.len() > 1 && args[1] == "--export-config" {
        let config = Config::load_from_env()?;

        let home_dir =
            dirs::home_dir().ok_or("Could not determine user's home directory")?;
        let config_path = home_dir.join(".config/vpn_handler/config.json");

        config.save_to_file(&config_path)?;
        info!("Configuration exported to: {:?}", config_path);
        return Ok(());
    }

    // Create and run the handle.
    let handle = Handle::new()?;
    info!("Handle created successfully.");
    handle.run().await?;

    Ok(())
}
