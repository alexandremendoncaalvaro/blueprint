// Package profile define perfis de instalacao que filtram modulos por tags.
package profile

// Profile define quais modulos incluir/excluir por tags.
type Profile struct {
	Name        string   // Nome do perfil (full, minimal, server)
	Description string   // Descricao para exibicao
	Tags        []string // Tags incluidas
	ExcludeTags []string // Tags excluidas
}
