use serde::Deserialize;

#[derive(Deserialize)]
pub struct AuthMemberRequest {
    pub member_id: String,
}
