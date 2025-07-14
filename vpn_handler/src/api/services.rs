use crate::{
    config::Config,
    models::api_models::{AuthMemberRequest, ErrorResponse, ServerIpResponse, VpnIdResponse},
    persistence::key_value_storage::get_global,
    vpn::vpn_factory::VpnFactory,
};
use actix_web::{get, post, web, HttpResponse, Responder};
use log::{error, info};

#[get("/status")]
async fn status() -> impl Responder {
    HttpResponse::Ok().body("1")
}

#[get("/server-ip")]
pub async fn server_ip() -> impl Responder {
    match get_global("server_ip").unwrap_or(None) {
        Some(ip) => HttpResponse::Ok().json(ServerIpResponse { server_ip: ip }),
        None => HttpResponse::InternalServerError()
            .json(ErrorResponse { message: "error".to_string() }),
    }
}

#[post("/auth-member")]
pub async fn auth_member(member_data: web::Json<AuthMemberRequest>) -> impl Responder {
    info!("Authorizing member: {}", member_data.member_id);
    
    let config = Config::instance();
    let vpn_result = VpnFactory::create_vpn(
        config.vpn_config().vpn_type().clone(),
        config.vpn_config().api_url().to_string(),
        config.vpn_config().api_key().to_string(),
        config.vpn_config().network_id().to_string(),
    );

    let vpn = match vpn_result {
        Ok(v) => v,
        Err(e) => {
            error!("Failed to create VPN client: {}", e);
            return HttpResponse::InternalServerError()
                .json(ErrorResponse { message: "error".to_string() });
        }
    };

    match vpn.auth_member(member_data.member_id.clone(), None).await {
        Ok(_) => HttpResponse::Ok().finish(),
        Err(e) => {
            error!("Failed to authorize member: {}", e);
            if e.to_string().contains("not found") {
                return HttpResponse::NotFound()
                    .json(ErrorResponse { message: "member not found".to_string() });
            }
            HttpResponse::InternalServerError()
                .json(ErrorResponse { message: "error".to_string() })
        }
    }
}

#[get("/vpn-id")]
pub async fn vpn_id() -> impl Responder {
    let config = Config::instance();
    HttpResponse::Ok().json(VpnIdResponse {
        vpn_id: config.vpn_config().network_id().to_string(),
    })
}