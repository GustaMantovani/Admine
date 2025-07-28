use crate::errors::VpnError;
use crate::vpn::public_ip::PublicIp;
use crate::vpn::vpn::DynVpn;
use crate::vpn::zerotier_vpn::ZerotierVpn;
use serde::Deserialize;
use strum::EnumString;
use zerotier_central_api::apis::configuration::Configuration;

#[derive(Clone, Debug, EnumString, Deserialize)]
pub enum VpnType {
    Zerotier,
    PublicIp,
}

pub struct VpnFactory;

impl VpnFactory {
    pub fn create_vpn(
        vpn_type: VpnType,
        api_url: String,
        api_key: String,
        network_id: String,
    ) -> Result<DynVpn, VpnError> {
        match vpn_type {
            VpnType::Zerotier => {
                let mut config = Configuration::new();
                config.base_path = api_url;
                config.api_key = Some(zerotier_central_api::apis::configuration::ApiKey {
                    prefix: None,
                    key: api_key.clone(),
                });

                config.bearer_access_token = Some(api_key.clone());

                print!("{:?}", config);

                Ok(Box::new(ZerotierVpn::new(config, network_id)))
            }

            VpnType::PublicIp => Ok(Box::new(PublicIp::new())),
        }
    }
}
