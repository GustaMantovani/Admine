/*
 * ZeroTier Central API
 *
 * ZeroTier Central Network Management Portal API.<p>All API requests must have an API token header specified in the <code>Authorization: token xxxxx</code> format.  You can generate your API key by logging into <a href=\"https://my.zerotier.com\">ZeroTier Central</a> and creating a token on the Account page.</p><p>eg. <code>curl -X GET -H \"Authorization: token xxxxx\" https://api.zerotier.com/api/v1/network</code></p><p><h3>Rate Limiting</h3></p><p>The ZeroTier Central API implements rate limiting.  Paid users are limited to 100 requests per second.  Free users are limited to 20 requests per second.</p> <p> You can get the OpenAPI spec here as well: <code>https://docs.zerotier.com/api/central/ref-v1.json</code></p>
 *
 * The version of the OpenAPI document: v1
 * 
 * Generated by: https://openapi-generator.tech
 */

use crate::models;
use serde::{Deserialize, Serialize};

#[derive(Clone, Default, Debug, PartialEq, Serialize, Deserialize)]
pub struct NetworkSsoConfig {
    /// SSO enabled/disabled on network
    #[serde(rename = "enabled", skip_serializing_if = "Option::is_none")]
    pub enabled: Option<bool>,
    /// SSO mode.  One of: `default`, `email`, `group`
    #[serde(rename = "mode", skip_serializing_if = "Option::is_none")]
    pub mode: Option<String>,
    /// SSO client ID.  Client ID must be already configured in the Org
    #[serde(rename = "clientId", skip_serializing_if = "Option::is_none")]
    pub client_id: Option<String>,
    /// URL of the OIDC issuer
    #[serde(rename = "issuer", skip_serializing_if = "Option::is_none")]
    pub issuer: Option<String>,
    /// Provider type
    #[serde(rename = "provider", skip_serializing_if = "Option::is_none")]
    pub provider: Option<String>,
    /// Authorization URL endpoint
    #[serde(rename = "authorizationEndpoint", skip_serializing_if = "Option::is_none")]
    pub authorization_endpoint: Option<String>,
    /// List of email addresses or group memberships that may SSO auth onto the network
    #[serde(rename = "allowList", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub allow_list: Option<Option<Vec<String>>>,
}

impl NetworkSsoConfig {
    pub fn new() -> NetworkSsoConfig {
        NetworkSsoConfig {
            enabled: None,
            mode: None,
            client_id: None,
            issuer: None,
            provider: None,
            authorization_endpoint: None,
            allow_list: None,
        }
    }
}

