use crate::config::Config;
use crate::persistence::key_value_storage_factory::{DynKeyValueStore, StoreFactory};
use std::sync::{Mutex, OnceLock};

pub struct StorageManager {
    store: Mutex<DynKeyValueStore>,
}

static STORAGE: OnceLock<StorageManager> = OnceLock::new();

impl StorageManager {
    fn new() -> Result<Self, String> {
        let config = Config::instance();
        let store = StoreFactory::create_store_instance(
            config.db_config().store_type().clone(),
            config.db_config().path(),
        )?;
        
        Ok(Self {
            store: Mutex::new(store),
        })
    }

    pub fn instance() -> &'static StorageManager {
        STORAGE.get_or_init(|| {
            Self::new().expect("Failed to initialize storage manager")
        })
    }
    
    pub fn get_store(&self) -> &Mutex<DynKeyValueStore> {
        &self.store
    }
}
