use super::vpn::TVpnClient;
use crate::errors::VpnError;
use async_trait::async_trait;
use std::net::IpAddr;
use std::str::FromStr;
use zerotier_central_api::apis::configuration::Configuration;
use zerotier_central_api::apis::network_member_api::{
    delete_network_member, get_network_member, update_network_member, DeleteNetworkMemberError,
    GetNetworkMemberError, UpdateNetworkMemberError,
};
use zerotier_central_api::apis::Error;

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
                    Ok(_) => Ok(()),
                    Err(Error::ResponseError(response)) => match &response.entity {
                        Some(DeleteNetworkMemberError::Status403()) => Err(
                            VpnError::InternalError("API authentication failed".to_string()),
                        ),
                        Some(DeleteNetworkMemberError::Status401()) => Err(
                            VpnError::InternalError("API authentication failed".to_string()),
                        ),
                        _ => Err(VpnError::DeletionError(response.content.clone())),
                    },
                    Err(e) => Err(VpnError::InternalError(format!("Network error: {}", e))),
                }
            }
            Err(Error::ResponseError(response)) => match &response.entity {
                Some(GetNetworkMemberError::Status404()) => {
                    Err(VpnError::MemberNotFoundError(response.content.clone()))
                }
                Some(GetNetworkMemberError::Status403()) => Err(VpnError::InternalError(
                    "API authentication failed".to_string(),
                )),
                Some(GetNetworkMemberError::Status401()) => Err(VpnError::InternalError(
                    "API authentication failed".to_string(),
                )),
                _ => Err(VpnError::InternalError(format!(
                    "Network error: {}",
                    response.content
                ))),
            },
            Err(e) => Err(VpnError::InternalError(format!("Network error: {}", e))),
        }
    }

    async fn get_member_ips_in_vpn(&self, member_id: String) -> Result<Vec<IpAddr>, VpnError> {
        if member_id.len() > 0 {
            let member = match get_network_member(&self.config, &self.network_id, &member_id).await
            {
                Ok(m) => m,
                Err(Error::ResponseError(response)) => match &response.entity {
                    Some(GetNetworkMemberError::Status404()) => {
                        return Err(VpnError::MemberNotFoundError(response.content.clone()));
                    }
                    Some(GetNetworkMemberError::Status403()) => {
                        return Err(VpnError::InternalError(
                            "API authentication failed".to_string(),
                        ));
                    }
                    Some(GetNetworkMemberError::Status401()) => {
                        return Err(VpnError::InternalError(
                            "API authentication failed".to_string(),
                        ));
                    }
                    _ => {
                        return Err(VpnError::InternalError(format!(
                            "Network error: {}",
                            response.content
                        )))
                    }
                },
                Err(e) => return Err(VpnError::InternalError(format!("Network error: {}", e))),
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
            Err(Error::ResponseError(response)) => match &response.entity {
                Some(GetNetworkMemberError::Status404()) => {
                    return Err(VpnError::MemberNotFoundError(response.content.clone()));
                }
                Some(GetNetworkMemberError::Status403()) => {
                    return Err(VpnError::InternalError(
                        "API authentication failed".to_string(),
                    ));
                }
                Some(GetNetworkMemberError::Status401()) => {
                    return Err(VpnError::InternalError(
                        "API authentication failed".to_string(),
                    ));
                }
                _ => {
                    return Err(VpnError::InternalError(format!(
                        "Network error: {}",
                        response.content
                    )))
                }
            },
            Err(e) => return Err(VpnError::InternalError(format!("Network error: {}", e))),
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
                    Err(Error::ResponseError(response)) => match &response.entity {
                        Some(UpdateNetworkMemberError::Status403()) => {
                            return Err(VpnError::InternalError(
                                "API authentication failed".to_string(),
                            ));
                        }
                        Some(UpdateNetworkMemberError::Status401()) => {
                            return Err(VpnError::InternalError(
                                "API authentication failed".to_string(),
                            ));
                        }
                        _ => return Err(VpnError::MemberUpdateError(response.content.clone())),
                    },
                    Err(e) => return Err(VpnError::InternalError(format!("Network error: {}", e))),
                };
            }
        }

        Ok(())
    }
}
