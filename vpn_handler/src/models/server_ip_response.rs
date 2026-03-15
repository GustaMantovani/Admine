use serde::Serialize;
use std::net::IpAddr;

#[derive(Serialize)]
pub struct ServerIpResponse {
    pub server_ips: Vec<IpAddr>,
}
