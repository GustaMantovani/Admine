use crate::models::admine_message::AdmineMessage;
use crate::persistence::{
    factories::{StoreFactory, StoreType},
    key_value_store::KeyValueStore,
};
use crate::pub_sub::factories::{PubSubFactory, PubSubType};
use crate::pub_sub::pub_sub::PubSubProvider;
use crate::vpn::{
    factories::{VpnFactory, VpnType},
    vpn::TVpnClient,
};
use dotenvy::dotenv;
use log::{error, info, warn};
use std::env;
use std::fmt;
use std::str::FromStr;
use std::sync::Arc;
use tokio::spawn;
use tokio::sync::{mpsc, Mutex};
use tokio::time::{sleep, Duration};

/// Structure holding channel names.
#[derive(Debug, Clone)]
struct AdmineChannelsMap {
    server_channel: String,
    command_channel: String,
    vpn_channel: String,
}

impl fmt::Display for AdmineChannelsMap {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(
            f,
            "AdmineChannelsMap {{ server_channel: {}, command_channel: {}, vpn_channel: {} }}",
            self.server_channel, self.command_channel, self.vpn_channel
        )
    }
}

/// Configuration for retry logic.
#[derive(Debug, Clone)]
pub struct RetryConfig {
    pub attempts: usize,
    pub delay: Duration,
}

/// Main handle structure.
pub struct Handle {
    pub_sub_listener: Arc<Mutex<Box<dyn PubSubProvider>>>,
    pub_sub_publisher: Arc<Mutex<Box<dyn PubSubProvider>>>,
    vpn: Arc<Box<dyn TVpnClient + Send + Sync>>,
    db: Box<dyn KeyValueStore + Send + Sync>,
    admine_channels_map: AdmineChannelsMap,
    retry_config: RetryConfig,
}

impl Handle {
    /// Initializes the Handle by loading configuration from environment variables.
    pub fn new() -> Result<Self, Box<dyn std::error::Error>> {
        dotenv().ok();

        fn fetch_env_var(var_name: &str) -> Result<String, Box<dyn std::error::Error>> {
            env::var(var_name).map_err(|_| {
                let msg = format!("Missing environment variable: {}", var_name);
                error!("{}", msg);
                msg.into()
            })
        }

        // Load configuration values.
        let pubsub_url = fetch_env_var("PUBSUB_URL")?;
        let pubsub_type = fetch_env_var("PUBSUB_TYPE")?;
        let api_url = fetch_env_var("VPN_API_URL")?;
        let api_key = fetch_env_var("VPN_API_KEY")?;
        let network_id = fetch_env_var("VPN_NETWORK_ID")?;
        let server_channel = fetch_env_var("SERVER_CHANNEL")?;
        let command_channel = fetch_env_var("COMMAND_CHANNEL")?;
        let vpn_channel = fetch_env_var("VPN_CHANNEL")?;
        let db_path = fetch_env_var("DB_PATH")?;
        let store_type = fetch_env_var("STORE_TYPE")?;
        let retry_attempts = fetch_env_var("VPN_RETRY_ATTEMPTS")?;
        let retry_delay_ms = fetch_env_var("VPN_RETRY_DELAY_MS")?;

        // Create VPN client using cloned values to avoid move errors.
        let vpn = VpnFactory::create_vpn(
            VpnType::Zerotier,
            api_url.clone(),
            api_key.clone(),
            network_id.clone(),
        )
        .map_err(|e| {
            error!(
                "Error creating VPN client with API URL: {}, API Key: {}, Network ID: {}: {}",
                api_url.clone(),
                api_key.clone(),
                network_id.clone(),
                e
            );
            e
        })?;

        // Parse PubSub type.
        let pub_sub_type = PubSubType::from_str(&pubsub_type).map_err(|_| {
            error!("Unsupported PubSub type: {}", pubsub_type);
            "Unsupported PubSub type"
        })?;

        // Create publisher and listener instances.
        let pub_sub_publisher =
            PubSubFactory::create_pubsub_instance(pub_sub_type.clone(), &pubsub_url).map_err(
                |e| {
                    error!(
                        "Error creating PubSub publisher with type: {}, Redis URL: {}: {}",
                        pubsub_type, pubsub_url, e
                    );
                    e
                },
            )?;
        let mut pub_sub_listener = PubSubFactory::create_pubsub_instance(pub_sub_type, &pubsub_url)
            .map_err(|e| {
                error!(
                    "Error creating PubSub listener with Redis URL: {}: {}",
                    pubsub_url, e
                );
                e
            })?;

        // Subscribe listener to channels.
        pub_sub_listener.subscribe(vec![server_channel.clone(), command_channel.clone()])?;

        // Parse store type and create DB instance.
        let store_type_enum = StoreType::from_str(&store_type).map_err(|_| {
            error!("Unsupported Store type: {}", store_type);
            "Unsupported Store type"
        })?;
        let db = StoreFactory::create_store_instance(store_type_enum, &db_path).map_err(|e| {
            error!("Error creating store instance: {}", e);
            e
        })?;

        let admine_channels_map = AdmineChannelsMap {
            server_channel,
            command_channel,
            vpn_channel,
        };

        let retry_config = RetryConfig {
            attempts: retry_attempts.parse()?,
            delay: Duration::from_millis(retry_delay_ms.parse()?),
        };

        info!(
            "Handle created successfully with channels: {:?} and retry config: {:?}",
            admine_channels_map, retry_config
        );

        Ok(Self {
            pub_sub_listener: Arc::new(Mutex::new(pub_sub_listener)),
            pub_sub_publisher: Arc::new(Mutex::new(pub_sub_publisher)),
            vpn: Arc::new(vpn),
            db,
            admine_channels_map,
            retry_config,
        })
    }

    /// Helper function to update the server member ID in the database.
    /// Note: Takes a mutable reference to the DB because set() requires mutable access.
    async fn update_server_id(
        db: &mut Box<dyn KeyValueStore + Send + Sync>,
        new_id: &str,
    ) -> Result<(), Box<dyn std::error::Error>> {
        db.set("server_member_id".to_string(), new_id.to_string())
            .map_err(|e| {
                error!("Error saving new server member id: {}", e);
                e.into()
            })
    }

    /// Main run loop.
    pub async fn run(self) -> Result<(), Box<dyn std::error::Error>> {
        info!("Handle run started.");

        // Create an ingestion queue for server messages.
        info!("Creating ingestion queue to handle server messages");
        let (tx, mut rx) = mpsc::channel::<Arc<AdmineMessage>>(32);

        info!("Spawning task to process ingestion messages");

        // Clone fields needed for the ingestion task.
        let ingestion_vpn = Arc::clone(&self.vpn);
        let ingestion_pubsub = Arc::clone(&self.pub_sub_publisher);
        let ingestion_vpn_channel = self.admine_channels_map.vpn_channel.clone();
        let ingestion_retry_config = self.retry_config.clone();
        // Make the DB mutable so we can perform update operations.
        let mut ingestion_db = self.db;

        spawn(async move {
            while let Some(ingest_message) = rx.recv().await {
                info!("Processing ingestion message: {:?}", ingest_message);
                let member_id = ingest_message.message.clone();

                // Authenticate the member.
                if let Err(e) = ingestion_vpn.auth_member(member_id.clone(), None).await {
                    error!("Error authenticating member {}: {}", member_id, e);
                    continue;
                }
                info!("Member {} authenticated successfully.", member_id);

                // Retry logic to fetch member IPs until available.
                let mut attempts = ingestion_retry_config.attempts;
                let member_ips = loop {
                    match ingestion_vpn.get_member_ips_in_vpn(member_id.clone()).await {
                        Ok(ips) if !ips.is_empty() => break ips,
                        Ok(_) | Err(_) => {
                            if attempts == 0 {
                                error!(
                                    "Exceeded retry attempts to fetch IPs for member {}",
                                    member_id
                                );
                                break Vec::new();
                            }
                            attempts -= 1;
                            info!(
                                "IPs not available yet for member {}. Retrying in {:?}...",
                                member_id, ingestion_retry_config.delay
                            );
                            sleep(ingestion_retry_config.delay).await;
                        }
                    }
                };

                // Prepare the message for PubSub publication.
                let pubsub_message = AdmineMessage {
                    tags: vec![String::from("new_server_up")],
                    message: member_ips
                        .iter()
                        .map(|ip| ip.to_string())
                        .collect::<Vec<String>>()
                        .join(","),
                };

                let pubsub_msg_str = match pubsub_message.to_json_string() {
                    Ok(msg) => msg,
                    Err(e) => {
                        error!("Error serializing message for member {}: {}", member_id, e);
                        continue;
                    }
                };

                {
                    let mut publisher = ingestion_pubsub.lock().await;
                    info!("Publishing ingestion message: {}", pubsub_msg_str);
                    if let Err(e) = publisher.publish(ingestion_vpn_channel.clone(), pubsub_msg_str)
                    {
                        error!("Error publishing ingestion message: {}", e);
                    } else {
                        info!("Ingestion message published successfully.");
                    }
                }

                // Retrieve the old server member ID from persistence.
                let old_member_id = match ingestion_db.get("server_member_id") {
                    Ok(Some(id)) => id,
                    Ok(None) => {
                        info!("No old server member id found in persistence.");
                        String::new()
                    }
                    Err(e) => {
                        error!("Error retrieving old server member id: {}", e);
                        String::new()
                    }
                };

                // If an old ID exists, delete it.
                if !old_member_id.is_empty() {
                    if let Err(e) = ingestion_vpn.delete_member(old_member_id.clone()).await {
                        error!("Error deleting old member {}: {}", old_member_id, e);
                    }
                }

                // Save the new server member ID in persistence.
                if let Err(e) = Self::update_server_id(&mut ingestion_db, &member_id).await {
                    error!("Failed to update server member id: {}", e);
                }
            }
        });

        // Main loop to listen for incoming messages.
        loop {
            info!("Waiting for a new message...");

            let raw_message = {
                let mut listener = self.pub_sub_listener.lock().await;
                match listener.listen_until_to_ricieve_message() {
                    Ok(msg) => {
                        info!("Message received: {:?}", msg);
                        msg
                    }
                    Err(e) => {
                        error!("Error receiving message: {}", e);
                        continue;
                    }
                }
            };

            let admine_message = match AdmineMessage::from_json_string(&raw_message.0) {
                Ok(msg) => msg,
                Err(e) => {
                    error!("Error deserializing message: {}", e);
                    continue;
                }
            };

            info!(
                "Processing message received on channel {}: {:?}",
                raw_message.1, admine_message
            );

            match raw_message.1.as_str() {
                // For messages from the server channel with "server_up" tag, enqueue for ingestion.
                s if s == self.admine_channels_map.server_channel => {
                    if admine_message.tags.contains(&"server_up".to_string()) {
                        if let Err(e) = tx.send(Arc::new(admine_message)).await {
                            error!("Error sending message to ingestion queue: {}", e);
                        }
                    }
                }
                // For messages from the command channel with "auth_member" tag, process in a separate task.
                s if s == self.admine_channels_map.command_channel => {
                    if admine_message.tags.contains(&"auth_member".to_string()) {
                        // Clone fields for the command task.
                        let command_vpn = Arc::clone(&self.vpn);
                        let command_pubsub = Arc::clone(&self.pub_sub_publisher);
                        let command_vpn_channel = self.admine_channels_map.vpn_channel.clone();
                        let member_id = admine_message.message.clone();

                        tokio::spawn(async move {
                            if let Err(e) = command_vpn.auth_member(member_id.clone(), None).await {
                                error!("Error authenticating member {}: {}", member_id, e);
                                return;
                            }
                            info!("Member {} authenticated successfully.", member_id);

                            let command_message = AdmineMessage {
                                tags: vec![String::from("auth_member_success")],
                                message: member_id.clone(),
                            };

                            let command_pubsub_msg = match command_message.to_json_string() {
                                Ok(msg) => msg,
                                Err(e) => {
                                    error!(
                                        "Error serializing message for member {}: {}",
                                        member_id, e
                                    );
                                    return;
                                }
                            };

                            let mut publisher = command_pubsub.lock().await;
                            info!("Publishing command channel message: {}", command_pubsub_msg);
                            if let Err(e) =
                                publisher.publish(command_vpn_channel.clone(), command_pubsub_msg)
                            {
                                error!("Error publishing command channel message: {}", e);
                            } else {
                                info!("Command channel message published successfully.");
                            }
                        });
                    }
                }
                other => {
                    warn!("Unsupported channel: {}", other);
                }
            }
        }
    }
}
