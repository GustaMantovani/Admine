use super::pub_sub::{TPublisher, TSubscriber};
use crate::errors::PubSubError;
use redis::Commands;

pub struct RedisPubSub {
    client: redis::Client,
    connection: redis::Connection,
    subscribed_topics: Vec<String>,
}

impl RedisPubSub {
    pub fn new(url: &str) -> Self {
        let client = redis::Client::open(url).unwrap();
        let connection = client.get_connection().unwrap();
        let subscribed_topics = Vec::new();
        Self {
            client,
            connection,
            subscribed_topics,
        }
    }
}

impl TSubscriber for RedisPubSub {
    fn subscribe(&mut self, topics: Vec<String>) -> Result<(), PubSubError> {
        topics.iter().for_each(|t| {
            self.subscribed_topics.push(t.clone());
        });
        Ok(())
    }

    fn unsubscribe(&mut self, topic: String) -> Result<(), PubSubError> {
        self.subscribed_topics.retain(|t| t != &topic);
        Ok(())
    }

    fn listen_until_to_ricieve_message(&mut self) -> Result<(String, String), PubSubError> {
        let mut pubsub = self.connection.as_pubsub();

        self.subscribed_topics.iter().for_each(|t| {
            pubsub.subscribe(t).unwrap();
        });

        let msg = pubsub.get_message().unwrap();

        match msg.get_payload() {
            Ok(p) => return Ok((p, msg.get_channel_name().to_string())),
            Err(e) => return Err(PubSubError::MessageError(e.to_string())),
        };
    }
}

impl TPublisher for RedisPubSub {
    fn publish(&mut self, topic: String, message: String) -> Result<(), PubSubError> {
        self.connection
            .publish::<String, String, ()>(topic, message)
            .unwrap();
        Ok(())
    }
}
