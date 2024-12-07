use std::{thread::sleep, time::Duration};
use redis::Commands;

use crate::zerotier::{
    apis::configuration::Configuration,
    apis::network_member_api::{delete_network_member, get_network_member, update_network_member},
    models::member::Member,
};
use crate::utils;

pub async fn remove_old_server_member(
    config: &Configuration, 
    network_id: &str, 
    record_file_path: &str
) -> Result<(), Box<dyn std::error::Error>> {
    let old_member_string = utils::read_file(record_file_path.to_string())?;
    if !old_member_string.is_empty() {
        let old_member: Member = serde_json::from_str(&old_member_string)?;
        if let Some(Some(node_id)) = old_member.node_id {
            delete_network_member(config, network_id, &node_id).await?;
        } else {
            println!("Node ID não encontrado no membro antigo");
        }
    }
    Ok(())
}

pub async fn authorize_new_server_member(
    config: &Configuration,
    network_id: &str,
    id: &str,
    record_file_path: &str,
    pub_connection: &mut redis::Connection,
    vpn_channel: &str,
) -> Result<(), Box<dyn std::error::Error>> {
    let mut new_member = get_network_member(config, network_id, id).await?;
    
    // Verifica se já está autorizado
    if let Some(ref mut config) = new_member.config {
        if let Some(Some(true)) = config.authorized {
            println!("Membro já autorizado");
            save_member_to_file(&new_member, record_file_path)?;
            return Ok(());
        }
    }

    // Autoriza o membro
    if let Some(mut member_config) = new_member.config {
        member_config.authorized = Some(Some(true));
        new_member.config = Some(member_config);

        if let Some(Some(node_id)) = new_member.node_id.clone() {
            println!("Novo node ID: {}", node_id);

            let updated_member = update_network_member(config, network_id, &node_id, new_member.clone()).await?;
            println!("Membro atualizado com sucesso");

            save_member_to_file(&updated_member, record_file_path)?;

            if let Some(_) = updated_member.config {
                sleep(Duration::from_secs(5));
                loop {
                    match get_network_member(config, network_id, &node_id).await {
                        Ok(member) => {
                            if let Some(Some(ip_assignments)) = member.config.as_ref().and_then(|config| config.ip_assignments.as_ref()) {
                                if !ip_assignments.is_empty() {
                                    println!("IP Assignments encontrado");
                                    match pub_connection.publish::<&str, &String, ()>(vpn_channel, &ip_assignments[0]) {
                                        Ok(_) => {
                                            println!("IP publicado com sucesso no canal {}", vpn_channel);
                                            break;
                                        }
                                        Err(e) => {
                                            println!("Erro ao publicar IP no canal {}: {}", vpn_channel, e);
                                            sleep(Duration::from_secs(5));
                                            continue;
                                        }
                                    }
                                }
                            }
                            println!("IP Assignments não encontrado");
                            sleep(Duration::from_secs(5));
                        }
                        Err(e) => {
                            println!("Erro ao obter membro atualizado: {}", e);
                        }
                    }
                }
            } else {
                println!("Configuração não encontrada no membro atualizado");
            }
        } else {
            println!("Node ID não encontrado no novo membro");
        }
    } else {
        println!("Configuração não encontrada no novo membro");
    }
    Ok(())
}

pub fn save_member_to_file(member: &Member, record_file_path: &str) -> Result<(), Box<dyn std::error::Error>> {
    let member_json_string = serde_json::to_string_pretty(member)?;
    utils::write_to_file(record_file_path.to_string(), member_json_string)?;
    println!("Membro salvo com sucesso");
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
            println!("Member already authorized");
            return Ok(member);
        }
    }

    // Set authorization
    if let Some(mut member_config) = member.config {
        member_config.authorized = Some(Some(true));
        member.config = Some(member_config);

        if let Some(Some(node_id)) = member.node_id.clone() {
            println!("Authorizing member with ID: {}", node_id);
            
            let updated_member = update_network_member(
                config,
                network_id, 
                &node_id,
                member.clone()
            ).await?;
            
            println!("Member authorized successfully");
            return Ok(updated_member);
        }
    }

    Err("Failed to authorize member - invalid configuration".into())
}