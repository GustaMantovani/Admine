use crate::utils;
use crate::zerotier::{
    apis::configuration::Configuration,
    apis::network_member_api::{delete_network_member, get_network_member, update_network_member},
    models::member::Member,
};
use log::{error, info, warn};
use std::{net::IpAddr, thread::sleep, time::Duration};

pub async fn remove_old_server_member(
    config: &Configuration,
    network_id: &str,
    record_file_path: &str,
) -> Result<(), Box<dyn std::error::Error>> {
    let old_member_string = utils::read_file(record_file_path.to_string())?;
    if !old_member_string.is_empty() {
        let old_member: Member = serde_json::from_str(&old_member_string)?;
        if let Some(Some(node_id)) = old_member.node_id {
            delete_network_member(config, network_id, &node_id).await?;
            info!("Old member with Node ID {} removed", node_id);
        } else {
            warn!("Node ID not found in old member");
        }
    }
    Ok(())
}

pub async fn authorize_new_server_member(
    config: &Configuration,
    network_id: &str,
    id: &str,
    record_file_path: &str,
    retry_interval: u64,
    retry_count: u64,
) -> Result<Vec<IpAddr>, Box<dyn std::error::Error>> {
    let mut new_member = get_network_member(config, network_id, id).await?;

    if let Some(ref mut config) = new_member.config {
        if let Some(Some(true)) = config.authorized {
            info!("Member already authorized with Node ID: {}", id);
            save_member_to_file(&new_member, record_file_path)?;
            return Ok(new_member.get_member_ips());
        }
    }

    if let Some(mut member_config) = new_member.config {
        member_config.authorized = Some(Some(true));
        new_member.config = Some(member_config);

        if let Some(Some(node_id)) = new_member.node_id.clone() {
            info!("New node ID: {}", node_id);

            let updated_member =
                update_network_member(config, network_id, &node_id, new_member.clone()).await?;
            info!("Member successfully updated with Node ID: {}", node_id);

            if updated_member.config.is_some() {
                sleep(Duration::from_secs(retry_interval));
                let mut retry_counter = 0;
                while retry_counter < retry_count {
                    retry_counter += 1;
                    match get_network_member(config, network_id, &node_id).await {
                        Ok(member) => {
                            if let Some(Some(ip_assignments)) = member
                                .config
                                .as_ref()
                                .and_then(|config| config.ip_assignments.as_ref())
                            {
                                if !ip_assignments.is_empty() {
                                    info!("IP Assignments found for Node ID: {}", node_id);

                                    if let Err(e) = remove_old_server_member(
                                        &config,
                                        &network_id,
                                        &record_file_path,
                                    )
                                    .await
                                    {
                                        error!("Error handling old member: {}", e);
                                    }

                                    save_member_to_file(&updated_member, record_file_path)?;

                                    return Ok(member.get_member_ips());
                                }
                            }
                            warn!("IP Assignments not found for Node ID: {}", node_id);
                            sleep(Duration::from_secs(retry_interval));
                        }
                        Err(e) => {
                            error!(
                                "Error getting updated member with Node ID {}: {}",
                                node_id, e
                            );
                            sleep(Duration::from_secs(retry_interval));
                        }
                    }
                }
            }
        }
    }
    Err("Failed to get member IPs".into())
}

pub fn save_member_to_file(
    member: &Member,
    record_file_path: &str,
) -> Result<(), Box<dyn std::error::Error>> {
    let member_json_string = serde_json::to_string_pretty(member)?;
    utils::write_to_file(record_file_path.to_string(), member_json_string)?;
    info!("Member successfully saved to file {}", record_file_path);
    Ok(())
}

pub async fn authorize_member_by_id(
    config: &Configuration,
    network_id: &str,
    member_id: &str,
) -> Result<Member, Box<dyn std::error::Error>> {
    // Get current member state
    let mut member = get_network_member(config, network_id, member_id).await?;

    // Check if already authorized
    if let Some(ref config) = member.config {
        if let Some(Some(true)) = config.authorized {
            info!("Member already authorized with ID: {}", member_id);
            return Ok(member);
        }
    }

    // Set authorization
    if let Some(mut member_config) = member.config {
        member_config.authorized = Some(Some(true));
        member.config = Some(member_config);

        if let Some(Some(node_id)) = member.node_id.clone() {
            info!("Authorizing member with Node ID: {}", node_id);

            let updated_member =
                update_network_member(config, network_id, &node_id, member.clone()).await?;

            info!("Member successfully authorized with Node ID: {}", node_id);
            return Ok(updated_member);
        }
    }

    Err("Failed to authorize member - invalid configuration".into())
}
