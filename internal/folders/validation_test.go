package folders

import "testing"

func TestValidateNameAcceptsNormalName(t *testing.T) {
	if err := ValidateName("Documents"); err != nil {
		t.Fatalf("expected valid folder name, got %v", err)
	}
}

func TestValidateNameRejectsEmptyName(t *testing.T) {
	if err := ValidateName(""); err == nil {
		t.Fatal("expected empty folder name to fail")
	}
}

func TestValidateNameRejectsTraversalNames(t *testing.T) {
	names := []string{
		".",
		"..",
		"../secret",
		"folder/subfolder",
		`folder\subfolder`,
	}

	for _, name := range names {
		if err := ValidateName(name); err == nil {
			t.Fatalf("expected %q to fail", name)
		}
	}
}

func TestValidateNameRejectsLongName(t *testing.T) {
	name := ""

	for i := 0; i < 256; i++ {
		name += "a"
	}

	if err := ValidateName(name); err == nil {
		t.Fatal("expected oversized folder name to fail")
	}
}
