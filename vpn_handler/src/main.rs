mod api;
mod app_context;
mod config;
mod errors;
mod models;
mod persistence;
mod pub_sub;
mod vpn;
use crate::{api::server, app_context::AppContext};
use actix_web::rt;
use anyhow::Result;
use log::{error, info};

#[actix_web::main]
async fn main() -> Result<()> {
    // Initialize logger using configuration file.
    log4rs::init_file("./etc/log4rs.yaml", Default::default()).map_err(|e| {
        error!("Error initializing logger: {}", e);
        e
    })?;

    info!("Starting the application.");

    info!("Loading configuration...");
    let _context = AppContext::instance();

    let (actix_server, server_handle) = server::create_server()?;

    rt::spawn(actix_server);

    tokio::signal::ctrl_c().await?;
    server_handle.stop(true).await;
    Ok(())
}
