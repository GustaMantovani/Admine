use std::net::IpAddr;

use crate::errors::VpnError;
use async_trait::async_trait;
use std::fmt;

pub enum MemberVpnStatus {
    Online,
    Offline,
    Unknown,
}

pub enum MemberVpnAuthStatus {
    Authenticated,
    NotAuthenticated,
    Unknown,
}

impl fmt::Display for MemberVpnAuthStatus {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        let status_str = match self {
            MemberVpnAuthStatus::Authenticated => "Authenticated",
            MemberVpnAuthStatus::NotAuthenticated => "NotAuthenticated",
            MemberVpnAuthStatus::Unknown => "Unknown",
        };
        write!(f, "{}", status_str)
    }
}

#[async_trait]

pub trait TVpnClient {
    async fn auth_member(
        &self,
        member_id: String,
        member_token: Option<String>,
    ) -> Result<(), VpnError>;
    async fn delete_member(&self, member_id: String) -> Result<(), VpnError>;
    async fn member_vpn_status(&self, member_id: String) -> Result<MemberVpnStatus, VpnError>;
    async fn member_vpn_auth_status(
        &self,
        member_id: String,
    ) -> Result<MemberVpnAuthStatus, VpnError>;
    async fn get_member_ips_in_vpn(&self, member_id: String) -> Result<Vec<IpAddr>, VpnError>;
}
