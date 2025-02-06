mod errors;
mod handle;
mod models;
mod persistence;
mod pub_sub;
mod vpn;
mod zerotier;

use handle::Handle;
use log::{error, info};
use log4rs;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Initialize logger using configuration file.
    log4rs::init_file("./config/log4rs.yaml", Default::default()).map_err(|e| {
        error!("Error initializing logger: {}", e);
        e
    })?;

    info!("Starting the application.");

    // Create and run the handle.
    let handle = Handle::new()?;
    info!("Handle created successfully.");
    handle.run().await?;

    Ok(())
}
