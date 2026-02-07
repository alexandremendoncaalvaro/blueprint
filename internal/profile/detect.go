package profile

import "github.com/ale/dotfiles/internal/module"

// Detect escolhe o perfil automaticamente baseado no ambiente.
//
// Regras:
//   - Container → minimal (sem desktop, sem system — só shell)
//   - Sem sessão gráfica ($DISPLAY e $WAYLAND_DISPLAY vazios) → server
//   - Todo o resto → full
func Detect(sys module.System) Profile {
	if sys.IsContainer() {
		return Minimal
	}

	if sys.Env("DISPLAY") == "" && sys.Env("WAYLAND_DISPLAY") == "" {
		return Server
	}

	return Full
}
