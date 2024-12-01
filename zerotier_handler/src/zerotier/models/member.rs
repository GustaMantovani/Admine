use crate::zerotier::models;
use serde::{Deserialize, Serialize};
use std::fmt;


#[derive(Clone, Default, Debug, PartialEq, Serialize, Deserialize)]
pub struct Member {
    /// concatenation of network ID and member ID
    #[serde(rename = "id", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub id: Option<Option<String>>,
    #[serde(rename = "clock", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub clock: Option<Option<i64>>,
    #[serde(rename = "networkId", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub network_id: Option<Option<String>>,
    /// ZeroTier ID of the member
    #[serde(rename = "nodeId", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub node_id: Option<Option<String>>,
    #[serde(rename = "controllerId", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub controller_id: Option<Option<String>>,
    /// Whether or not the member is hidden in the UI
    #[serde(rename = "hidden", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub hidden: Option<Option<bool>>,
    /// User defined name of the member
    #[serde(rename = "name", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub name: Option<Option<String>>,
    /// User defined description of the member
    #[serde(rename = "description", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub description: Option<Option<String>>,
    #[serde(rename = "config", skip_serializing_if = "Option::is_none")]
    pub config: Option<Box<models::MemberConfig>>,
    /// Last seen time of the member (milliseconds since epoch).  Note: This data is considered ephemeral and may be reset to 0 at any time without warning.
    #[serde(rename = "lastOnline", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub last_online: Option<Option<i64>>,
    /// Time the member last checked in with the network controller in milliseconds since epoch. Note: This data is considered ephemeral and may be reset to 0 at any time without warning.
    #[serde(rename = "lastSeen", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub last_seen: Option<Option<i64>>,
    /// IP address the member last spoke to the controller via (milliseconds since epoch).  Note: This data is considered ephemeral and may be reset to 0 at any time without warning.
    #[serde(rename = "physicalAddress", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub physical_address: Option<Option<String>>,
    /// ZeroTier version the member is running
    #[serde(rename = "clientVersion", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub client_version: Option<Option<String>>,
    /// ZeroTier protocol version
    #[serde(rename = "protocolVersion", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub protocol_version: Option<Option<i32>>,
    /// Whether or not the client version is new enough to support the rules engine (1.4.0+)
    #[serde(rename = "supportsRulesEngine", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub supports_rules_engine: Option<Option<bool>>,
}

impl fmt::Display for Member {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(
            f,
            "Member {{
    ID: {},
    Clock: {},
    Network ID: {},
    Node ID: {},
    Controller ID: {},
    Hidden: {},
    Name: {},
    Description: {},
    Last Online: {},
    Last Seen: {},
    Physical Address: {},
    Client Version: {},
    Protocol Version: {},
    Supports Rules Engine: {}
}}",
            self.id
                .as_ref()
                .and_then(|x| x.clone())
                .unwrap_or_else(|| String::from("None")),
            self.clock
                .as_ref()
                .and_then(|x| x.map(|v| v.to_string()))
                .unwrap_or_else(|| String::from("None")),
            self.network_id
                .as_ref()
                .and_then(|x| x.clone())
                .unwrap_or_else(|| String::from("None")),
            self.node_id
                .as_ref()
                .and_then(|x| x.clone())
                .unwrap_or_else(|| String::from("None")),
            self.controller_id
                .as_ref()
                .and_then(|x| x.clone())
                .unwrap_or_else(|| String::from("None")),
            self.hidden
                .as_ref()
                .and_then(|x| x.map(|v| v.to_string()))
                .unwrap_or_else(|| String::from("None")),
            self.name
                .as_ref()
                .and_then(|x| x.clone())
                .unwrap_or_else(|| String::from("None")),
            self.description
                .as_ref()
                .and_then(|x| x.clone())
                .unwrap_or_else(|| String::from("None")),
            self.last_online
                .as_ref()
                .and_then(|x| x.map(|v| v.to_string()))
                .unwrap_or_else(|| String::from("None")),
            self.last_seen
                .as_ref()
                .and_then(|x| x.map(|v| v.to_string()))
                .unwrap_or_else(|| String::from("None")),
            self.physical_address
                .as_ref()
                .and_then(|x| x.clone())
                .unwrap_or_else(|| String::from("None")),
            self.client_version
                .as_ref()
                .and_then(|x| x.clone())
                .unwrap_or_else(|| String::from("None")),
            self.protocol_version
                .as_ref()
                .and_then(|x| x.map(|v| v.to_string()))
                .unwrap_or_else(|| String::from("None")),
            self.supports_rules_engine
                .as_ref()
                .and_then(|x| x.map(|v| v.to_string()))
                .unwrap_or_else(|| String::from("None")),
        )
    }
}
