use thiserror::Error;

#[derive(Error, Debug)]
pub enum PubSubError {
    #[error("Failed to get message payload: {0}")]
    MessageError(String),

    #[error("Connection error: {0}")]
    ConnectionError(String),

    #[error("Creation error: {0}")]
    CreationError(String),
}

#[derive(Error, Debug)]
pub enum VpnError {
    #[error("Member not found: {0}")]
    MemberNotFoundError(String),

    #[error("Deletion error: {0}")]
    DeletionError(String),

    #[error("Authorization error: {0}")]
    MemberUpdateError(String),

    #[error("Internal VPN error: {0}")]
    InternalError(String),
}
