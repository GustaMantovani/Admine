use crate::persistence::key_value_storage::DynKeyValueStore;
use crate::vpn::vpn::DynVpn;
use crate::{api::services, config::Config, models::error_response::ErrorResponse};
use actix_web::{
    dev::Server, dev::ServerHandle, error::JsonPayloadError, web, App, HttpRequest, HttpResponse,
    HttpServer, Result,
};
use log::error;
use std::sync::Arc;

/// Custom error handler for JSON parsing errors
fn json_error_handler(err: JsonPayloadError, _req: &HttpRequest) -> actix_web::Error {
    let error_message = match &err {
        JsonPayloadError::ContentType => "Content type must be application/json".to_string(),
        JsonPayloadError::Deserialize(json_err) => {
            format!("Invalid JSON format: {}", json_err)
        }
        JsonPayloadError::Overflow { .. } => "Request payload too large".to_string(),
        JsonPayloadError::Payload(payload_err) => {
            format!("Payload error: {}", payload_err)
        }
        _ => "Invalid JSON request".to_string(),
    };

    error!("JSON parsing error: {}", error_message);

    let response = HttpResponse::BadRequest().json(ErrorResponse {
        message: error_message,
    });

    actix_web::error::InternalError::from_response(err, response).into()
}

pub fn create_server(
    vpn_client: Arc<DynVpn>,
    storage: Arc<DynKeyValueStore>,
    config: Arc<Config>,
) -> Result<(Server, ServerHandle), std::io::Error> {
    let host = config.api_config().host().clone();
    let port = *config.api_config().port();

    let server = HttpServer::new(move || {
        App::new()
            .app_data(
                web::JsonConfig::default()
                    .limit(4096) // 4KB limit
                    .error_handler(json_error_handler),
            )
            .app_data(web::Data::from(Arc::clone(&vpn_client)))
            .app_data(web::Data::from(Arc::clone(&storage)))
            .app_data(web::Data::from(Arc::clone(&config)))
            .service(services::server_ip)
            .service(services::auth_member)
            .service(services::vpn_id)
    })
    .bind((host.as_str(), port))?
    .workers(1)
    .disable_signals()
    .run();

    let handle = server.handle();
    Ok((server, handle))
}
