use serde::Serialize;

#[derive(Serialize)]
pub struct VpnIdResponse {
    pub vpn_id: String,
}
