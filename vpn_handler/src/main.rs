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
    vpn::{vpn_factory::VpnFactory, zerotier_vpn::ZerotierVpn},
};
use actix_web::rt;
use log::{error, info};
use zerotier_central_api::apis::configuration::Configuration;
use std::error;

#[tokio::main]
async fn main() -> Result<(), Box<dyn error::Error>> {
    // Initialize logger using configuration file.
    log4rs::init_file("./etc/log4rs.yaml", Default::default()).map_err(|e| {
        error!("Error initializing logger: {}", e);
        e
    })?;

    info!("Starting the application.");

    info!("Loading configuration...");
    Config::instance();

    // let (actix_server, server_handle) = server::create_server()?;

    // rt::spawn(actix_server);

    // tokio::signal::ctrl_c().await?;
    // server_handle.stop(true).await;

    // let vpn = VpnFactory::create_vpn(
    //     Config::instance().vpn_config().vpn_type().clone(),
    //     Config::instance().vpn_config().api_url().to_string(),
    //     Config::instance().vpn_config().api_key().to_string(),
    //     Config::instance().vpn_config().network_id().to_string(),
    // )
    // .unwrap();

    let mut config = Configuration::new();
    config.base_path = Config::instance().vpn_config().api_url().to_string();
    config.api_key = Some(zerotier_central_api::apis::configuration::ApiKey {
        prefix: None,
        key: Config::instance().vpn_config().api_key().to_string(),
    });
    config.bearer_access_token = Some(Config::instance().vpn_config().api_key().to_string());

    print!("{:?}", config);

    let vpn = ZerotierVpn::new(
        config,
        Config::instance().vpn_config().network_id().to_string(),
    );

    // vpn.auth_member(String::from("a41a6f919c"), None).await?;

    vpn.get_all().await;

    Ok(())
}
