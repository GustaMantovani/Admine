use crate::errors::PubSubError;

pub trait TSubscriber {
    fn subscribe(&mut self, topics: Vec<String>) -> Result<(), PubSubError>;
    fn unsubscribe(&mut self, topic: String) -> Result<(), PubSubError>;
    fn listen_until_to_ricieve_message(&mut self) -> Result<(String, String), PubSubError>; // Return a tuple with message payload and channel name
}

pub trait TPublisher {
    fn publish(&mut self, topic: String, message: String) -> Result<(), PubSubError>;
}

pub trait PubSubProvider: TPublisher + TSubscriber + Send + Sync {}

impl<T: TPublisher + TSubscriber + Send + Sync> PubSubProvider for T {}
