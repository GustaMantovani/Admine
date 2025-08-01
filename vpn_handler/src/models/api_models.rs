use std::net::IpAddr;

use serde::{Deserialize, Serialize};

#[derive(Deserialize)]
pub struct AuthMemberRequest {
    pub member_id: String,
}

#[derive(Serialize)]
pub struct ServerIpResponse {
    pub server_ips: Vec<IpAddr>,
}

#[derive(Serialize)]
pub struct VpnIdResponse {
    pub vpn_id: String,
}

#[derive(Serialize)]
pub struct ErrorResponse {
    pub message: String,
}
