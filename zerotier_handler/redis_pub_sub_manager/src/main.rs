use dotenvy::dotenv;
use std::env;
use zerotier_api_client::apis::network_member_api::get_network_member_list;
use zerotier_api_client::apis::configuration::Configuration;

#[tokio::main]
async fn main() {
    dotenv().ok();

    let base_path = env::var("ZEROTIER_API_BASE_URL").expect("ZEROTIER_API_BASE_URL not set");
    let api_key = env::var("ZEROTIER_API_TOKEN").expect("ZEROTIER_API_TOKEN not set");
    let network_id = env::var("ZEROTIER_NETWORK_ID").expect("ZEROTIER_NETWORK_ID not set");

    let config = Configuration::new(base_path, api_key);

    println!("{:?}", config);

    match get_network_member_list(&config, &network_id).await {
        Ok(response) => {
            for member in response {
                println!("{:?}", member);
            }
        },
        Err(e) => eprintln!("Error: {:?}", e),
    }
}
