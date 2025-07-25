pub type DynKeyValueStore = Box<dyn KeyValueStore + Send + Sync>;
pub trait KeyValueStore {
    fn set(&self, key: String, value: String) -> Result<(), Box<dyn std::error::Error>>;
    fn get(&self, key: &str) -> Option<String>;
}
