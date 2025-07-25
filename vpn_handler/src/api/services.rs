use crate::{
    app_context::AppContext,
    models::api_models::{AuthMemberRequest, ErrorResponse, ServerIpResponse, VpnIdResponse},
};
use actix_web::{get, post, web, HttpResponse, Responder};
use log::info;

#[get("/status")]
async fn status() -> impl Responder {
    HttpResponse::Ok().body("1")
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
        Err(e) => HttpResponse::InternalServerError().json(ErrorResponse {
            message: e.to_string(),
        }),
    }
}

#[post("/auth-member")]
pub async fn auth_member(member_data: web::Json<AuthMemberRequest>) -> impl Responder {
    info!("Authorizing member: {}", member_data.member_id);

    let vpn = AppContext::instance().vpn_client();

    match vpn.auth_member(member_data.member_id.clone(), None).await {
        Ok(_) => HttpResponse::Ok().finish(),
        Err(e) => HttpResponse::InternalServerError().json(ErrorResponse {
            message: e.to_string(),
        }),
    }
}

#[get("/vpn-id")]
pub async fn vpn_id() -> impl Responder {
    let config = AppContext::instance().config();
    HttpResponse::Ok().json(VpnIdResponse {
        vpn_id: config.vpn_config.network_id.clone(),
    })
}
