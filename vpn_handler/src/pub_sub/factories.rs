use crate::errors::PubSubError;
use crate::pub_sub::redis_pubsub::RedisPubSub;
use crate::pub_sub::traits::PubSubProvider;
use std::str::FromStr;

pub type DynPubSub = Box<dyn PubSubProvider>;

#[derive(Clone)]
pub enum PubSubType {
    Redis,
}

impl FromStr for PubSubType {
    type Err = ();

    fn from_str(input: &str) -> Result<PubSubType, Self::Err> {
        match input {
            "Redis" => Ok(PubSubType::Redis),
            _ => Err(()),
        }
    }
}

pub struct PubSubFactory;

impl PubSubFactory {
    pub fn create_pubsub_instance(
        pubsub_type: PubSubType,
        url: &str,
    ) -> Result<DynPubSub, PubSubError> {
        match pubsub_type {
            PubSubType::Redis => Ok(Box::new(RedisPubSub::new(url))),
        }
    }
}
