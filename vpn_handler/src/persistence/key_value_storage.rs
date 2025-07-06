pub trait KeyValueStore {
    fn set(&mut self, key: String, value: String) -> Result<(), String>;
    fn get(&self, key: &str) -> Result<Option<String>, String>;
}
