use crate::persistence::storage_manager::StorageManager;

pub trait KeyValueStore {
    fn set(&self, key: String, value: String) -> Result<(), String>;
    fn get(&self, key: &str) -> Result<Option<String>, String>;
}

// Funções helper globais
pub fn set_global(key: String, value: String) -> Result<(), String> {
    let store = StorageManager::instance().get_store();
    let guard = store.lock().unwrap();
    guard.set(key, value)
}

pub fn get_global(key: &str) -> Result<Option<String>, String> {
    let store = StorageManager::instance().get_store();
    let guard = store.lock().unwrap();
    guard.get(key)
}
