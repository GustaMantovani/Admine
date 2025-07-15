use crate::app_context::AppContext;

pub trait KeyValueStore {
    fn set(&self, key: String, value: String) -> Result<(), String>;
    fn get(&self, key: &str) -> Result<Option<String>, String>;
}

// Funções helper globais
pub fn set_global(key: String, value: String) -> Result<(), String> {
    let storage = AppContext::instance().storage();
    let guard = storage.lock().unwrap();
    guard.set(key, value)
}

pub fn get_global(key: &str) -> Result<Option<String>, String> {
    let storage = AppContext::instance().storage();
    let guard = storage.lock().unwrap();
    guard.get(key)
}
