# dotfiles

CLI para configurar e manter um desktop [Bluefin](https://projectbluefin.io) (Fedora Atomic). Escrito em Go, com TUI interativo via [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## O que faz

O programa organiza configurações em **módulos** independentes, filtrados por **perfil**:

| Módulo | Tags | O que configura |
|--------|------|-----------------|
| `starship` | shell | Instala o prompt [Starship](https://starship.rs), cria symlink do config e adiciona init ao `.bashrc`/`.zshrc` |
| `cedilla-fix` | desktop | Corrige cedilha no Bluefin (Wayland/GNOME) via `~/.XCompose` |
| `bluefin-update` | system | Atualiza rpm-ostree, Flatpak, firmware (fwupd) e Distrobox |

### Perfis

| Perfil | Tags incluídas | Excluídas | Caso de uso |
|--------|---------------|-----------|-------------|
| `full` | shell, desktop, system | — | Desktop Bluefin completo |
| `minimal` | shell | desktop, system | Devcontainer / CI |
| `server` | shell, system | desktop | Servidor sem desktop |

## Instalação

```bash
git clone https://github.com/ale/dotfiles.git ~/dotfiles
cd ~/dotfiles
make build
```

O binário fica em `bin/dotfiles`. Para colocar no PATH:

```bash
mkdir -p ~/.local/bin
ln -sf ~/dotfiles/bin/dotfiles ~/.local/bin/dotfiles
```

> Requer Go 1.24+. No Bluefin, compile dentro de um distrobox com Go instalado.

## Uso

```bash
# TUI interativo (perfil full)
dotfiles apply

# Headless com perfil específico
dotfiles apply -p minimal --headless

# Simular sem executar
dotfiles apply --dry-run --headless

# Ver status dos módulos
dotfiles status

# Atualizar sistema (atalho para bluefin-update)
dotfiles update

# Versão
dotfiles version
```

## Estrutura

```
cmd/dotfiles/       → Entry point
internal/
  module/           → Interfaces de domínio (Module, System, Guard, Checker, Applier)
  modules/          → Implementações dos módulos
  profile/          → Perfis e resolução de tags
  orchestrator/     → Guard → Check → Apply pipeline
  system/           → Implementações de System (Real, Mock, DryRun)
  cli/              → Comandos Cobra
  tui/              → Interface Bubble Tea
configs/            → Arquivos de configuração (starship.toml)
```

## Adicionando um módulo

1. Crie `internal/modules/nome/nome.go` implementando `module.Module` + as interfaces necessárias (`Checker`, `Applier`, opcionalmente `Guard`)
2. Registre em `cmd/dotfiles/main.go` com `reg.Register(nome.New())`

O módulo só precisa implementar o que usa — o orchestrator detecta as interfaces por type assertion.

## Testes

```bash
make test
```
