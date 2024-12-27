use log::{error, info};
use models::message::AdmineMessage;
use redis::Commands;
use std::env;
use std::sync::{Arc, Mutex};
use tokio;
use tokio::spawn;
use tokio::sync::mpsc;
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
           \tZEROTIER_API_BASE_URL: {}\n\
           \tZEROTIER_API_TOKEN: [REDACTED]\n\
           \tZEROTIER_NETWORK_ID: {}\n\
           \tZEROTIER_HANDLER_RETRY_COUNT: {}\n\
           \tZEROTIER_HANDLER_RETRY_INTERVAL: {}\n\
           \tRECORD_FILE_PATH: {}\n\
           \tREDIS_URL: {}\n\
           \tREDIS_SERVER_CHANNEL: {}\n\
           \tREDIS_COMMAND_CHANNEL: {}\n\
           \tREDIS_VPN_CHANNEL: {}",
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
    let client = match redis::Client::open(redis_url.clone()) {
        Ok(client) => {
            info!("Connected to Redis at {}", redis_url);
            client
        }
        Err(e) => {
            error!("Failed to connect to Redis at {}: {}", redis_url, e);
            return;
        }
    };

    let mut sub_connection = match client.get_connection() {
        Ok(conn) => {
            info!("Subscription connection established");
            conn
        }
        Err(e) => {
            error!("Failed to establish subscription connection: {}", e);
            return;
        }
    };

    let pub_connection = match client.get_connection() {
        Ok(conn) => {
            info!("Publish connection established");
            Arc::new(Mutex::new(conn))
        }
        Err(e) => {
            error!("Failed to establish publish connection: {}", e);
            return;
        }
    };

    let mut pubsub = sub_connection.as_pubsub();

    if let Err(e) = pubsub.subscribe(&server_channel) {
        error!(
            "Failed to subscribe to server channel {}: {}",
            server_channel, e
        );
        return;
    } else {
        info!("Subscribed to server channel: {}", server_channel);
    }

    if let Err(e) = pubsub.subscribe(&command_channel) {
        error!(
            "Failed to subscribe to command channel {}: {}",
            command_channel, e
        );
        return;
    } else {
        info!("Subscribed to command channel: {}", command_channel);
    }

    let (tx, mut rx) = mpsc::channel::<Arc<AdmineMessage>>(32);

    // Spawn a task to process messages from the queue
    let config_clone = config.clone();
    let network_id_clone = network_id.clone();
    let record_file_path_clone = record_file_path.clone();
    let pub_connection_clone = pub_connection.clone();
    let vpn_channel_clone = vpn_channel.clone();
    spawn(async move {
        while let Some(admine_message) = rx.recv().await {
            let id = admine_message.message.as_str();
            info!("Server starting with ID: {}", id);
            
            match handle::authorize_new_server_member(
                &config_clone,
                &network_id_clone,
                id,
                &record_file_path_clone,
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

                        match pub_connection_clone
                            .lock()
                            .unwrap()
                            .publish::<&str, &String, ()>(vpn_channel_clone.as_str(), &ip_string)
                        {
                            Ok(_) => {
                                info!("IP successfully published to channel {}", vpn_channel_clone)
                            }
                            Err(e) => {
                                error!(
                                    "Error publishing IP to channel {}: {}",
                                    vpn_channel_clone, e
                                )
                            }
                        }
                    }
                }
                Err(e) => error!("Error handling new member: {}", e),
            }
        }
    });

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
            Ok(msg) => Arc::new(msg),
            Err(e) => {
                error!("Error parsing message: {}", e);
                continue;
            }
        };



        // Enqueue the message for processing
        if admine_message.tags.contains(&"server_up".to_string()) {
            let admine_message = admine_message.clone();
            if let Err(e) = tx.send(admine_message).await {
                error!("Error sending message to queue: {}", e);
            }
        }

        if admine_message.tags.contains(&"new_member".to_string()) {
            let id = admine_message.message.clone();
            let config = config.clone();
            let network_id = network_id.clone();
            let vpn_channel = vpn_channel.clone();
            let pub_connection_clone = pub_connection.clone();

            spawn(async move {
                match handle::authorize_member_by_id(&config, &network_id, &id).await {
                    Ok(member) => {
                        let member_json = match serde_json::to_string(&member) {
                            Ok(json) => json,
                            Err(e) => {
                                error!("Error serializing member: {}", e);
                                return;
                            }
                        };

                        match pub_connection_clone
                            .lock()
                            .unwrap()
                            .publish::<&str, &String, ()>(vpn_channel.as_str(), &member_json)
                        {
                            Ok(_) => info!(
                                "Member successfully published to channel {}",
                                vpn_channel
                            ),
                            Err(e) => {
                                error!(
                                    "Error publishing member to channel {}: {}",
                                    vpn_channel, e
                                )
                            }
                        }
                    }
                    Err(e) => error!("Error authorizing new member: {}", e),
                }
            });
        }
    }
}
