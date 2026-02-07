package profile

import (
	"testing"

	"github.com/ale/blueprint/internal/module"
)

// stubModule e um modulo simples para testes.
type stubModule struct {
	name string
	tags []string
}

func (s *stubModule) Name() string        { return s.name }
func (s *stubModule) Description() string { return s.name }
func (s *stubModule) Tags() []string      { return s.tags }

func TestResolve(t *testing.T) {
	tests := []struct {
		name     string
		profile  Profile
		modules  []module.Module
		expected []string
	}{
		{
			name:    "tag incluida retorna modulo",
			profile: Profile{Tags: []string{"shell"}},
			modules: []module.Module{
				&stubModule{name: "starship", tags: []string{"shell"}},
			},
			expected: []string{"starship"},
		},
		{
			name:    "tag excluida remove modulo",
			profile: Profile{Tags: []string{"shell", "desktop"}, ExcludeTags: []string{"desktop"}},
			modules: []module.Module{
				&stubModule{name: "starship", tags: []string{"shell"}},
				&stubModule{name: "cedilla", tags: []string{"desktop"}},
			},
			expected: []string{"starship"},
		},
		{
			name:    "modulo com tag incluida e excluida e removido",
			profile: Profile{Tags: []string{"shell", "system"}, ExcludeTags: []string{"system"}},
			modules: []module.Module{
				&stubModule{name: "update", tags: []string{"system"}},
			},
			expected: nil,
		},
		{
			name:    "modulo sem tag incluida e ignorado",
			profile: Profile{Tags: []string{"desktop"}},
			modules: []module.Module{
				&stubModule{name: "starship", tags: []string{"shell"}},
			},
			expected: nil,
		},
		{
			name:    "perfil full inclui tudo",
			profile: Full,
			modules: []module.Module{
				&stubModule{name: "starship", tags: []string{"shell"}},
				&stubModule{name: "cedilla", tags: []string{"desktop"}},
				&stubModule{name: "update", tags: []string{"system"}},
			},
			expected: []string{"starship", "cedilla", "update"},
		},
		{
			name:    "perfil minimal exclui desktop e system",
			profile: Minimal,
			modules: []module.Module{
				&stubModule{name: "starship", tags: []string{"shell"}},
				&stubModule{name: "cedilla", tags: []string{"desktop"}},
				&stubModule{name: "update", tags: []string{"system"}},
			},
			expected: []string{"starship"},
		},
		{
			name:    "perfil server exclui desktop",
			profile: Server,
			modules: []module.Module{
				&stubModule{name: "starship", tags: []string{"shell"}},
				&stubModule{name: "cedilla", tags: []string{"desktop"}},
				&stubModule{name: "update", tags: []string{"system"}},
			},
			expected: []string{"starship", "update"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := module.NewRegistry()
			for _, m := range tt.modules {
				if err := reg.Register(m); err != nil {
					t.Fatalf("erro ao registrar modulo: %v", err)
				}
			}

			result := Resolve(tt.profile, reg)

			if len(result) != len(tt.expected) {
				names := make([]string, len(result))
				for i, r := range result {
					names[i] = r.Name()
				}
				t.Fatalf("esperava %d modulos %v, obteve %d: %v", len(tt.expected), tt.expected, len(result), names)
			}

			for i, name := range tt.expected {
				if result[i].Name() != name {
					t.Errorf("posicao %d: esperava %s, obteve %s", i, name, result[i].Name())
				}
			}
		})
	}
}

func TestMatchesTags(t *testing.T) {
	tests := []struct {
		name        string
		moduleTags  []string
		includeTags []string
		excludeTags []string
		expected    bool
	}{
		{
			name:        "tag incluida",
			moduleTags:  []string{"shell"},
			includeTags: []string{"shell"},
			excludeTags: nil,
			expected:    true,
		},
		{
			name:        "tag excluida",
			moduleTags:  []string{"desktop"},
			includeTags: []string{"shell", "desktop"},
			excludeTags: []string{"desktop"},
			expected:    false,
		},
		{
			name:        "nenhuma tag incluida",
			moduleTags:  []string{"other"},
			includeTags: []string{"shell"},
			excludeTags: nil,
			expected:    false,
		},
		{
			name:        "tag incluida e outra excluida diferentes",
			moduleTags:  []string{"shell"},
			includeTags: []string{"shell", "desktop"},
			excludeTags: []string{"desktop"},
			expected:    true,
		},
		{
			name:        "sem tags incluidas",
			moduleTags:  []string{"shell"},
			includeTags: nil,
			excludeTags: nil,
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchesTags(tt.moduleTags, tt.includeTags, tt.excludeTags)
			if got != tt.expected {
				t.Errorf("matchesTags(%v, %v, %v) = %v, esperava %v",
					tt.moduleTags, tt.includeTags, tt.excludeTags, got, tt.expected)
			}
		})
	}
}
