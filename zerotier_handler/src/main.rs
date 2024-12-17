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

    let base_path = env::var("ZEROTIER_API_BASE_URL").expect("ZEROTIER_API_BASE_URL not set");
    let api_key = env::var("ZEROTIER_API_TOKEN").expect("ZEROTIER_API_TOKEN not set");
    let network_id = env::var("ZEROTIER_NETWORK_ID").expect("ZEROTIER_NETWORK_ID not set");
    let record_file_path = env::var("RECORD_FILE_PATH").expect("RECORD_FILE_PATH not set");
    let redis_url = env::var("REDIS_URL").expect("REDIS_URL not set");
    let server_channel = env::var("REDIS_SERVER_CHANNEL").expect("REDIS_SERVER_CHANNEL not set");
    let command_channel = env::var("REDIS_COMMAND_CHANNEL").expect("REDIS_COMMAND_CHANNEL not set");
    let vpn_channel = env::var("REDIS_VPN_CHANNEL").expect("REDIS_VPN_CHANNEL not set");

    // Gera a configuração para o cliente da API
    let config = Configuration::new(base_path.clone(), api_key.clone());

    // Conecta-se ao Redis
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
            let id = admine_message.message.as_str();
            println!("ID: {}", id);

            if let Err(e) =
                handle::remove_old_server_member(&config, &network_id, &record_file_path).await
            {
                println!("Erro ao lidar com membro antigo: {}", e);
            }

            match handle::authorize_new_server_member(&config, &network_id, id, &record_file_path)
                .await
            {
                Ok(ips) => {
                    if !ips.is_empty() {
                        let ip_string = ips
                            .iter()
                            .map(|ip| ip.to_string())
                            .collect::<Vec<String>>()
                            .join(", ");

                        match pub_connection
                            .publish::<&str, &String, ()>(vpn_channel.as_str(), &ip_string)
                        {
                            Ok(_) => println!("IP publicado com sucesso no canal {}", vpn_channel),
                            Err(e) => {
                                println!("Erro ao publicar IP no canal {}: {}", vpn_channel, e)
                            }
                        }
                    }
                }
                Err(e) => println!("Erro ao lidar com novo membro: {}", e),
            }
        }

        // New Member
        if admine_message.tags.contains(&"new_member".to_string()) {
            let id = admine_message.message.as_str();
            println!("ID: {}", id);

            match handle::authorize_member_by_id(&config, &network_id, id).await {
                Ok(member) => {
                    let member_json = match serde_json::to_string(&member) {
                        Ok(json) => json,
                        Err(e) => {
                            println!("Erro ao serializar membro: {}", e);
                            continue;
                        }
                    };

                    match pub_connection
                        .publish::<&str, &String, ()>(vpn_channel.as_str(), &member_json)
                    {
                        Ok(_) => println!("Membro publicado com sucesso no canal {}", vpn_channel),
                        Err(e) => {
                            println!("Erro ao publicar membro no canal {}: {}", vpn_channel, e)
                        }
                    }
                }
                Err(e) => println!("Erro ao autorizar novo membro: {}", e),
            }
        }
    }
}
