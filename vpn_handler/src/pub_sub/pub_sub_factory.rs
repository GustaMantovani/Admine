use crate::errors::PubSubError;
use crate::pub_sub::pub_sub::PubSubProvider;
use crate::pub_sub::redis_pubsub::RedisPubSub;
use strum::EnumString;

pub type DynPubSub = Box<dyn PubSubProvider>;

#[derive(Clone, Debug, PartialEq, EnumString)]
pub enum PubSubType {
    Redis,
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
