use crate::persistence::key_value_store::KeyValueStore;
use crate::persistence::sled_store::SledStore;
use serde::{Deserialize, Serialize};
use std::fmt;
use std::str::FromStr;

pub type DynKeyValueStore = Box<dyn KeyValueStore + Send + Sync>;

#[derive(Clone, Debug, Serialize, Deserialize)]
pub enum StoreType {
    Sled,
}

impl fmt::Display for StoreType {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            StoreType::Sled => write!(f, "Sled"),
        }
    }
}

impl FromStr for StoreType {
    type Err = ();

    fn from_str(input: &str) -> Result<StoreType, Self::Err> {
        match input {
            "Sled" => Ok(StoreType::Sled),
            _ => Err(()),
        }
    }
}

pub struct StoreFactory;

impl StoreFactory {
    pub fn create_store_instance(
        store_type: StoreType,
        path: &str,
    ) -> Result<DynKeyValueStore, String> {
        match store_type {
            StoreType::Sled => Ok(Box::new(SledStore::new(path)?)),
        }
    }
}
