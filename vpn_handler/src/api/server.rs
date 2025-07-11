use crate::{api::services, config::Config};
use actix_web::{dev::Server, dev::ServerHandle, App, HttpServer};

pub fn create_server() -> Result<(Server, ServerHandle), std::io::Error> {
    let config = Config::instance();
    let host = config.api_config().host();
    let port = *config.api_config().port();

    let server = HttpServer::new(|| App::new().service(services::hello))
        .bind((host, port))?
        .workers(1)
        .run();

    let handle = server.handle();
    Ok((server, handle))
}
