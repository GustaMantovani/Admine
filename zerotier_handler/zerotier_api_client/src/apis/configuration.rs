#[derive(Debug, Clone)]
pub struct Configuration {
    pub base_path: String,
    pub user_agent: Option<String>,
    pub client: reqwest::Client,
    pub basic_auth: Option<BasicAuth>,
    pub oauth_access_token: Option<String>,
    pub bearer_access_token: Option<String>,
    pub api_key: ApiKey,
}

pub type BasicAuth = (String, Option<String>);

#[derive(Debug, Clone)]
pub struct ApiKey {
    pub prefix: Option<String>,
    pub key: String,
}

impl ApiKey{
    pub fn new(api_key: String) -> ApiKey {
        ApiKey { prefix: None, key: api_key }
    }
}

impl Configuration {
    pub fn new(base_path: String, api_key: String) -> Configuration {
        Configuration {
            base_path: base_path,
            user_agent: Some("Rust Zerotier Handler v1".to_owned()),
            client: reqwest::Client::new(),
            basic_auth: None,
            oauth_access_token: None,
            bearer_access_token: None,
            api_key: ApiKey::new(api_key),
        }
    }
}
