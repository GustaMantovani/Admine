package server

import (
	"server-handler/internal/config"
	"testing"
)

// Testes end to end
func Test_Server_Should_Not_Fatal(t *testing.T) {
	err := StartServerCompose()

	if err != nil {
		t.Errorf("erro ao iniciar o servidor: %v", err)
	}

	StopServerCompose()

	t.Logf("servidor inicializado e derrubado com sucesso")
}

func Test_Server_ShouldFail(t *testing.T) {
	config.GetInstance().ComposeAbsPath = "non-existing-path"

	err := StartServerCompose()

	if err == nil {
		t.Errorf("servidor startado com sucesso quando deveria falhar")
	}

	t.Logf("erro retornado como esperado: %v", err)
}
