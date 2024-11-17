package internal

import "os/exec"

// Executa o comando 'docker compose up -d' em determinado diretorio e retorna o output do cmd
func StartServerDockerCompose(composeDirectory string) ([]byte, error) {
	cmd := exec.Command("docker", "compose", "up", "-d")
	cmd.Dir = composeDirectory

	return cmd.CombinedOutput()
}
