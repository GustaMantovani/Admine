use serde::{Deserialize, Serialize};
use serde_json;

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct AdmineMessage {
    pub tags: Vec<String>,
    pub message: String,
}

impl AdmineMessage {
    pub fn to_json_string(&self) -> Result<String, serde_json::Error> {
        serde_json::to_string(self)
    }

    pub fn from_json_string(json_str: &str) -> Result<AdmineMessage, serde_json::Error> {
        serde_json::from_str(json_str)
    }
}
