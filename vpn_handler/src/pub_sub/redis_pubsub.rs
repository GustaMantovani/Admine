use super::pub_sub::{TPublisher, TSubscriber};
use crate::errors::PubSubError;
use async_trait::async_trait;
use futures_util::StreamExt;

pub struct RedisPubSub {
    subscriber: redis::aio::PubSub,
    publisher: redis::aio::MultiplexedConnection,
    subscribed_topics: Vec<String>,
    subscribed: bool,
}

impl RedisPubSub {
    pub async fn new(url: &str) -> Result<Self, PubSubError> {
        let client = redis::Client::open(url).map_err(|e| {
            log::error!("Failed to create Redis client with URL '{}': {}", url, e);
            PubSubError::CreationError(format!("Failed to create Redis client: {}", e))
        })?;

        let subscriber = client.get_async_pubsub().await.map_err(|e| {
            log::error!("Failed to establish Redis subscriber connection: {}", e);
            PubSubError::ConnectionError(format!("Failed to connect subscriber: {}", e))
        })?;

        let publisher = client
            .get_multiplexed_async_connection()
            .await
            .map_err(|e| {
                log::error!("Failed to establish Redis publisher connection: {}", e);
                PubSubError::ConnectionError(format!("Failed to connect publisher: {}", e))
            })?;

        Ok(Self {
            subscriber,
            publisher,
            subscribed_topics: Vec::new(),
            subscribed: false,
        })
    }
}

#[async_trait]
impl TSubscriber for RedisPubSub {
    fn subscribe(&mut self, topics: Vec<String>) -> Result<(), PubSubError> {
        self.subscribed_topics.extend(topics);
        Ok(())
    }

    async fn listen_until_receive_message(&mut self) -> Result<(String, String), PubSubError> {
        if !self.subscribed {
            for topic in self.subscribed_topics.clone() {
                self.subscriber
                    .subscribe(topic.as_str())
                    .await
                    .map_err(|e| {
                        log::error!("Failed to subscribe to topic '{}': {}", topic, e);
                        PubSubError::ConnectionError(format!(
                            "Failed to subscribe to '{}': {}",
                            topic, e
                        ))
                    })?;
            }
            self.subscribed = true;
        }

        let msg = self
            .subscriber
            .on_message()
            .next()
            .await
            .ok_or_else(|| PubSubError::ConnectionError("PubSub stream closed".to_string()))?;

        let payload: String = msg.get_payload().map_err(|e| {
            log::error!("Failed to get message payload: {}", e);
            PubSubError::MessageError(format!("Failed to get message payload: {}", e))
        })?;

        Ok((payload, msg.get_channel_name().to_string()))
    }
}

#[async_trait]
impl TPublisher for RedisPubSub {
    async fn publish(&mut self, topic: String, message: String) -> Result<(), PubSubError> {
        use redis::AsyncCommands;
        self.publisher
            .publish::<_, _, ()>(topic.clone(), message)
            .await
            .map_err(|e| {
                log::error!("Failed to publish message to topic '{}': {}", topic, e);
                PubSubError::MessageError(format!("Failed to publish to '{}': {}", topic, e))
            })
    }
}
