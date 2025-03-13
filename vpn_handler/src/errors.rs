use thiserror::Error;

#[derive(Error, Debug)]
pub enum PubSubError {
    #[error("Failed to connect to Redis: {0}")]
    ConnectionError(String),

    #[error("Failed to subscribe to channel: {0}")]
    SubscriptionError(String),

    #[error("Failed to publish message: {0}")]
    PublishError(String),

    #[error("Failed to receive message: {0}")]
    ReceiveError(String),

    #[error("Unknown provider: {0}")]
    UnknownProvider(String),

    #[error("Failed to get message payload: {0}")]
    MessageError(String),
}

#[derive(Error, Debug)]
pub enum VpnError {
    #[error("VPN client error: {0}")]
    VpnClientError(String),

    #[error("Member not found: {0}")]
    MemberNotFoundError(String),

    #[error("Authentication error: {0}")]
    AuthError(String),

    #[error("Deletion error: {0}")]
    DeletionError(String),

    #[error("Status error: {0}")]
    StatusError(String),

    #[error("Authorization error: {0}")]
    MemberUpdateError(String),
}
