mod api;
mod config;
mod errors;
mod models;
mod persistence;
mod pub_sub;
mod queue_handler;
mod vpn;
use crate::{
    api::server, config::Config, persistence::key_value_storage_factory::StoreFactory,
    pub_sub::pub_sub_factory::PubSubFactory, queue_handler::Handle, vpn::vpn_factory::VpnFactory,
};
use actix_web::rt;
use log::info;
use log::LevelFilter;
use log4rs::append::console::ConsoleAppender;
use log4rs::append::file::FileAppender;
use log4rs::config::{Appender, Root};
use log4rs::encode::pattern::PatternEncoder;
use std::sync::Arc;

#[actix_web::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    init_logger()?;

    info!("Starting the application.");

    let config = Arc::new(Config::new()?);

    let storage = Arc::new(StoreFactory::create_store_instance(
        config.db_config().store_type().clone(),
        config.db_config().path(),
    )?);

    let vpn_client = Arc::new(
        VpnFactory::create_vpn(
            config.vpn_config().vpn_type().clone(),
            config.vpn_config().api_url().clone(),
            config.vpn_config().api_key().clone(),
            config.vpn_config().network_id().clone(),
        )
        .map_err(|e| format!("Failed to create VPN client: {:?}", e))?,
    );

    let mut pub_sub = PubSubFactory::create_pubsub_instance(
        config.pub_sub_config().pub_sub_type().clone(),
        config.pub_sub_config().url(),
    )
    .await
    .map_err(|e| format!("Failed to create PubSub client: {:?}", e))?;

    pub_sub
        .subscribe(vec![
            config.admine_channels_map().server_channel().clone(),
            config.admine_channels_map().command_channel().clone(),
        ])
        .map_err(|e| format!("Failed to subscribe to channels: {:?}", e))?;

    let (actix_server, server_handle) = server::create_server(
        Arc::clone(&vpn_client),
        Arc::clone(&storage),
        Arc::clone(&config),
    )?;

    let (shutdown_tx, shutdown_rx) = tokio::sync::watch::channel(false);

    info!("Starting queue handler...");
    let queue_handle = Handle::new(pub_sub, vpn_client, storage, config, shutdown_rx);
    let queue_task = rt::spawn(queue_handle.run());

    let shutdown_task = rt::spawn(async move {
        shutdown_signal().await;
        info!("Initiating graceful shutdown...");
        server_handle.stop(true).await;
        let _ = shutdown_tx.send(true);
    });

    actix_server.await?;

    let _ = queue_task.await;
    let _ = shutdown_task.await;

    info!("Application shutdown complete.");
    Ok(())
}

async fn shutdown_signal() {
    use tokio::signal;

    #[cfg(unix)]
    {
        use signal::unix::{signal, SignalKind};
        let mut sigterm =
            signal(SignalKind::terminate()).expect("Failed to register SIGTERM handler");
        tokio::select! {
            _ = signal::ctrl_c() => { info!("SIGINT received."); }
            _ = sigterm.recv() => { info!("SIGTERM received."); }
        }
    }
}

fn init_logger() -> Result<(), Box<dyn std::error::Error>> {
    if let Err(err) = log4rs::init_file("./etc/log4rs.yaml", Default::default()) {
        eprintln!("Error initializing logger from file: {}", err);

        const PATTERN_ENCONDER: &str = "{d} - {l} - {m}{n}";

        let stdout = ConsoleAppender::builder()
            .encoder(Box::new(PatternEncoder::new(PATTERN_ENCONDER)))
            .build();

        let file = FileAppender::builder()
            .encoder(Box::new(PatternEncoder::new(PATTERN_ENCONDER)))
            .build("/tmp/admine/logs/vpn_handler.log")?;

        let config = log4rs::config::Config::builder()
            .appender(Appender::builder().build("stdout", Box::new(stdout)))
            .appender(Appender::builder().build("file", Box::new(file)))
            .build(
                Root::builder()
                    .appender("stdout")
                    .appender("file")
                    .build(LevelFilter::Info),
            )?;

        log4rs::init_config(config)?;
    }

    Ok(())
}
