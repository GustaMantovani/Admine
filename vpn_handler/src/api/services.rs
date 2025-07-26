use crate::{
    app_context::AppContext,
    errors::VpnError,
    models::api_models::{AuthMemberRequest, ErrorResponse, ServerIpResponse, VpnIdResponse},
};
use actix_web::{get, post, web, HttpResponse, Responder};
use log::{error, info};

fn map_error_to_http_response(error: VpnError) -> HttpResponse {
    match error {
        VpnError::MemberNotFoundError(_) => {
            error!("Member not found: {}", error);
            HttpResponse::NotFound().json(ErrorResponse {
                message: error.to_string(),
            })
        }
        VpnError::DeletionError(_) => {
            error!("Failed to delete member: {}", error);
            HttpResponse::InternalServerError().json(ErrorResponse {
                message: error.to_string(),
            })
        }
        VpnError::MemberUpdateError(_) => {
            error!("Failed to update/authorize member: {}", error);
            HttpResponse::InternalServerError().json(ErrorResponse {
                message: error.to_string(),
            })
        }
        VpnError::InternalError(_) => {
            error!("Internal VPN error: {}", error);
            HttpResponse::InternalServerError().json(ErrorResponse {
                message: error.to_string(),
            })
        }
    }
}

#[get("/server-ip")]
pub async fn server_ip() -> impl Responder {
    let server_vpn_id = AppContext::instance()
        .storage()
        .get("server_vpn_id")
        .unwrap_or("".to_string());

    match AppContext::instance()
        .vpn_client()
        .get_member_ips_in_vpn(server_vpn_id)
        .await
    {
        Ok(ips) => HttpResponse::Ok().json(ServerIpResponse { server_ips: ips }),
        Err(vpn_error) => map_error_to_http_response(vpn_error),
    }
}

#[post("/auth-member")]
pub async fn auth_member(member_data: web::Json<AuthMemberRequest>) -> impl Responder {
    info!("Authorizing member: {}", member_data.member_id);

    let vpn = AppContext::instance().vpn_client();

    match vpn.auth_member(member_data.member_id.clone(), None).await {
        Ok(_) => {
            info!("Member {} authorized successfully", member_data.member_id);
            HttpResponse::NoContent().finish()
        }
        Err(vpn_error) => {
            error!(
                "Failed to authorize member {}: {}",
                member_data.member_id, vpn_error
            );
            map_error_to_http_response(vpn_error)
        }
    }
}

#[get("/vpn-id")]
pub async fn vpn_id() -> impl Responder {
    let config = AppContext::instance().config();
    HttpResponse::Ok().json(VpnIdResponse {
        vpn_id: config.vpn_config.network_id.clone(),
    })
}
