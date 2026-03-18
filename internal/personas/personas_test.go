package personas

import (
	"testing"
)

func TestLoadAllPersonas(t *testing.T) {
	all, err := LoadAll()
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}
	if len(all) != 11 {
		t.Errorf("expected 11 personas, got %d", len(all))
	}
}

func TestPersonaHasRequiredFields(t *testing.T) {
	all, err := LoadAll()
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}
	for _, persona := range all {
		if persona.ID == "" {
			t.Errorf("persona missing ID")
		}
		if persona.Name == "" {
			t.Errorf("persona %s missing Name", persona.ID)
		}
		if persona.Description == "" {
			t.Errorf("persona %s missing Description", persona.ID)
		}
		if len(persona.Triggers) == 0 {
			t.Errorf("persona %s missing Triggers", persona.ID)
		}
	}
}

func TestPersonaContent(t *testing.T) {
	all, err := LoadAll()
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}
	for _, persona := range all {
		if persona.Content == "" {
			t.Errorf("persona %s has empty Content", persona.ID)
		}
	}
}
