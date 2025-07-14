use super::vpn::TVpnClient;
use crate::errors::VpnError;
use async_trait::async_trait;
use std::net::IpAddr;
use std::str::FromStr;
use zerotier_central_api::apis::configuration::Configuration;
use zerotier_central_api::apis::network_member_api::{
    delete_network_member, get_network_member, update_network_member
};

pub struct ZerotierVpn {
    config: Configuration,
    network_id: String,
}

impl ZerotierVpn {
    pub fn new(config: Configuration, network_id: String) -> Self {
        Self { config, network_id }
    }
}

#[async_trait]
impl TVpnClient for ZerotierVpn {
    async fn delete_member(&self, member_id: String) -> Result<(), VpnError> {
        match get_network_member(&self.config, &self.network_id, &member_id).await {
            Ok(_) => {
                match delete_network_member(&self.config, &self.network_id, &member_id).await {
                    Ok(_) => return Ok(()),
                    Err(e) => return Err(VpnError::DeletionError(e.to_string())),
                };
            }
            Err(e) => return Err(VpnError::MemberNotFoundError(e.to_string())),
        };
    }

    async fn get_member_ips_in_vpn(&self, member_id: String) -> Result<Vec<IpAddr>, VpnError> {
        let member = match get_network_member(&self.config, &self.network_id, &member_id).await {
            Ok(m) => m,
            Err(e) => return Err(VpnError::MemberNotFoundError(e.to_string())),
        };

        if let Some(config) = member.config {
            if let Some(ip_assignments) = config.ip_assignments {
                if !ip_assignments.is_empty() {
                    return Ok(ip_assignments
                        .iter()
                        .filter_map(|ip| IpAddr::from_str(ip).ok())
                        .collect());
                }
            }
        }

        Ok(vec![])
    }

    async fn auth_member(
        &self,
        member_id: String,
        _member_token: Option<String>,
    ) -> Result<(), VpnError> {
        print!("{}", &self.network_id);
        print!("{:?}", &self.config);
        let mut member = match get_network_member(&self.config, &self.network_id, &member_id).await
        {
            Ok(m) => m,
            Err(e) => return Err(VpnError::MemberNotFoundError(e.to_string())),
        };

        // Check if already authorized
        if let Some(ref config) = member.config {
            if let Some(true) = config.authorized {
                return Ok(());
            }
        }

        // Set authorization
        if let Some(mut member_config) = member.config {
            member_config.authorized = Some(true);
            member.config = Some(member_config);

            if let Some(node_id) = member.node_id.clone() {
                match update_network_member(
                    &self.config,
                    &self.network_id,
                    &node_id,
                    member.clone(),
                )
                .await
                {
                    Ok(_) => return Ok(()),
                    Err(e) => return Err(VpnError::MemberUpdateError(e.to_string())),
                };
            }
        }

        Ok(())
    }
}
