mod api;
mod config;
mod errors;
mod models;
mod persistence;
mod pub_sub;
mod vpn;
use crate::{
    api::server,
    config::Config,
    persistence::storage_manager::StorageManager,
    vpn::vpn_factory::VpnFactory,
};
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
    let config = Config::instance();

    StorageManager::instance();

    // let (actix_server, server_handle) = server::create_server()?;

    // rt::spawn(actix_server);

    let vpn = VpnFactory::create_vpn(
        config.vpn_config().vpn_type().clone(),
        config.vpn_config().api_url().to_string(),
        config.vpn_config().api_key().to_string(),
        config.vpn_config().network_id().to_string(),
    ).unwrap();

    let a = vpn.get_member_ips_in_vpn(String::from("a41a6f919c")).await.unwrap();

    print!("{:?}", a);

    // tokio::signal::ctrl_c().await?;
    // server_handle.stop(true).await;
    Ok(())
}
