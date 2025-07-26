use crate::app_context::AppContext;
use getset::{Getters, MutGetters, Setters};
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, Clone, Getters, Setters, MutGetters)]
pub struct AdmineMessage {
    #[getset(get = "pub", set = "pub")]
    origin: String,

    #[getset(get = "pub", set = "pub", get_mut = "pub")]
    tags: Vec<String>,

    #[getset(get = "pub", set = "pub")]
    message: String,
}

impl AdmineMessage {
    pub fn new(tags: Vec<String>, message: String) -> Self {
        Self {
            origin: AppContext::instance().config().self_origin_name.clone(),
            tags,
            message,
        }
    }

    pub fn has_tag(&self, tag: &str) -> bool {
        self.tags.contains(&tag.to_string())
    }
}
