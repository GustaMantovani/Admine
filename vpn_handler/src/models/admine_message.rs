use getset::{Getters, Setters};
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, Clone, Getters, Setters)]
#[getset(get = "pub", set = "pub")]
pub struct AdmineMessage {
    origin: String,
    tags: Vec<String>,
    message: String,
}

impl AdmineMessage {
    pub fn new(origin: String, tags: Vec<String>, message: String) -> Self {
        Self {
            origin,
            tags,
            message,
        }
    }

    pub fn has_tag(&self, tag: &str) -> bool {
        self.tags.contains(&tag.to_string())
    }
}
