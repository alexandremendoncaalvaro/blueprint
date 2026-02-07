package module

// StatusKind representa o estado de um modulo no sistema.
type StatusKind int

const (
	Installed StatusKind = iota // Modulo ja esta configurado
	Missing                    // Modulo nao esta configurado
	Partial                    // Modulo parcialmente configurado
	Skipped                    // Modulo pulado (guard retornou false)
)

// String retorna a representacao textual do status.
func (s StatusKind) String() string {
	switch s {
	case Installed:
		return "instalado"
	case Missing:
		return "ausente"
	case Partial:
		return "parcial"
	case Skipped:
		return "pulado"
	default:
		return "desconhecido"
	}
}

// Status agrupa o estado e uma mensagem descritiva.
type Status struct {
	Kind    StatusKind
	Message string
}
