use crate::{
    config::Config,
    errors::VpnError,
    models::{
        auth_member_request::AuthMemberRequest, error_response::ErrorResponse,
        server_ip_response::ServerIpResponse, vpn_id_response::VpnIdResponse,
    },
    persistence::key_value_storage::DynKeyValueStore,
    vpn::vpn::DynVpn,
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

#[get("/server-ips")]
pub async fn server_ip(
    storage: web::Data<DynKeyValueStore>,
    vpn_client: web::Data<DynVpn>,
) -> impl Responder {
    let server_vpn_id = storage.get("server_member_id").unwrap_or_default();

    match vpn_client.get_member_ips_in_vpn(server_vpn_id).await {
        Ok(ips) => HttpResponse::Ok().json(ServerIpResponse { server_ips: ips }),
        Err(vpn_error) => map_error_to_http_response(vpn_error),
    }
}

#[post("/auth-member")]
pub async fn auth_member(
    member_data: web::Json<AuthMemberRequest>,
    vpn_client: web::Data<DynVpn>,
) -> impl Responder {
    info!("Authorizing member: {}", member_data.member_id);

    match vpn_client
        .auth_member(member_data.member_id.clone(), None)
        .await
    {
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
pub async fn vpn_id(config: web::Data<Config>) -> impl Responder {
    HttpResponse::Ok().json(VpnIdResponse {
        vpn_id: config.vpn_config().network_id().clone(),
    })
}
