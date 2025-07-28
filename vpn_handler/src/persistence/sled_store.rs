use crate::persistence::key_value_storage::KeyValueStore;
use sled::Db;

pub struct SledStore {
    db: Db,
}

impl SledStore {
    pub fn new(path: &str) -> Result<Self, String> {
        let db = sled::open(path).map_err(|e| e.to_string())?;
        Ok(SledStore { db })
    }
}

impl KeyValueStore for SledStore {
    fn set(&self, key: String, value: String) -> Result<(), Box<dyn std::error::Error>> {
        self.db
            .insert(key, value.as_bytes())
            .map(|_| ())
            .map_err(|e| Box::new(e) as Box<dyn std::error::Error>)
    }

    fn get(&self, key: &str) -> Option<String> {
        self.db
            .get(key)
            .ok()
            .flatten()
            .and_then(|ivec| String::from_utf8(ivec.to_vec()).ok())
    }
}
