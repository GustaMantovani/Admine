use super::vpn::TVpnClient;
use crate::errors::VpnError;
use async_trait::async_trait;
use reqwest::{Client, StatusCode};
use serde::Deserialize;
use std::net::IpAddr;
use std::str::FromStr;

pub struct TailscaleVpn {
    client: Client,
    base_url: String,
    api_key: String,
    _tailnet: String,
}

impl TailscaleVpn {
    pub fn new(base_url: String, api_key: String, tailnet: String) -> Self {
        Self {
            client: Client::new(),
            base_url,
            api_key,
            _tailnet: tailnet,
        }
    }
}

#[derive(Deserialize)]
struct TailscaleDevice {
    addresses: Vec<String>,
}

#[async_trait]
impl TVpnClient for TailscaleVpn {
    async fn auth_member(
        &self,
        member_id: String,
        _member_token: Option<String>,
    ) -> Result<(), VpnError> {
        let url = format!("{}/device/{}/authorized", self.base_url, member_id);

        let resp = self
            .client
            .post(&url)
            .bearer_auth(&self.api_key)
            .json(&serde_json::json!({"authorized": true}))
            .send()
            .await
            .map_err(|e| VpnError::InternalError(format!("Network error: {}", e)))?;

        match resp.status() {
            StatusCode::OK | StatusCode::NO_CONTENT => Ok(()),
            StatusCode::NOT_FOUND => Err(VpnError::MemberNotFoundError(format!(
                "Device {} not found",
                member_id
            ))),
            StatusCode::UNAUTHORIZED | StatusCode::FORBIDDEN => {
                Err(VpnError::InternalError("API authentication failed".to_string()))
            }
            s => Err(VpnError::MemberUpdateError(format!(
                "Unexpected status: {}",
                s
            ))),
        }
    }

    async fn delete_member(&self, member_id: String) -> Result<(), VpnError> {
        let url = format!("{}/device/{}", self.base_url, member_id);

        let resp = self
            .client
            .delete(&url)
            .bearer_auth(&self.api_key)
            .send()
            .await
            .map_err(|e| VpnError::InternalError(format!("Network error: {}", e)))?;

        match resp.status() {
            StatusCode::OK | StatusCode::NO_CONTENT => Ok(()),
            StatusCode::NOT_FOUND => Err(VpnError::MemberNotFoundError(format!(
                "Device {} not found",
                member_id
            ))),
            StatusCode::UNAUTHORIZED | StatusCode::FORBIDDEN => {
                Err(VpnError::InternalError("API authentication failed".to_string()))
            }
            s => Err(VpnError::DeletionError(format!("Unexpected status: {}", s))),
        }
    }

    async fn get_member_ips_in_vpn(&self, member_id: String) -> Result<Vec<IpAddr>, VpnError> {
        if member_id.is_empty() {
            return Ok(vec![]);
        }

        let url = format!("{}/device/{}", self.base_url, member_id);

        let resp = self
            .client
            .get(&url)
            .bearer_auth(&self.api_key)
            .send()
            .await
            .map_err(|e| VpnError::InternalError(format!("Network error: {}", e)))?;

        match resp.status() {
            StatusCode::OK => {
                let device: TailscaleDevice = resp
                    .json()
                    .await
                    .map_err(|e| VpnError::InternalError(format!("Failed to parse response: {}", e)))?;

                Ok(device
                    .addresses
                    .iter()
                    .filter_map(|a| IpAddr::from_str(a).ok())
                    .collect())
            }
            StatusCode::NOT_FOUND => Err(VpnError::MemberNotFoundError(format!(
                "Device {} not found",
                member_id
            ))),
            StatusCode::UNAUTHORIZED | StatusCode::FORBIDDEN => {
                Err(VpnError::InternalError("API authentication failed".to_string()))
            }
            s => Err(VpnError::InternalError(format!("Unexpected status: {}", s))),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::errors::VpnError;
    use async_trait::async_trait;
    use mockall::mock;
    use mockall::predicate::*;

    mock! {
        TailscaleClient {}

        #[async_trait]
        impl TVpnClient for TailscaleClient {
            async fn auth_member(&self, member_id: String, member_token: Option<String>) -> Result<(), VpnError>;
            async fn delete_member(&self, member_id: String) -> Result<(), VpnError>;
            async fn get_member_ips_in_vpn(&self, member_id: String) -> Result<Vec<IpAddr>, VpnError>;
        }
    }

    #[tokio::test]
    async fn test_get_member_ips_empty_id_returns_empty_vec() {
        let vpn = TailscaleVpn::new(
            "https://api.tailscale.com".to_string(),
            "dummy_key".to_string(),
            "example.com".to_string(),
        );
        let result = vpn.get_member_ips_in_vpn("".to_string()).await;
        assert!(result.is_ok());
        assert!(result.unwrap().is_empty());
    }

    #[tokio::test]
    async fn test_auth_member_success() {
        let mut mock = MockTailscaleClient::new();
        mock.expect_auth_member()
            .with(eq("12345678".to_string()), eq(None))
            .times(1)
            .returning(|_, _| Ok(()));

        let result = mock.auth_member("12345678".to_string(), None).await;
        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_auth_member_not_found() {
        let mut mock = MockTailscaleClient::new();
        mock.expect_auth_member()
            .times(1)
            .returning(|_, _| Err(VpnError::MemberNotFoundError("Device not found".to_string())));

        let result = mock.auth_member("99999999".to_string(), None).await;
        assert!(matches!(result, Err(VpnError::MemberNotFoundError(_))));
    }

    #[tokio::test]
    async fn test_auth_member_internal_error() {
        let mut mock = MockTailscaleClient::new();
        mock.expect_auth_member()
            .times(1)
            .returning(|_, _| Err(VpnError::InternalError("API authentication failed".to_string())));

        let result = mock.auth_member("12345678".to_string(), None).await;
        assert!(matches!(result, Err(VpnError::InternalError(_))));
    }

    #[tokio::test]
    async fn test_delete_member_success() {
        let mut mock = MockTailscaleClient::new();
        mock.expect_delete_member()
            .with(eq("12345678".to_string()))
            .times(1)
            .returning(|_| Ok(()));

        let result = mock.delete_member("12345678".to_string()).await;
        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_delete_member_not_found() {
        let mut mock = MockTailscaleClient::new();
        mock.expect_delete_member()
            .times(1)
            .returning(|_| Err(VpnError::MemberNotFoundError("Device not found".to_string())));

        let result = mock.delete_member("00000000".to_string()).await;
        assert!(matches!(result, Err(VpnError::MemberNotFoundError(_))));
    }

    #[tokio::test]
    async fn test_get_member_ips_success() {
        let mut mock = MockTailscaleClient::new();
        let ips: Vec<IpAddr> = vec!["100.64.0.1".parse().unwrap()];
        mock.expect_get_member_ips_in_vpn()
            .with(eq("12345678".to_string()))
            .times(1)
            .returning(move |_| Ok(ips.clone()));

        let result = mock.get_member_ips_in_vpn("12345678".to_string()).await;
        assert!(result.is_ok());
        assert_eq!(result.unwrap().len(), 1);
    }

    #[tokio::test]
    async fn test_get_member_ips_not_found() {
        let mut mock = MockTailscaleClient::new();
        mock.expect_get_member_ips_in_vpn()
            .times(1)
            .returning(|_| Err(VpnError::MemberNotFoundError("Device not found".to_string())));

        let result = mock.get_member_ips_in_vpn("12345678".to_string()).await;
        assert!(matches!(result, Err(VpnError::MemberNotFoundError(_))));
    }
}
