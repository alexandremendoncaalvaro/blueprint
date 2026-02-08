package usb_audio

import (
	"context"
	"fmt"
	"testing"

	"github.com/ale/blueprint/internal/module"
	"github.com/ale/blueprint/internal/module/moduletest"
	"github.com/ale/blueprint/internal/system"
)

func TestShouldRun_SkipInContainer(t *testing.T) {
	mock := system.NewMock()
	mock.Container = true

	mod := New()
	ok, reason := mod.ShouldRun(context.Background(), mock)

	if ok {
		t.Error("deveria pular em container")
	}
	if reason == "" {
		t.Error("deveria ter motivo")
	}
}

func TestShouldRun_RunOutsideContainer(t *testing.T) {
	mock := system.NewMock()
	mod := New()

	ok, _ := mod.ShouldRun(context.Background(), mock)
	if !ok {
		t.Error("deveria rodar fora de container")
	}
}

func TestCheck_Missing(t *testing.T) {
	mock := system.NewMock()
	mod := New()

	status, err := mod.Check(context.Background(), mock)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if status.Kind != module.Missing {
		t.Errorf("esperava Missing, obteve %s", status.Kind)
	}
}

func TestCheck_Installed(t *testing.T) {
	mock := system.NewMock()
	mock.Files[rulesPath] = []byte(rulesContent)

	mod := New()
	status, err := mod.Check(context.Background(), mock)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if status.Kind != module.Installed {
		t.Errorf("esperava Installed, obteve %s", status.Kind)
	}
}

func TestCheck_Partial(t *testing.T) {
	mock := system.NewMock()
	mock.Files[rulesPath] = []byte("# regra antiga\n")

	mod := New()
	status, err := mod.Check(context.Background(), mock)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if status.Kind != module.Partial {
		t.Errorf("esperava Partial, obteve %s", status.Kind)
	}
}

func TestApply_InstallsRules(t *testing.T) {
	mock := system.NewMock()

	mod := New()
	reporter := moduletest.NoopReporter()

	err := mod.Apply(context.Background(), mock, reporter)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	// Verifica arquivo temporario
	tmpFile := "/home/test/.cache/blueprint-usb-audio.rules"
	data, ok := mock.Files[tmpFile]
	if !ok {
		t.Fatal("arquivo temporario nao criado")
	}
	if string(data) != rulesContent {
		t.Errorf("conteudo inesperado: %q", string(data))
	}

	// Verifica comandos executados
	expectedCmds := []string{
		"sudo cp " + tmpFile + " " + rulesPath,
		"sudo udevadm control --reload-rules",
		"sudo udevadm trigger --subsystem-match=usb",
	}
	for _, cmd := range expectedCmds {
		found := false
		for _, logged := range mock.ExecLog {
			if logged == cmd {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("comando esperado nao executado: %s", cmd)
		}
	}
}

func TestApply_WriteFileFails(t *testing.T) {
	mock := system.NewMock()
	mock.WriteFileErr = fmt.Errorf("disco cheio")

	mod := New()
	reporter := moduletest.NoopReporter()

	err := mod.Apply(context.Background(), mock, reporter)
	if err == nil {
		t.Error("esperava erro quando WriteFile falha")
	}
}

func TestApply_CopyFails(t *testing.T) {
	mock := system.NewMock()
	tmpFile := "/home/test/.cache/blueprint-usb-audio.rules"
	mock.ExecResults["sudo cp "+tmpFile+" "+rulesPath] = system.ExecResult{
		Err: fmt.Errorf("permissao negada"),
	}

	mod := New()
	reporter := moduletest.NoopReporter()

	err := mod.Apply(context.Background(), mock, reporter)
	if err == nil {
		t.Error("esperava erro quando sudo cp falha")
	}
}

func TestApply_ReloadFails(t *testing.T) {
	mock := system.NewMock()
	mock.ExecResults["sudo udevadm control --reload-rules"] = system.ExecResult{
		Err: fmt.Errorf("udevadm falhou"),
	}

	mod := New()
	reporter := moduletest.NoopReporter()

	err := mod.Apply(context.Background(), mock, reporter)
	if err == nil {
		t.Error("esperava erro quando udevadm reload falha")
	}
}

func TestApply_TriggerFails(t *testing.T) {
	mock := system.NewMock()
	mock.ExecResults["sudo udevadm trigger --subsystem-match=usb"] = system.ExecResult{
		Err: fmt.Errorf("trigger falhou"),
	}

	mod := New()
	reporter := moduletest.NoopReporter()

	err := mod.Apply(context.Background(), mock, reporter)
	if err == nil {
		t.Error("esperava erro quando udevadm trigger falha")
	}
}
