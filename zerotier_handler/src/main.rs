use log::{error, info};
use redis::Commands;
use std::env;
use tokio;
use zerotier::apis::configuration::Configuration;
mod handle;
mod models;
mod utils;
mod zerotier;

#[tokio::main]
async fn main() {
    dotenvy::dotenv().ok();
    log4rs::init_file("log4rs.yaml", Default::default()).unwrap();

    let base_path = env::var("ZEROTIER_API_BASE_URL").expect("ZEROTIER_API_BASE_URL not set");
    let api_key = env::var("ZEROTIER_API_TOKEN").expect("ZEROTIER_API_TOKEN not set");
    let network_id = env::var("ZEROTIER_NETWORK_ID").expect("ZEROTIER_NETWORK_ID not set");
    let retry_count: u64 = env::var("ZEROTIER_HANDLER_RETRY_COUNT")
        .expect("ZEROTIER_HANDLER_RETRY_COUNT not set")
        .parse()
        .expect("Failed to parse ZEROTIER_HANDLER_RETRY_COUNT as u64");

    let retry_interval: u64 = env::var("ZEROTIER_HANDLER_RETRY_INTERVAL")
        .expect("ZEROTIER_HANDLER_RETRY_INTERVAL not set")
        .parse()
        .expect("Failed to parse ZEROTIER_HANDLER_RETRY_INTERVAL as u64");
    let record_file_path = env::var("RECORD_FILE_PATH").expect("RECORD_FILE_PATH not set");
    let redis_url = env::var("REDIS_URL").expect("REDIS_URL not set");
    let server_channel = env::var("REDIS_SERVER_CHANNEL").expect("REDIS_SERVER_CHANNEL not set");
    let command_channel = env::var("REDIS_COMMAND_CHANNEL").expect("REDIS_COMMAND_CHANNEL not set");
    let vpn_channel = env::var("REDIS_VPN_CHANNEL").expect("REDIS_VPN_CHANNEL not set");

    info!(
        "Starting application with the following environment variables:\n\
           ZEROTIER_API_BASE_URL: {}\n\
           ZEROTIER_API_TOKEN: [REDACTED]\n\
           ZEROTIER_NETWORK_ID: {}\n\
           ZEROTIER_HANDLER_RETRY_COUNT: {}\n\
           ZEROTIER_HANDLER_RETRY_INTERVAL: {}\n\
           RECORD_FILE_PATH: {}\n\
           REDIS_URL: {}\n\
           REDIS_SERVER_CHANNEL: {}\n\
           REDIS_COMMAND_CHANNEL: {}\n\
           REDIS_VPN_CHANNEL: {}",
        base_path,
        network_id,
        retry_count,
        retry_interval,
        record_file_path,
        redis_url,
        server_channel,
        command_channel,
        vpn_channel
    );

    // Generate configuration for the API client
    let config = Configuration::new(base_path.clone(), api_key.clone());

    // Connect to Redis
    let client = redis::Client::open(redis_url.clone()).unwrap();
    let mut sub_connection = client.get_connection().unwrap();
    let mut pub_connection = client.get_connection().unwrap();
    let mut pubsub = sub_connection.as_pubsub();

    pubsub.subscribe(&server_channel).unwrap();
    pubsub.subscribe(&command_channel).unwrap();

    loop {
        let msg = match pubsub.get_message() {
            Ok(msg) => msg,
            Err(e) => {
                error!("Error receiving message: {}", e);
                continue;
            }
        };

        let payload: String = match msg.get_payload() {
            Ok(payload) => payload,
            Err(e) => {
                error!("Error getting payload: {}", e);
                continue;
            }
        };

        // Parse AdmineMessage
        let admine_message = match serde_json::from_str::<models::message::AdmineMessage>(&payload)
        {
            Ok(msg) => msg,
            Err(e) => {
                error!("Error parsing message: {}", e);
                continue;
            }
        };

        // Server Up
        if admine_message.tags.contains(&"server_up".to_string()) {
            let id = admine_message.message.as_str();
            info!("Server starting with ID: {}", id);

            if let Err(e) =
                handle::remove_old_server_member(&config, &network_id, &record_file_path).await
            {
                error!("Error handling old member: {}", e);
            }

            match handle::authorize_new_server_member(
                &config,
                &network_id,
                id,
                &record_file_path,
                retry_interval,
                retry_count,
            )
            .await
            {
                Ok(ips) => {
                    if !ips.is_empty() {
                        let ip_string = ips
                            .iter()
                            .map(|ip| ip.to_string())
                            .collect::<Vec<String>>()
                            .join(", ");
                        info!("Authorized IPs: {}", ip_string);

                        match pub_connection
                            .publish::<&str, &String, ()>(vpn_channel.as_str(), &ip_string)
                        {
                            Ok(_) => info!("IP successfully published to channel {}", vpn_channel),
                            Err(e) => {
                                error!("Error publishing IP to channel {}: {}", vpn_channel, e)
                            }
                        }
                    }
                }
                Err(e) => error!("Error handling new member: {}", e),
            }
        }

        // New Member
        if admine_message.tags.contains(&"new_member".to_string()) {
            let id = admine_message.message.as_str();
            info!("New member with ID: {}", id);

            match handle::authorize_member_by_id(&config, &network_id, id).await {
                Ok(member) => {
                    let member_json = match serde_json::to_string(&member) {
                        Ok(json) => json,
                        Err(e) => {
                            error!("Error serializing member: {}", e);
                            continue;
                        }
                    };

                    match pub_connection
                        .publish::<&str, &String, ()>(vpn_channel.as_str(), &member_json)
                    {
                        Ok(_) => info!("Member successfully published to channel {}", vpn_channel),
                        Err(e) => {
                            error!("Error publishing member to channel {}: {}", vpn_channel, e)
                        }
                    }
                }
                Err(e) => error!("Error authorizing new member: {}", e),
            }
        }
    }
}
