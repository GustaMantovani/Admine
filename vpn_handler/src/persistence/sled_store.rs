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
    fn set(&self, key: String, value: String) -> Result<(), String> {
        self.db
            .insert(key, value.into_bytes())
            .map_err(|e| e.to_string())?;
        self.db.flush().map_err(|e| e.to_string())?;
        Ok(())
    }

    fn get(&self, key: &str) -> Result<Option<String>, String> {
        match self.db.get(key).map_err(|e| e.to_string())? {
            Some(value) => Ok(Some(
                String::from_utf8(value.to_vec()).map_err(|e| e.to_string())?,
            )),
            None => Ok(None),
        }
    }
}
