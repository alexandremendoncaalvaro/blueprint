package module

// Reporter abstrai a comunicacao de progresso com o usuario.
// Implementado de forma diferente pelo TUI (spinners) e headless (log).
type Reporter interface {
	Info(msg string)
	Success(msg string)
	Warn(msg string)
	Error(msg string)
	Step(current, total int, msg string)
}
