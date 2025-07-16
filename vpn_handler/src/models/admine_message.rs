use serde::{Deserialize, Serialize};
use std::fmt;

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct AdmineMessage {
    pub origin: String,
    pub tags: Vec<String>,
    pub message: String,
}

impl fmt::Display for AdmineMessage {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(
            f,
            "AdmineMessage {{ origin: \"{}\", tags: [{}], message: \"{}\" }}",
            self.origin,
            self.tags.join(", "),
            self.message
        )
    }
}
