use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug)]
pub struct AdmineMessage {
    pub tags: Vec<String>,
    pub message: String,
}

impl AdmineMessage {
    pub fn new(tags: Vec<String>, message: &str) -> Self {
        AdmineMessage {
            tags,
            message: message.to_string(),
        }
    }
}