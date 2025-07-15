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
use log::{error, info};
use std::error;

#[actix_web::main]
async fn main() -> Result<(), Box<dyn error::Error>> {
    // Initialize logger using configuration file.
    log4rs::init_file("./etc/log4rs.yaml", Default::default()).map_err(|e| {
        error!("Error initializing logger: {}", e);
        e
    })?;

    info!("Starting the application.");

    info!("Loading configuration...");
    let context = AppContext::instance();

    // Teste do storage singleton
    info!("Testando storage singleton...");
    context.set_storage("teste".to_string(), "valor_teste".to_string())?;

    if let Some(valor) = context.get_storage("teste")? {
        info!("✅ Storage funcionando! Valor recuperado: {}", valor);
    } else {
        error!("❌ Storage não funcionou!");
    }

    let (actix_server, server_handle) = server::create_server()?;

    rt::spawn(actix_server);

    tokio::signal::ctrl_c().await?;
    server_handle.stop(true).await;
    Ok(())
}
