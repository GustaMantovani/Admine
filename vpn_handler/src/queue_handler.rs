use crate::app_context::AppContext;
use crate::models::admine_message::AdmineMessage;
use log::{error, info, warn};
use tokio::time::sleep;

/// Main handle structure.
pub struct Handle;

impl Handle {
    /// Initializes the Handle using the shared AppContext.
    pub fn new() -> Result<Self, Box<dyn std::error::Error>> {
        AppContext::instance();
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

    /// Process server_up messages with retry logic and IP fetching
    async fn process_server_up(member_id: String) {
        let vpn_client = AppContext::instance().vpn_client();
        let config = AppContext::instance().config();
        let retry_config = config.retry_config();

        // Authenticate the member
        if let Err(e) = vpn_client.auth_member(member_id.clone(), None).await {
            error!("Error authenticating member {}: {}", member_id, e);
            return;
        }
        info!("Member {} authenticated successfully.", member_id);

        // Retry logic to fetch member IPs until available
        let mut attempts = *retry_config.attempts();
        let member_ips = loop {
            match vpn_client.get_member_ips_in_vpn(member_id.clone()).await {
                Ok(ips) if !ips.is_empty() => break ips,
                Ok(_) | Err(_) => {
                    if attempts == 0 {
                        error!(
                            "Exceeded retry attempts to fetch IPs for member {}",
                            member_id
                        );
                        return;
                    }
                    attempts -= 1;
                    info!(
                        "IPs not available yet for member {}. Retrying in {:?}...",
                        member_id,
                        retry_config.delay()
                    );
                    sleep(*retry_config.delay()).await;
                }
            }
        };

        // Publish new server IPs
        let new_message = AdmineMessage::new(
            vec!["new_server_up".to_string()],
            member_ips
                .iter()
                .map(|ip| ip.to_string())
                .collect::<Vec<String>>()
                .join(","),
        );

        let serialized_message = match serde_json::to_string(&new_message) {
            Ok(json) => json,
            Err(e) => {
                error!("Failed to serialize message: {}", e);
                return;
            }
        };

        if let Err(e) = AppContext::instance().pub_sub().lock().unwrap().publish(
            config.admine_channels_map().vpn_channel().clone(),
            serialized_message,
        ) {
            error!("Failed to publish message: {}", e);
        } else {
            info!("New server up message published successfully.");
        }

        // Handle old server cleanup
        let old_member_id = AppContext::instance()
            .storage()
            .get("server_member_id")
            .unwrap_or_default();

        if !old_member_id.is_empty() && old_member_id != member_id {
            if let Err(e) = vpn_client.delete_member(old_member_id.clone()).await {
                error!("Error deleting old member {}: {}", old_member_id, e);
            }
        }

        // Save the new server member ID
        if let Err(e) = Self::update_server_id(&member_id).await {
            error!("Failed to update server member id: {}", e);
        }
    }

    /// Process auth_member command
    async fn process_auth_member(member_id: String) {
        let vpn_client = AppContext::instance().vpn_client();
        let config = AppContext::instance().config();

        if let Err(e) = vpn_client.auth_member(member_id.clone(), None).await {
            error!("Error authenticating member {}: {}", member_id, e);
            return;
        }
        info!("Member {} authenticated successfully.", member_id);

        // Publish success message
        let success_message =
            AdmineMessage::new(vec!["auth_member_success".to_string()], member_id);

        let serialized_message = match serde_json::to_string(&success_message) {
            Ok(json) => json,
            Err(e) => {
                error!("Failed to serialize success message: {}", e);
                return;
            }
        };

        if let Err(e) = AppContext::instance().pub_sub().lock().unwrap().publish(
            config.admine_channels_map().vpn_channel().clone(),
            serialized_message,
        ) {
            error!("Failed to publish success message: {}", e);
        } else {
            info!("Auth member success message published successfully.");
        }
    }

    /// Process incoming messages based on channel and tags
    async fn process_message(admine_message: AdmineMessage) {
        match admine_message.origin() {
            // Server channel - handle server_up messages
            org if org == "server" => {
                if admine_message.has_tag("server_up") && !admine_message.message().is_empty() {
                    let member_id = admine_message.message().clone();
                    Self::process_server_up(member_id).await;
                }
            }
            // Command channel - handle auth_member commands
            org if org == "bot" => {
                if admine_message.has_tag("auth_member") && !admine_message.message().is_empty() {
                    let member_id = admine_message.message().clone();
                    Self::process_auth_member(member_id).await;
                }
            }
            other => {
                warn!("Unsupported channel: {}", other);
            }
        }
    }

    /// Main run loop.
    pub async fn run(self) {
        info!("Handle run started.");

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

            Self::process_message(admine_message).await;
        }
    }
}
