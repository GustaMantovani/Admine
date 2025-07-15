use zerotier_central_api::apis::configuration::Configuration;

use crate::errors::VpnError;
use crate::vpn::public_ip::PublicIp;
use crate::vpn::vpn::TVpnClient;
use crate::vpn::zerotier_vpn::ZerotierVpn;

#[derive(Clone, Debug)]
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
    ) -> Result<Box<dyn TVpnClient + Send + Sync>, VpnError> {
        match vpn_type {
            VpnType::Zerotier => {
                let mut config = Configuration::new();
                config.base_path = api_url;
                config.api_key = Some(zerotier_central_api::apis::configuration::ApiKey {
                    prefix: None,
                    key: api_key,
                });

                print!("{:?}", config);

                Ok(Box::new(ZerotierVpn::new(config, network_id)))
            }

            VpnType::PublicIp => Ok(Box::new(PublicIp::new())),
        }
    }
}
