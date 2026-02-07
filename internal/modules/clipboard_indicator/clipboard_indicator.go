// Package clipboard_indicator instala e ativa a extensão Clipboard Indicator (Tudmotu) no GNOME.
// Gerenciador de clipboard com histórico, busca e atalhos.
package clipboard_indicator

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ale/dotfiles/internal/module"
)

const extensionUUID = "clipboard-indicator@tudmotu.com"

// Module implementa a instalação do Clipboard Indicator.
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string        { return "clipboard-indicator" }
func (m *Module) Description() string { return "Clipboard Indicator (historico de clipboard no GNOME)" }
func (m *Module) Tags() []string      { return []string{"desktop"} }

func (m *Module) ShouldRun(_ context.Context, sys module.System) (bool, string) {
	if sys.IsContainer() {
		return false, "dentro de container"
	}
	if sys.Env("WAYLAND_DISPLAY") == "" && sys.Env("DISPLAY") == "" {
		return false, "sem sessao grafica (faca login no desktop primeiro)"
	}
	if !sys.CommandExists("gnome-extensions") {
		return false, "gnome-extensions nao disponivel (requer GNOME Shell)"
	}
	return true, ""
}

func (m *Module) Check(ctx context.Context, sys module.System) (module.Status, error) {
	out, err := sys.Exec(ctx, "gnome-extensions", "show", extensionUUID)
	if err != nil {
		return module.Status{Kind: module.Missing, Message: "Clipboard Indicator nao instalado"}, nil
	}
	if strings.Contains(out, "ENABLED") {
		return module.Status{Kind: module.Installed, Message: "Clipboard Indicator instalado e ativo"}, nil
	}
	return module.Status{Kind: module.Partial, Message: "Clipboard Indicator instalado mas desativado"}, nil
}

func (m *Module) Apply(ctx context.Context, sys module.System, reporter module.Reporter) error {
	total := 3

	// 1. Detectar versão do GNOME
	reporter.Step(1, total, "Detectando versao do GNOME...")
	gnomeVer, err := detectGnomeVersion(ctx, sys)
	if err != nil {
		return err
	}
	reporter.Info(fmt.Sprintf("GNOME Shell %s", gnomeVer))

	// 2. Instalar se necessário
	reporter.Step(2, total, "Verificando Clipboard Indicator...")
	out, _ := sys.Exec(ctx, "gnome-extensions", "show", extensionUUID)
	if !strings.Contains(out, extensionUUID) {
		reporter.Info("Baixando Clipboard Indicator de extensions.gnome.org...")
		if err := m.installFromGnomeExtensions(ctx, sys, gnomeVer); err != nil {
			return fmt.Errorf("erro ao instalar Clipboard Indicator: %w", err)
		}
		reporter.Success("Clipboard Indicator instalado")
	} else {
		reporter.Info("Clipboard Indicator ja instalado")
	}

	// 3. Ativar
	reporter.Step(3, total, "Ativando Clipboard Indicator...")
	if _, err := sys.Exec(ctx, "gnome-extensions", "enable", extensionUUID); err != nil {
		reporter.Warn("Clipboard Indicator sera ativado apos re-login")
	} else {
		reporter.Success("Clipboard Indicator ativo")
	}

	reporter.Info("Faca logout e login se nao aparecer imediatamente")
	return nil
}

// extensionInfo representa a resposta da API do extensions.gnome.org.
type extensionInfo struct {
	DownloadURL string `json:"download_url"`
}

func (m *Module) installFromGnomeExtensions(ctx context.Context, sys module.System, gnomeVer string) error {
	apiURL := fmt.Sprintf(
		"https://extensions.gnome.org/extension-info/?uuid=%s&shell_version=%s",
		extensionUUID, gnomeVer,
	)

	jsonOut, err := sys.Exec(ctx, "curl", "-sfL", apiURL)
	if err != nil {
		return fmt.Errorf("erro ao consultar extensions.gnome.org (verifique sua conexao com a internet): %w", err)
	}

	var info extensionInfo
	if err := json.Unmarshal([]byte(jsonOut), &info); err != nil {
		return fmt.Errorf("resposta inesperada da API: %w", err)
	}
	if info.DownloadURL == "" {
		return fmt.Errorf("Clipboard Indicator nao disponivel para GNOME Shell %s — verifique se ha uma versao compativel em https://extensions.gnome.org", gnomeVer)
	}

	downloadURL := "https://extensions.gnome.org" + info.DownloadURL
	zipPath := "/tmp/clipboard-indicator.zip"

	if _, err := sys.Exec(ctx, "curl", "-sfL", "-o", zipPath, downloadURL); err != nil {
		return fmt.Errorf("erro ao baixar Clipboard Indicator: %w", err)
	}

	if _, err := sys.Exec(ctx, "gnome-extensions", "install", "--force", zipPath); err != nil {
		return fmt.Errorf("erro ao instalar extensao: %w", err)
	}

	return nil
}

func detectGnomeVersion(ctx context.Context, sys module.System) (string, error) {
	out, err := sys.Exec(ctx, "gnome-shell", "--version")
	if err != nil {
		return "", fmt.Errorf("gnome-shell nao encontrado: %w", err)
	}
	// "GNOME Shell 46.2" → "46"
	parts := strings.Fields(out)
	if len(parts) < 3 {
		return "", fmt.Errorf("saida inesperada: %s", out)
	}
	ver := parts[2]
	if dot := strings.Index(ver, "."); dot > 0 {
		ver = ver[:dot]
	}
	return ver, nil
}
