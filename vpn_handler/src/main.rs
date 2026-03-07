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
use log::LevelFilter;
use log::{debug, info};
use log4rs::append::console::ConsoleAppender;
use log4rs::append::file::FileAppender;
use log4rs::config::{Appender, Config, Root};
use log4rs::encode::pattern::PatternEncoder;

#[actix_web::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    init_logger()?;

    info!("Starting the application.");

    info!("Loading application context...");
    let _context = AppContext::instance();

    info!("Application context load sucefully!");
    debug!("{:?}", AppContext::instance().config());

    let (actix_server, _) = server::create_server()?;

    info!("Starting queue handler...");
    let queue_handle = Handle::new()?;

    rt::spawn(queue_handle.run());
    let _ = actix_server.await;

    Ok(())
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
            .build("./etc/zerotier_handler.log")?;

        let config = Config::builder()
            .appender(Appender::builder().build("stdout", Box::new(stdout)))
            .appender(Appender::builder().build("file", Box::new(file)))
            .build(Root::builder().appender("stdout").build(LevelFilter::Info))?;

        log4rs::init_config(config)?;
    }

    Ok(())
}
