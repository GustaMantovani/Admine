use std::io::{prelude::*, self, Write};
use std::fs::{File, OpenOptions};

pub fn write_to_file(file_path: String, content: String) -> io::Result<()> {
    // Abre o arquivo no modo append ou cria o arquivo se ele não existir
    let mut file = OpenOptions::new()
        .write(true) // Permite escrita
        .create(true) // Cria o arquivo se ele não existir
        .open(file_path)?;

    // Escreve o conteúdo no arquivo
    file.write_all(content.as_bytes())?;
    Ok(())
}

pub fn read_file(path: String) -> io::Result<String> {
    let mut file = File::open(path)?;
    let mut contents = String::new();
    file.read_to_string(&mut contents)?;
    Ok(contents)
}
