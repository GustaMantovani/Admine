package internal

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Dados do arquivo de configuração
type ConfigFile struct {
	ServerName       string `yaml:"serverName"`
	ComposeDirectory string `yaml:"composeDirectory"`
}

// Verifica se o arquivo de configuração ~/.config/admine/server.yaml existe
func VerifyIfConfigFileExists() bool {
	configFilePath := getConfigFilePath()

	_, err := os.ReadFile(configFilePath)
	if err != nil {
		return false
	}

	return true
}

// Retorna o caminho do arquivo de configuração do servidor
func getConfigFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("erro: ", err)
	}

	configFilePath := home + "/.config/admine/server.yaml"

	return configFilePath
}

// Passa os dados do arquivo de configuração do servidor para uma struct
func GetConfigFileData() (ConfigFile, error) {
	var configFileData ConfigFile

	configFilePath := getConfigFilePath

	file, err := os.Open(configFilePath())
	if err != nil {
		return configFileData, err
	}

	decoder := yaml.NewDecoder(file)

	err = decoder.Decode(&configFileData)

	if err != nil {
		return configFileData, err
	}

	file.Close()

	return configFileData, nil
}

// Retorna o path da área de trabalho
func getLocalDirectory() string {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Erro ao obter o diretório de trabalho: %v\n", err)
		return ""
	}

	return wd
}
