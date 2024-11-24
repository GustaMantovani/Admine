use tokio;
use zerotier_api_client::{
    apis::configuration::Configuration,
    models::member::Member,
    apis::network_member_api::{
        delete_network_member,
        get_network_member,
        update_network_member,
        get_network_member_list
    },
};
use std::env;

mod utils;
mod models;

#[tokio::main]
async fn main() {

    dotenvy::dotenv().ok();

    // Conecta-se ao Redis
    let client = redis::Client::open("redis://127.0.0.1/").unwrap();
    let mut connection = client.get_connection().unwrap();
    let mut pubsub = connection.as_pubsub();

    pubsub.subscribe("server_channel").unwrap();

    loop {
        let msg = pubsub.get_message().unwrap();
        let payload: String = msg.get_payload().unwrap();

        // println!("Channel '{}': {}", msg.get_channel_name(), payload);

        // Tentar parsear a mensagem no formato padrão do Admine
        match serde_json::from_str::<models::message::AdmineMessage>(&payload) {
            Ok(admine_message) => {
                // println!("AdmineMessage: {:?}", admine_message);

                if admine_message.tags.contains(&"server_up".to_string()) {

                    // Extrai o ID da mensagem do comando zerotier-cli info
                    let parts: Vec<&str> = admine_message.message.split_whitespace().collect();
                    if parts.len() > 2 {

                        let id = parts[2];

                        print!("{}\n", id);

                        let base_path = env::var("ZEROTIER_API_BASE_URL").expect("ZEROTIER_API_BASE_URL not set");
                        let api_key = env::var("ZEROTIER_API_TOKEN").expect("ZEROTIER_API_TOKEN not set");
                        let network_id = env::var("ZEROTIER_NETWORK_ID").expect("ZEROTIER_NETWORK_ID not set");
                        let record_file_path = env::var("RECORD_FILE_PATH").expect("RECORD_FILE_PATH not set");
                    
                        // Load the configuration
                        let config = Configuration::new(base_path, api_key);
    
                        // Read the old member from json file
                        let old_member_as_string = utils::read_file(record_file_path).unwrap_or_default();
    
                        if !old_member_as_string.is_empty() {
                            let old_member: Member = serde_json::from_str(&old_member_as_string).unwrap();
    
                            // Removing old member from zerotier network
                            delete_network_member(
                                &config,
                                network_id.as_str(),
                                old_member
                                    .node_id
                                    .unwrap()
                                    .expect("Node ID not found")
                                    .as_str(),
                            )
                            .await
                            .unwrap();
                        }

                        // Get new server member from zerotier network

                        let new_member = get_network_member(&config, network_id.as_str(), id).await.unwrap();

                        // Update new member authorization
                        let mut new_member_config = new_member.config.clone().unwrap();

                        new_member_config.authorized = Some(Some((true)));

                        let new_member_updatted = Member{
                            config: Some(new_member_config),
                            ..new_member
                        };

                        print!("{}\n", new_member_updatted.clone().node_id.unwrap().unwrap());

                        // Update new member in zerotier network
                        update_network_member(&config, network_id.as_str(), new_member_updatted.clone().node_id.unwrap().unwrap().as_str(), new_member_updatted.clone()).await.unwrap();

                        

                    } else {
                        println!("Formato da mensagem inválido");
                    }
                    

                }
            }
            Err(e) => {
                print!("Mensagem não é um AdmineMessage {}", e);
            }
        }
    }

}
