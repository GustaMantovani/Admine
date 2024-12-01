use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, PartialEq, Serialize, Deserialize)]
pub enum MemberConfigTagsInnerInner {
    Variant0(i64),
    Variant1(bool),
}
