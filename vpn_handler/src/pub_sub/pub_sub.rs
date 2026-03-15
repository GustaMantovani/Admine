use crate::errors::PubSubError;
use async_trait::async_trait;

pub type DynPubSub = Box<dyn PubSubProvider>;
pub trait PubSubProvider: TPublisher + TSubscriber + Send + Sync {}
impl<T: TPublisher + TSubscriber + Send + Sync> PubSubProvider for T {}

#[async_trait]
pub trait TSubscriber {
    fn subscribe(&mut self, topics: Vec<String>) -> Result<(), PubSubError>;
    async fn listen_until_receive_message(&mut self) -> Result<(String, String), PubSubError>;
}

#[async_trait]
pub trait TPublisher {
    async fn publish(&mut self, topic: String, message: String) -> Result<(), PubSubError>;
}
