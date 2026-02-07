package clipboard_indicator

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
	ok, _ := mod.ShouldRun(context.Background(), mock)
	if ok {
		t.Error("deveria pular em container")
	}
}

func TestShouldRun_SkipWithoutDisplay(t *testing.T) {
	mock := system.NewMock()

	mod := New()
	ok, _ := mod.ShouldRun(context.Background(), mock)
	if ok {
		t.Error("deveria pular sem sessao grafica")
	}
}

func TestShouldRun_SkipWithoutGnomeExtensions(t *testing.T) {
	mock := system.NewMock()
	mock.EnvVars["WAYLAND_DISPLAY"] = "wayland-0"

	mod := New()
	ok, _ := mod.ShouldRun(context.Background(), mock)
	if ok {
		t.Error("deveria pular sem gnome-extensions")
	}
}

func TestShouldRun_RunOnGnomeDesktop(t *testing.T) {
	mock := system.NewMock()
	mock.EnvVars["WAYLAND_DISPLAY"] = "wayland-0"
	mock.Commands["gnome-extensions"] = true

	mod := New()
	ok, _ := mod.ShouldRun(context.Background(), mock)
	if !ok {
		t.Error("deveria rodar em desktop GNOME")
	}
}

func TestCheck_Missing(t *testing.T) {
	mock := system.NewMock()
	mock.ExecResults["gnome-extensions show clipboard-indicator@tudmotu.com"] = system.ExecResult{
		Err: fmt.Errorf("not found"),
	}

	mod := New()
	status, err := mod.Check(context.Background(), mock)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if status.Kind != module.Missing {
		t.Errorf("esperava Missing, obteve %v", status.Kind)
	}
}

func TestCheck_Installed(t *testing.T) {
	mock := system.NewMock()
	mock.ExecResults["gnome-extensions show clipboard-indicator@tudmotu.com"] = system.ExecResult{
		Output: "clipboard-indicator@tudmotu.com\n  Enabled: Yes\n  State: ACTIVE\n",
	}

	mod := New()
	status, err := mod.Check(context.Background(), mock)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if status.Kind != module.Installed {
		t.Errorf("esperava Installed, obteve %v", status.Kind)
	}
}

func TestCheck_OutOfDate(t *testing.T) {
	mock := system.NewMock()
	mock.ExecResults["gnome-extensions show clipboard-indicator@tudmotu.com"] = system.ExecResult{
		Output: "clipboard-indicator@tudmotu.com\n  Enabled: Yes\n  State: OUT OF DATE\n",
	}

	mod := New()
	status, err := mod.Check(context.Background(), mock)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if status.Kind != module.Partial {
		t.Errorf("esperava Partial para OUT OF DATE, obteve %v", status.Kind)
	}
}

func TestCheck_Partial(t *testing.T) {
	mock := system.NewMock()
	mock.ExecResults["gnome-extensions show clipboard-indicator@tudmotu.com"] = system.ExecResult{
		Output: "clipboard-indicator@tudmotu.com\n  Enabled: No\n  State: INACTIVE\n",
	}

	mod := New()
	status, err := mod.Check(context.Background(), mock)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if status.Kind != module.Partial {
		t.Errorf("esperava Partial, obteve %v", status.Kind)
	}
}

func TestApply_AlreadyInstalled(t *testing.T) {
	mock := system.NewMock()
	mock.Commands["gnome-extensions"] = true
	mock.Commands["gnome-shell"] = true
	mock.ExecResults["gnome-shell --version"] = system.ExecResult{Output: "GNOME Shell 46.2"}
	mock.ExecResults["gnome-extensions show clipboard-indicator@tudmotu.com"] = system.ExecResult{
		Output: "clipboard-indicator@tudmotu.com\n  State: ENABLED\n",
	}

	mod := New()
	reporter := moduletest.NoopReporter()

	err := mod.Apply(context.Background(), mock, reporter)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
}
