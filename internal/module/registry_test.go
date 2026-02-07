package module

import "testing"

// stubModule para testes do Registry.
type stubModule struct {
	name string
}

func (s *stubModule) Name() string        { return s.name }
func (s *stubModule) Description() string { return s.name }
func (s *stubModule) Tags() []string      { return []string{"test"} }

func TestRegistry_RegisterAndAll(t *testing.T) {
	reg := NewRegistry()

	if got := reg.All(); len(got) != 0 {
		t.Fatalf("registry novo deveria estar vazio, tem %d", len(got))
	}

	err := reg.Register(&stubModule{name: "a"})
	if err != nil {
		t.Fatalf("erro ao registrar: %v", err)
	}

	err = reg.Register(&stubModule{name: "b"})
	if err != nil {
		t.Fatalf("erro ao registrar: %v", err)
	}

	all := reg.All()
	if len(all) != 2 {
		t.Fatalf("esperava 2 modulos, obteve %d", len(all))
	}

	if all[0].Name() != "a" || all[1].Name() != "b" {
		t.Errorf("ordem incorreta: %s, %s", all[0].Name(), all[1].Name())
	}
}

func TestRegistry_RegisterDuplicate(t *testing.T) {
	reg := NewRegistry()

	_ = reg.Register(&stubModule{name: "a"})
	err := reg.Register(&stubModule{name: "a"})

	if err == nil {
		t.Fatal("esperava erro ao registrar modulo duplicado")
	}
}

func TestRegistry_ByName(t *testing.T) {
	reg := NewRegistry()
	_ = reg.Register(&stubModule{name: "starship"})

	mod, ok := reg.ByName("starship")
	if !ok {
		t.Fatal("modulo starship deveria existir")
	}
	if mod.Name() != "starship" {
		t.Errorf("esperava starship, obteve %s", mod.Name())
	}
}

func TestRegistry_ByName_NotFound(t *testing.T) {
	reg := NewRegistry()

	_, ok := reg.ByName("inexistente")
	if ok {
		t.Fatal("modulo inexistente nao deveria ser encontrado")
	}
}
