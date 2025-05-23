use thiserror::Error;

#[derive(Error, Debug)]
pub enum PubSubError {
    #[error("Failed to get message payload: {0}")]
    MessageError(String),
}

#[derive(Error, Debug)]
pub enum VpnError {
    #[error("Member not found: {0}")]
    MemberNotFoundError(String),

    #[error("Deletion error: {0}")]
    DeletionError(String),

    #[error("Authorization error: {0}")]
    MemberUpdateError(String),
}
