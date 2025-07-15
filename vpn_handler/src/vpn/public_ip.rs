use crate::{errors::VpnError, vpn::vpn::TVpnClient};
use async_trait::async_trait;
use std::net::IpAddr;

pub struct PublicIp {}

impl PublicIp {
    pub fn new() -> Self {
        Self {}
    }
}

#[async_trait]
impl TVpnClient for PublicIp {
    async fn delete_member(&self, _member_id: String) -> Result<(), VpnError> {
        Ok(())
    }

    async fn get_member_ips_in_vpn(&self, _member_id: String) -> Result<Vec<IpAddr>, VpnError> {
        match public_ip_address::perform_lookup(None).await {
            Ok(response) => Ok(vec![response.ip]),
            Err(e) => Err(VpnError::MemberNotFoundError(e.to_string())),
        }
    }

    async fn auth_member(
        &self,
        _member_id: String,
        _member_token: Option<String>,
    ) -> Result<(), VpnError> {
        Ok(())
    }
}
