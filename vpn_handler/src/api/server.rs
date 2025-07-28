use crate::{api::services, app_context::AppContext};
use actix_web::{dev::Server, dev::ServerHandle, App, HttpServer};

pub fn create_server() -> Result<(Server, ServerHandle), std::io::Error> {
    let config = AppContext::instance().config();
    let host = &config.api_config().host();
    let port = *config.api_config().port();

    let server = HttpServer::new(|| {
        App::new()
            .service(services::server_ip)
            .service(services::auth_member)
            .service(services::vpn_id)
    })
    .bind((host.as_str(), port))?
    .workers(1)
    .run();

    let handle = server.handle();
    Ok((server, handle))
}
