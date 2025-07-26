mod api;
mod app_context;
mod config;
mod errors;
mod models;
mod persistence;
mod pub_sub;
mod queue_handler;
mod vpn;
use crate::{api::server, app_context::AppContext, queue_handler::Handle};
use actix_web::rt;
use log::{debug, error, info};

#[actix_web::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Initialize logger using configuration file.
    log4rs::init_file("./etc/log4rs.yaml", Default::default()).map_err(|e| {
        error!("Error initializing logger: {}", e);
        e
    })?;

    info!("Starting the application.");

    info!("Loading application context...");
    let _context = AppContext::instance();

    info!("Application context load sucefully!");
    debug!("{:?}", AppContext::instance().config());

    let (actix_server, server_handle) = server::create_server()?;

    info!("Starting queue handler...");
    let queue_handle = Handle::new()?;

    // Spawn both the HTTP server and queue handler
    rt::spawn(actix_server);
    rt::spawn(async move {
        if let Err(e) = queue_handle.run().await {
            error!("Queue handler error: {}", e);
        }
    });

    tokio::signal::ctrl_c().await?;
    server_handle.stop(true).await;
    Ok(())
}
