// Package gnome_focus configura o modo foco no GNOME:
// F11 (fullscreen) move a janela para um workspace exclusivo (estilo macOS).
// Instala uma extensão GNOME Shell leve que gerencia isso automaticamente.
package gnome_focus

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ale/blueprint/internal/gnome"
	"github.com/ale/blueprint/internal/module"
)

const extensionUUID = "focus-mode@blueprint"

// Module implementa o modo foco via extensão GNOME Shell.
type Module struct {
	// ExtensionSource é o caminho para configs/gnome-extensions/focus-mode@blueprint/
	ExtensionSource string
}

func New(extensionSource string) *Module {
	return &Module{ExtensionSource: extensionSource}
}

func (m *Module) Name() string        { return "gnome-focus-mode" }
func (m *Module) Description() string { return "Modo foco: F11 = fullscreen + workspace exclusivo" }
func (m *Module) Tags() []string      { return []string{"desktop"} }

func (m *Module) ShouldRun(_ context.Context, sys module.System) (bool, string) {
	return gnome.ShouldRunGuard(sys)
}

func (m *Module) Check(ctx context.Context, sys module.System) (module.Status, error) {
	// Usa o checker compartilhado para extensao
	status, err := gnome.CheckExtension(ctx, sys, extensionUUID, "Focus-mode")
	if err != nil || status.Kind != module.Installed {
		return status, err
	}

	// Check adicional: dynamic workspaces
	dynWs, _ := sys.Exec(ctx, "dconf", "read", "/org/gnome/mutter/dynamic-workspaces")
	if strings.TrimSpace(dynWs) == "false" {
		return module.Status{Kind: module.Partial, Message: "Workspaces dinamicos desativados"}, nil
	}

	return status, nil
}

func (m *Module) Apply(ctx context.Context, sys module.System, reporter module.Reporter) error {
	// Limpar UUID legado (rename dotfiles → blueprint)
	sys.Exec(ctx, "gnome-extensions", "disable", "focus-mode@dotfiles")

	// 1. Ativar workspaces dinamicos
	reporter.Step(1, 3, "Ativando workspaces dinamicos...")
	if _, err := sys.Exec(ctx, "dconf", "write", "/org/gnome/mutter/dynamic-workspaces", "true"); err != nil {
		return fmt.Errorf("erro ao ativar workspaces dinamicos: %w", err)
	}
	reporter.Success("Workspaces dinamicos ativados")

	// 2. Instalar ou atualizar extensao
	_, showErr := sys.Exec(ctx, "gnome-extensions", "show", extensionUUID)
	if showErr != nil {
		// Primeira instalacao: zip + install para registrar no GNOME Shell
		reporter.Step(2, 3, "Instalando extensao focus-mode...")
		zipPath := filepath.Join(os.TempDir(), extensionUUID+".zip")
		if _, err := sys.Exec(ctx, "zip", "-j", zipPath, filepath.Join(m.ExtensionSource, "metadata.json"), filepath.Join(m.ExtensionSource, "extension.js")); err != nil {
			return fmt.Errorf("erro ao criar zip da extensao: %w", err)
		}
		if _, err := sys.Exec(ctx, "gnome-extensions", "install", "--force", zipPath); err != nil {
			return fmt.Errorf("erro ao instalar extensao: %w", err)
		}
		reporter.Success("Extensao instalada")
	} else {
		// Extensao ja existe: copiar arquivos sem install --force (que reseta o estado)
		reporter.Step(2, 3, "Atualizando extensao focus-mode...")
		extDir := filepath.Join(sys.HomeDir(), ".local", "share", "gnome-shell", "extensions", extensionUUID)
		if err := sys.MkdirAll(extDir, 0o755); err != nil {
			return fmt.Errorf("erro ao criar diretorio da extensao: %w", err)
		}
		for _, fname := range []string{"metadata.json", "extension.js"} {
			data, err := sys.ReadFile(filepath.Join(m.ExtensionSource, fname))
			if err != nil {
				return fmt.Errorf("erro ao ler %s: %w", fname, err)
			}
			if err := sys.WriteFile(filepath.Join(extDir, fname), data, 0o644); err != nil {
				return fmt.Errorf("erro ao copiar %s: %w", fname, err)
			}
		}
		reporter.Success("Extensao atualizada")
	}

	// 3. Ativar
	reporter.Step(3, 3, "Ativando extensao...")
	_, err := sys.Exec(ctx, "gnome-extensions", "enable", extensionUUID)
	if err != nil {
		reporter.Warn("Extensao sera ativada apos re-login")
	} else {
		reporter.Success("Focus mode ativo")
	}

	reporter.Info("F11 agora envia a janela para um workspace exclusivo")
	return nil
}
