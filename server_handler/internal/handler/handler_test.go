package handler

import (
	"testing"
)

func Test_Handler_Should_Start_Server_And_Stop(t *testing.T) {
	err := ManageCommand("start")

	if err != nil {
		t.Errorf("A função deveria ocorrer sem haver nenhum erro.")
	}

	err = ManageCommand("stop")

	if err != nil {
		t.Errorf("A função deveria ocorrer sem haver nenhum erro.")
	}
}

func Test_Handler_Should_Do_Nothing(t *testing.T) {
	err := ManageCommand("brainless")
	print("opa")

	if err == nil {
		t.Errorf("The function should return an error.")
	}
}
