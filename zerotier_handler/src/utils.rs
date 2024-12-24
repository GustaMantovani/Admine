use std::fs::{File, OpenOptions};
use std::io::{self, prelude::*, Write};

pub fn write_to_file(file_path: String, content: String) -> io::Result<()> {
    // Open the file in append mode or create it if it doesn't exist
    let mut file = OpenOptions::new()
        .write(true) // Allow writing
        .create(true) // Create the file if it doesn't exist
        .open(file_path)?;

    // Write the content to the file
    file.write_all(content.as_bytes())?;
    Ok(())
}

pub fn read_file(path: String) -> io::Result<String> {
    let mut file = File::open(path)?;
    let mut contents = String::new();
    file.read_to_string(&mut contents)?;
    Ok(contents)
}
