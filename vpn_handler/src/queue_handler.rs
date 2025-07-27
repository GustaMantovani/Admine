use crate::app_context::AppContext;
use crate::models::admine_message::AdmineMessage;
use log::{error, info};
use std::sync::Arc;
use tokio::spawn;
use tokio::sync::mpsc;
use tokio::time::sleep;

/// Main handle structure.
pub struct Handle;

impl Handle {
    /// Initializes the Handle using the shared AppContext.
    pub fn new() -> Result<Self, Box<dyn std::error::Error>> {
        // Apenas inicializa o contexto global
        AppContext::instance();

        info!(
            "Handle created successfully with channels: {:?} and retry config: {:?}",
            AppContext::instance().config().admine_channels_map(),
            AppContext::instance().config().retry_config()
        );

        Ok(Self)
    }

    /// Helper function to update the server member ID in the database.
    async fn update_server_id(new_id: &str) -> Result<(), Box<dyn std::error::Error>> {
        AppContext::instance()
            .storage()
            .set("server_member_id".to_string(), new_id.to_string())
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

        // Clone context for the ingestion task.
        spawn(async move {
            while let Some(ingest_message) = rx.recv().await {
                info!("Processing ingestion message: {:?}", ingest_message);
                let member_id = ingest_message.message().clone();

                // Authenticate the member.
                if let Err(e) = AppContext::instance()
                    .vpn_client()
                    .auth_member(member_id.clone(), None)
                    .await
                {
                    error!("Error authenticating member {}: {}", member_id, e);
                    continue;
                }
                info!("Member {} authenticated successfully.", member_id);

                // Retry logic to fetch member IPs until available.
                let mut attempts = *AppContext::instance().config().retry_config().attempts();
                let member_ips = loop {
                    match AppContext::instance()
                        .vpn_client()
                        .get_member_ips_in_vpn(member_id.clone())
                        .await
                    {
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
                                member_id, AppContext::instance().config().retry_config().delay()
                            );
                            sleep(*AppContext::instance().config().retry_config().delay()).await;
                        }
                    }
                };

                // Prepare the message for PubSub publication.
                let pubsub_message = AdmineMessage::new(
                    vec![String::from("new_server_up")],
                    member_ips
                        .iter()
                        .map(|ip| ip.to_string())
                        .collect::<Vec<String>>()
                        .join(","),
                );

                let pubsub_msg_str = match serde_json::to_string(&pubsub_message) {
                    Ok(msg) => msg,
                    Err(e) => {
                        error!("Error serializing message for member {}: {}", member_id, e);
                        continue;
                    }
                };

                info!("Publishing ingestion message: {}", pubsub_msg_str);
                if let Err(e) = AppContext::instance().pub_sub().lock().unwrap().publish(
                    AppContext::instance()
                        .config()
                        .admine_channels_map()
                        .vpn_channel()
                        .clone(),
                    pubsub_msg_str,
                ) {
                    error!("Error publishing ingestion message: {}", e);
                } else {
                    info!("Ingestion message published successfully.");
                }

                // Retrieve the old server member ID from persistence.
                let old_member_id = match AppContext::instance().storage().get("server_member_id") {
                    Some(id) => id,
                    None => {
                        info!("No old server member id found in persistence.");
                        String::new()
                    }
                };

                // If an old ID exists, delete it.
                if !old_member_id.is_empty() {
                    if let Err(e) = AppContext::instance()
                        .vpn_client()
                        .delete_member(old_member_id.clone())
                        .await
                    {
                        error!("Error deleting old member {}: {}", old_member_id, e);
                    }
                }

                // Save the new server member ID in persistence.
                if let Err(e) = Self::update_server_id(&member_id).await {
                    error!("Failed to update server member id: {}", e);
                }
            }
        });

        // Main loop to listen for incoming messages.
        loop {
            info!("Waiting for a new message...");

            let raw_message = match AppContext::instance()
                .pub_sub()
                .lock()
                .unwrap()
                .listen_until_receive_message()
            {
                Ok(msg) => {
                    info!("Message received: {:?}", msg);
                    msg
                }
                Err(e) => {
                    error!("Error receiving message: {}", e);
                    continue;
                }
            };

            let admine_message = match serde_json::from_str::<AdmineMessage>(&raw_message.0) {
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

            let _ = tx.send(Arc::new(admine_message)).await;
        }
    }
}
