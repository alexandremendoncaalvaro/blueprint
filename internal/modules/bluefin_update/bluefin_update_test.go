package bluefin_update

import (
	"context"
	"fmt"
	"testing"

	"github.com/ale/dotfiles/internal/module"
	"github.com/ale/dotfiles/internal/module/moduletest"
	"github.com/ale/dotfiles/internal/system"
)

func TestShouldRun_SkipInContainer(t *testing.T) {
	mock := system.NewMock()
	mock.Container = true

	mod := New()
	ok, _ := mod.ShouldRun(context.Background(), mock)
	if ok {
		t.Error("deveria pular em container")
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

func TestCheck_AlwaysMissing(t *testing.T) {
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

func TestApply_AllCommandsAvailable(t *testing.T) {
	mock := system.NewMock()
	mock.Commands["rpm-ostree"] = true
	mock.Commands["flatpak"] = true
	mock.Commands["fwupdmgr"] = true
	mock.Commands["distrobox"] = true

	mod := New()
	reporter := moduletest.NoopReporter()

	err := mod.Apply(context.Background(), mock, reporter)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	// Verifica que todos os comandos foram executados
	if len(mock.ExecLog) != 5 {
		t.Errorf("esperava 5 comandos executados, obteve %d: %v", len(mock.ExecLog), mock.ExecLog)
	}
}

func TestApply_SkipOptionalMissing(t *testing.T) {
	mock := system.NewMock()
	mock.Commands["rpm-ostree"] = true
	mock.Commands["flatpak"] = true
	// fwupdmgr e distrobox nao disponiveis

	mod := New()
	reporter := moduletest.NoopReporter()

	err := mod.Apply(context.Background(), mock, reporter)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	// Apenas rpm-ostree e flatpak devem ter sido executados
	if len(mock.ExecLog) != 2 {
		t.Errorf("esperava 2 comandos executados, obteve %d: %v", len(mock.ExecLog), mock.ExecLog)
	}
}

func TestApply_FailOnMissingRequired(t *testing.T) {
	mock := system.NewMock()
	// rpm-ostree nao disponivel (obrigatorio)
	mock.Commands["flatpak"] = true

	mod := New()
	reporter := moduletest.NoopReporter()

	err := mod.Apply(context.Background(), mock, reporter)
	if err == nil {
		t.Error("esperava erro quando comando obrigatorio esta ausente")
	}
}

func TestApply_OptionalFailureContinues(t *testing.T) {
	mock := system.NewMock()
	mock.Commands["rpm-ostree"] = true
	mock.Commands["flatpak"] = true
	mock.Commands["fwupdmgr"] = true
	mock.ExecResults["fwupdmgr refresh"] = system.ExecResult{Err: fmt.Errorf("falha")}

	mod := New()
	reporter := moduletest.NoopReporter()

	err := mod.Apply(context.Background(), mock, reporter)
	if err != nil {
		t.Fatalf("nao deveria falhar com erro opcional: %v", err)
	}
}
