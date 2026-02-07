package module

import "fmt"

// Registry armazena e organiza os modulos disponiveis.
type Registry struct {
	modules []Module
	byName  map[string]Module
}

// NewRegistry cria um registry vazio.
func NewRegistry() *Registry {
	return &Registry{
		byName: make(map[string]Module),
	}
}

// Register adiciona um modulo ao registry.
// Retorna erro se ja existir um modulo com o mesmo nome.
func (r *Registry) Register(m Module) error {
	if _, exists := r.byName[m.Name()]; exists {
		return fmt.Errorf("modulo ja registrado: %s", m.Name())
	}
	r.modules = append(r.modules, m)
	r.byName[m.Name()] = m
	return nil
}

// All retorna todos os modulos na ordem de registro.
func (r *Registry) All() []Module {
	return r.modules
}

// ByName busca um modulo pelo nome.
func (r *Registry) ByName(name string) (Module, bool) {
	m, ok := r.byName[name]
	return m, ok
}
