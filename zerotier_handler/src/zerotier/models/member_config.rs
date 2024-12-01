use crate::zerotier::models;
use serde::{Deserialize, Serialize};

#[derive(Clone, Default, Debug, PartialEq, Serialize, Deserialize)]
pub struct MemberConfig {
    /// Allow the member to be a bridge on the network
    #[serde(rename = "activeBridge", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub active_bridge: Option<Option<bool>>,
    /// Is the member authorized on the network
    #[serde(rename = "authorized", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub authorized: Option<Option<bool>>,
    #[serde(rename = "capabilities", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub capabilities: Option<Option<Vec<i32>>>,
    /// Time the member was created or first tried to join the network
    #[serde(rename = "creationTime", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub creation_time: Option<Option<i64>>,
    /// ID of the member node.  This is the 10 digit identifier that identifies a ZeroTier node.
    #[serde(rename = "id", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub id: Option<Option<String>>,
    /// Public Key of the member's Identity
    #[serde(rename = "identity", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub identity: Option<Option<String>>,
    /// List of assigned IP addresses
    #[serde(rename = "ipAssignments", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub ip_assignments: Option<Option<Vec<String>>>,
    /// Time the member was authorized on the network
    #[serde(rename = "lastAuthorizedTime", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub last_authorized_time: Option<Option<i64>>,
    /// Time the member was deauthorized on the network
    #[serde(rename = "lastDeauthorizedTime", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub last_deauthorized_time: Option<Option<i64>>,
    /// Exempt this member from the IP auto assignment pool on a Network
    #[serde(rename = "noAutoAssignIps", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub no_auto_assign_ips: Option<Option<bool>>,
    /// Member record revision count
    #[serde(rename = "revision", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub revision: Option<Option<i32>>,
    /// Allow the member to be authorized without OIDC/SSO
    #[serde(rename = "ssoExempt", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub sso_exempt: Option<Option<bool>>,
    /// Array of 2 member tuples of tag [ID, tag value]
    #[serde(rename = "tags", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub tags: Option<Option<Vec<Vec<models::MemberConfigTagsInnerInner>>>>,
    /// Major version of the client
    #[serde(rename = "vMajor", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub v_major: Option<Option<i32>>,
    /// Minor version of the client
    #[serde(rename = "vMinor", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub v_minor: Option<Option<i32>>,
    /// Revision number of the client
    #[serde(rename = "vRev", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub v_rev: Option<Option<i32>>,
    /// Protocol version of the client
    #[serde(rename = "vProto", default, with = "::serde_with::rust::double_option", skip_serializing_if = "Option::is_none")]
    pub v_proto: Option<Option<i32>>,
}
