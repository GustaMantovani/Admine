use serde::Serialize;

#[derive(Serialize)]
pub struct ErrorResponse {
    pub message: String,
}
