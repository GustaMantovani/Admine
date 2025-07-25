use crate::persistence::key_value_storage::DynKeyValueStore;
use crate::persistence::sled_store::SledStore;
use serde::Deserialize;
use strum::EnumString;

#[derive(Clone, Debug, PartialEq, EnumString, Deserialize)]
pub enum StoreType {
    Sled,
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
