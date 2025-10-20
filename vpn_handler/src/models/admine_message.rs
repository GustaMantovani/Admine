#[cfg(not(test))]
use crate::app_context::AppContext;
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
    pub fn new(tags: Vec<String>, message: String) -> Self {
        #[cfg(test)]
        {
            return Self {
                origin: "test".to_string(),
                tags,
                message,
            };
        }

        #[cfg(not(test))]
        {
            Self {
                origin: AppContext::instance().config().self_origin_name().clone(),
                tags,
                message,
            }
        }
    }

    pub fn has_tag(&self, tag: &str) -> bool {
        self.tags.contains(&tag.to_string())
    }
}
