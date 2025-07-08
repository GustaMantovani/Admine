use std::net::IpAddr;
use crate::errors::VpnError;
use async_trait::async_trait;

#[async_trait]
pub trait TVpnClient {
    async fn auth_member(
        &self,
        member_id: String,
        member_token: Option<String>,
    ) -> Result<(), VpnError>;
    async fn delete_member(&self, member_id: String) -> Result<(), VpnError>;
    async fn get_member_ips_in_vpn(&self, member_id: String) -> Result<Vec<IpAddr>, VpnError>;
}
