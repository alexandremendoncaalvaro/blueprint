// Package usb_audio configura regras udev para desabilitar autosuspend
// em dispositivos de audio USB, evitando desconexoes e stuttering.
package usb_audio

import (
	"context"
	"fmt"
	"strings"

	"github.com/ale/blueprint/internal/module"
)

const (
	rulesPath = "/etc/udev/rules.d/99-usb-audio-no-autosuspend.rules"

	rulesContent = `# Desabilita autosuspend em dispositivos de audio USB para evitar desconexoes e stuttering.
ACTION=="add", SUBSYSTEM=="usb", ATTR{bInterfaceClass}=="01", ATTR{bInterfaceSubClass}=="01", RUN+="/bin/sh -c 'echo on > /sys$DEVPATH/../power/control; echo -1 > /sys$DEVPATH/../power/autosuspend_delay_ms'"

# Desabilita autosuspend no hub Genesys Logic que hospeda os dispositivos de audio.
ACTION=="add", SUBSYSTEM=="usb", ATTR{idVendor}=="05e3", ATTR{idProduct}=="0610", ATTR{power/control}="on", ATTR{power/autosuspend_delay_ms}="-1"
`
)

// Module implementa regras udev para audio USB sem autosuspend.
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string        { return "usb-audio" }
func (m *Module) Description() string { return "Regras udev para desabilitar autosuspend em audio USB" }
func (m *Module) Tags() []string      { return []string{"system"} }

// ShouldRun retorna false dentro de containers.
func (m *Module) ShouldRun(_ context.Context, sys module.System) (bool, string) {
	if sys.IsContainer() {
		return false, "dentro de container (configuracao de sistema)"
	}
	return true, ""
}

// Check verifica se as regras udev estao instaladas com o conteudo correto.
func (m *Module) Check(_ context.Context, sys module.System) (module.Status, error) {
	data, err := sys.ReadFile(rulesPath)
	if err != nil {
		return module.Status{Kind: module.Missing, Message: "Regras udev nao instaladas"}, nil
	}

	if strings.TrimSpace(string(data)) == strings.TrimSpace(rulesContent) {
		return module.Status{Kind: module.Installed, Message: "Regras udev configuradas"}, nil
	}

	return module.Status{Kind: module.Partial, Message: "Regras udev desatualizadas"}, nil
}

// Apply instala as regras udev e recarrega o udevadm.
func (m *Module) Apply(ctx context.Context, sys module.System, reporter module.Reporter) error {
	cacheDir := sys.HomeDir() + "/.cache"
	tmpFile := cacheDir + "/blueprint-usb-audio.rules"

	// Step 1 — Escreve regras no arquivo temporario
	reporter.Step(1, 3, "Escrevendo regras udev...")

	if err := sys.WriteFile(tmpFile, []byte(rulesContent), 0o644); err != nil {
		return fmt.Errorf("erro ao escrever arquivo temporario: %w", err)
	}

	// Step 2 — Copia para /etc com sudo
	reporter.Step(2, 3, "Instalando regras udev...")

	if _, err := sys.Exec(ctx, "sudo", "cp", tmpFile, rulesPath); err != nil {
		return fmt.Errorf("erro ao copiar regras udev: %w", err)
	}

	reporter.Success("Regras udev instaladas")

	// Step 3 — Recarrega regras
	reporter.Step(3, 3, "Recarregando udev...")

	if _, err := sys.Exec(ctx, "sudo", "udevadm", "control", "--reload-rules"); err != nil {
		return fmt.Errorf("erro ao recarregar regras udev: %w", err)
	}

	if _, err := sys.Exec(ctx, "sudo", "udevadm", "trigger", "--subsystem-match=usb"); err != nil {
		return fmt.Errorf("erro ao aplicar regras udev: %w", err)
	}

	reporter.Success("Regras udev recarregadas")

	return nil
}
