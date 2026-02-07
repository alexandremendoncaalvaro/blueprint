package cedilla

import (
	"context"
	"strings"
	"testing"

	"github.com/ale/dotfiles/internal/module"
	"github.com/ale/dotfiles/internal/module/moduletest"
	"github.com/ale/dotfiles/internal/system"
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
	mock.Files["/home/test/.XCompose"] = []byte(`include "%L"

# BEGIN BLUEFIN CEDILLA
<dead_acute> <c> : "ç"
<dead_acute> <C> : "Ç"
# END BLUEFIN CEDILLA
`)

	mod := New()
	status, err := mod.Check(context.Background(), mock)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if status.Kind != module.Installed {
		t.Errorf("esperava Installed, obteve %s", status.Kind)
	}
}

func TestApply_CreatesNewFile(t *testing.T) {
	mock := system.NewMock()
	mock.EnvVars["XDG_SESSION_TYPE"] = "wayland"

	mod := New()
	reporter := moduletest.NoopReporter()

	err := mod.Apply(context.Background(), mock, reporter)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	data, ok := mock.Files["/home/test/.XCompose"]
	if !ok {
		t.Fatal("~/.XCompose nao foi criado")
	}

	content := string(data)
	if !strings.Contains(content, `include "%L"`) {
		t.Error("falta include")
	}
	if !strings.Contains(content, beginMarker) {
		t.Error("falta bloco de cedilha")
	}
	if !strings.Contains(content, `<dead_acute> <c> : "ç"`) {
		t.Error("falta regra para c minusculo")
	}
}

func TestApply_IdempotentUpdate(t *testing.T) {
	mock := system.NewMock()
	mock.Files["/home/test/.XCompose"] = []byte(`include "%L"

# BEGIN BLUEFIN CEDILLA
<dead_acute> <c> : "old"
# END BLUEFIN CEDILLA
`)

	mod := New()
	reporter := moduletest.NoopReporter()

	err := mod.Apply(context.Background(), mock, reporter)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	content := string(mock.Files["/home/test/.XCompose"])

	// Deve ter apenas um bloco
	count := strings.Count(content, beginMarker)
	if count != 1 {
		t.Errorf("esperava 1 bloco BEGIN, encontrou %d", count)
	}

	// Deve ter as regras novas, nao as antigas
	if strings.Contains(content, "old") {
		t.Error("regra antiga nao foi removida")
	}
	if !strings.Contains(content, `<dead_acute> <c> : "ç"`) {
		t.Error("regra nova nao foi adicionada")
	}
}

func TestRemoveBlock(t *testing.T) {
	input := `line1
# BEGIN BLUEFIN CEDILLA
some content
# END BLUEFIN CEDILLA
line2`

	result := removeBlock(input, beginMarker, endMarker)

	if strings.Contains(result, beginMarker) {
		t.Error("bloco nao foi removido")
	}
	if !strings.Contains(result, "line1") || !strings.Contains(result, "line2") {
		t.Error("conteudo ao redor foi removido incorretamente")
	}
}
