use crate::errors::PubSubError;

pub type DynPubSub = Box<dyn PubSubProvider>;
pub trait PubSubProvider: TPublisher + TSubscriber + Send + Sync {}
impl<T: TPublisher + TSubscriber + Send + Sync> PubSubProvider for T {}

pub trait TSubscriber {
    fn subscribe(&mut self, topics: Vec<String>) -> Result<(), PubSubError>;
    fn listen_until_receive_message(&mut self) -> Result<(String, String), PubSubError>; // Return a tuple with message payload and channel name
}

pub trait TPublisher {
    fn publish(&mut self, topic: String, message: String) -> Result<(), PubSubError>;
}
