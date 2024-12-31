package internal

import (
	"os/exec"
	"strings"
)

// Executa o comando 'docker compose up -d' em determinado diretorio e retorna o output do cmd
func StartServerDockerCompose(composeDirectory string) ([]byte, error) {
	cmd := exec.Command("docker", "compose", "up", "-d")
	cmd.Dir = composeDirectory

	return cmd.CombinedOutput()
}

func GetZeroTierNodeID() string {
	cmd := exec.Command("docker", "exec", "-i", "minecraft-server-mine_server-1", "/bin/bash", "-c", "zerotier-cli info")
	output, err := cmd.CombinedOutput()

	if err != nil {
		panic(err)
	}

	outputStr := string(output)

	parts := strings.Split(outputStr, " ")

	return parts[2]
}
