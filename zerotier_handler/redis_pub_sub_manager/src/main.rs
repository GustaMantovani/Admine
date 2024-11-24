use std::{env, thread::sleep, time::Duration};
use tokio;
use zerotier_api_client::{
    apis::configuration::Configuration,
    apis::network_member_api::{delete_network_member, get_network_member, update_network_member},
    models::member::Member,
};

mod models;
mod utils;

#[tokio::main]
async fn main() {
    dotenvy::dotenv().ok();

    let base_path = env::var("ZEROTIER_API_BASE_URL").expect("ZEROTIER_API_BASE_URL not set");
    let api_key = env::var("ZEROTIER_API_TOKEN").expect("ZEROTIER_API_TOKEN not set");
    let network_id = env::var("ZEROTIER_NETWORK_ID").expect("ZEROTIER_NETWORK_ID not set");
    let record_file_path = env::var("RECORD_FILE_PATH").expect("RECORD_FILE_PATH not set");

    // Gera a configuração para o cliente da API
    let config = Configuration::new(base_path.clone(), api_key.clone());

    // Conecta-se ao Redis
    let client = redis::Client::open("redis://127.0.0.1/").unwrap();
    let mut connection = client.get_connection().unwrap();
    let mut pubsub = connection.as_pubsub();

    pubsub.subscribe("server_channel").unwrap();

    loop {
        let msg = match pubsub.get_message() {
            Ok(msg) => msg,
            Err(e) => {
                println!("Erro ao receber mensagem: {}", e);
                continue;
            }
        };

        let payload: String = match msg.get_payload() {
            Ok(payload) => payload,
            Err(e) => {
                println!("Erro ao obter payload: {}", e);
                continue;
            }
        };

        // Parse AdmineMessage
        let admine_message = match serde_json::from_str::<models::message::AdmineMessage>(&payload)
        {
            Ok(msg) => msg,
            Err(e) => {
                println!("Erro ao parsear mensagem: {}", e);
                continue;
            }
        };

        // Server Up
        if admine_message.tags.contains(&"server_up".to_string()) {
            let parts: Vec<&str> = admine_message.message.split_whitespace().collect();
            if parts.len() <= 2 {
                println!("Formato da mensagem inválido");
                continue;
            }

            let id = parts[2];
            println!("ID: {}", id);

            // Read old member
            match utils::read_file(record_file_path.clone()) {
                Ok(old_member_string) => {
                    if !old_member_string.is_empty() {
                        match serde_json::from_str::<Member>(&old_member_string) {
                            Ok(old_member) => {
                                if let Some(Some(node_id)) = old_member.node_id {
                                    if let Err(e) =
                                        delete_network_member(&config, &network_id, &node_id).await
                                    {
                                        println!("Erro ao remover membro antigo: {}", e);
                                    }
                                } else {
                                    println!("Node ID não encontrado no membro antigo");
                                }
                            }
                            Err(e) => println!("Erro ao desserializar membro antigo: {}", e),
                        }
                    }
                }
                Err(e) => println!("Erro ao ler arquivo de registro: {}", e),
            }

            // Get and update new member
            match get_network_member(&config, &network_id, id).await {
                Ok(mut new_member) => {
                    match new_member.config {
                        Some(ref mut config) => {
                            if let Some(Some(true)) = config.authorized {
                                println!("Membro já autorizado");
                                continue;
                            }
                        }
                        None => {
                            println!("Configuração não encontrada no novo membro");
                            continue;
                        }
                    }

                    if let Some(mut member_config) = new_member.config {
                        member_config.authorized = Some(Some(true));
                        new_member.config = Some(member_config);

                        if let Some(Some(node_id)) = new_member.node_id.clone() {
                            println!("Novo node ID: {}", node_id);

                            match update_network_member(
                                &config,
                                &network_id,
                                &node_id,
                                new_member.clone(),
                            )
                            .await
                            {
                                Ok(updated_member) => {
                                    println!("Membro atualizado com sucesso");

                                    // Serializa o membro atualizado para JSON
                                    match serde_json::to_string(&updated_member) {
                                        Ok(json) => {
                                            // Salva no arquivo
                                            if let Err(e) =
                                                utils::write_to_file(record_file_path.clone(), json)
                                            {
                                                println!("Erro ao salvar membro atualizado: {}", e);
                                            } else {
                                                println!("Membro atualizado salvo com sucesso");
                                            }

                                            match updated_member.config {
                                                Some(ref member_config) => {
                                                    
                                                    loop {
                                                        
                                                        match get_network_member(&config.clone(), network_id.as_str(), node_id.as_str()).await {
                                                            Ok(member) => {
                                                                match member.config {
                                                                    Some(ref config) => {
                                                                        match config.ip_assignments {
                                                                            Some(ref ip_assignments) => {
                                                                                if !ip_assignments.is_none() {
                                                                                    println!("IP Assignments encontrado");
                                                                                    // Publicar o novo ip no canal do bot
                                                                                    break;
                                                                                } else {
                                                                                    println!("IP Assignments não encontrado");
                                                                                    sleep(Duration::from_secs(5));
                                                                                    continue;
                                                                                }
                                                                            }
                                                                            None => {
                                                                                println!("IP Assignments não encontrado");
                                                                            }
                                                                            
                                                                        }
                                                                    }
                                                                    None => {
                                                                        println!("Configuração não encontrada no membro atualizado");
                                                                        continue;
                                                                    }
                                                                }
                                                            }
                                                            Err(e) => {
                                                                println!("Erro ao obter membro atualizado: {}", e);
                                                                continue;
                                                            }
                                                            
                                                        }
                                                    }

                                                }
                                                None => {
                                                    println!("Configuração não encontrada no membro atualizado");
                                                    continue;
                                                }
                                            }
                                            
                                        }
                                        Err(e) => {
                                            println!("Erro ao serializar membro atualizado: {}", e)
                                        }
                                    }
                                }
                                Err(e) => println!("Erro ao atualizar novo membro: {}", e),
                            }
                        } else {
                            println!("Node ID não encontrado no novo membro");
                        }
                    } else {
                        println!("Configuração não encontrada no novo membro");
                    }
                }
                Err(e) => println!("Erro ao obter novo membro: {}", e),
            }
        }
    }
}
