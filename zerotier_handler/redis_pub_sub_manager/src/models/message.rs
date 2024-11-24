use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug)]
pub struct AdmineMessage {
    pub tags: Vec<String>,
    pub message: String,
}