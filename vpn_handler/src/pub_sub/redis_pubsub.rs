use super::pub_sub::{TPublisher, TSubscriber};
use crate::errors::PubSubError;
use redis::Commands;

pub struct RedisPubSub {
    connection: redis::Connection,
    subscribed_topics: Vec<String>,
}

impl RedisPubSub {
    pub fn new(url: &str) -> Result<Self, PubSubError> {
        let client = redis::Client::open(url).map_err(|e| {
            log::error!("Failed to create Redis client with URL '{}': {}", url, e);
            PubSubError::CreationError(format!("Failed to create Redis client: {}", e))
        })?;
        
        let connection = client.get_connection().map_err(|e| {
            log::error!("Failed to establish Redis connection: {}", e);
            PubSubError::ConnectionError(format!("Failed to connect to Redis: {}", e))
        })?;
        
        let subscribed_topics = Vec::new();
        Ok(Self {
            connection,
            subscribed_topics,
        })
    }
}

impl TSubscriber for RedisPubSub {
    fn subscribe(&mut self, topics: Vec<String>) -> Result<(), PubSubError> {
        topics.iter().for_each(|t| {
            self.subscribed_topics.push(t.clone());
        });
        Ok(())
    }

    fn listen_until_receive_message(&mut self) -> Result<(String, String), PubSubError> {
        let mut pubsub = self.connection.as_pubsub();

        self.subscribed_topics.iter().for_each(|t| {
            if let Err(e) = pubsub.subscribe(t) {
                log::error!("Failed to subscribe to topic '{}': {}", t, e);
            }
        });

        let msg = pubsub.get_message().map_err(|e| {
            log::error!("Failed to receive message from Redis: {}", e);
            PubSubError::MessageError(format!("Failed to get message: {}", e))
        })?;

        match msg.get_payload() {
            Ok(p) => return Ok((p, msg.get_channel_name().to_string())),
            Err(e) => {
                log::error!("Failed to parse message payload: {}", e);
                return Err(PubSubError::MessageError(e.to_string()));
            }
        };
    }
}

impl TPublisher for RedisPubSub {
    fn publish(&mut self, topic: String, message: String) -> Result<(), PubSubError> {
        self.connection
            .publish::<String, String, ()>(topic.clone(), message.clone())
            .map_err(|e| {
                log::error!("Failed to publish message to topic '{}': {}", topic, e);
                PubSubError::MessageError(format!("Failed to publish message: {}", e))
            })?;
        Ok(())
    }
}
