use crate::errors::VpnError;
use crate::vpn::vpn::TVpnClient;
use crate::vpn::zerotier::apis::configuration::Configuration;
use crate::vpn::zerotier_vpn::ZerotierVpn;

#[derive(Clone, Debug)]
pub enum VpnType {
    Zerotier,
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
            VpnType::Zerotier => Ok(Box::new(ZerotierVpn::new(
                Configuration::new(api_url, api_key),
                network_id,
            ))),
        }
    }
}
