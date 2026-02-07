package profile

import "github.com/ale/blueprint/internal/module"

// Resolve filtra os modulos do registry baseado no perfil.
// Retorna apenas os modulos cujas tags sao incluidas e nao excluidas pelo perfil.
func Resolve(p Profile, reg *module.Registry) []module.Module {
	var result []module.Module

	for _, m := range reg.All() {
		if matchesTags(m.Tags(), p.Tags, p.ExcludeTags) {
			result = append(result, m)
		}
	}

	return result
}

// matchesTags verifica se as tags de um modulo sao compativeis com o perfil.
// O modulo precisa ter pelo menos uma tag incluida e nenhuma tag excluida.
func matchesTags(moduleTags, includeTags, excludeTags []string) bool {
	// Verifica se alguma tag do modulo esta excluida
	for _, mt := range moduleTags {
		for _, et := range excludeTags {
			if mt == et {
				return false
			}
		}
	}

	// Verifica se alguma tag do modulo esta incluida
	for _, mt := range moduleTags {
		for _, it := range includeTags {
			if mt == it {
				return true
			}
		}
	}

	return false
}
