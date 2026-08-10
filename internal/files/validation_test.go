package files

import "testing"

func TestValidateFilenameAcceptsNormalFilename(t *testing.T) {
	if err := ValidateFilename("report.pdf"); err != nil {
		t.Fatalf("expected valid filename, got %v", err)
	}
}

func TestValidateFilenameRejectsEmptyFilename(t *testing.T) {
	if err := ValidateFilename(""); err == nil {
		t.Fatal("expected empty filename to fail")
	}
}

func TestValidateFilenameRejectsTraversal(t *testing.T) {
	names := []string{
		"../secret.txt",
		"../../secret.txt",
		"folder/file.txt",
		`folder\file.txt`,
	}

	for _, name := range names {
		if err := ValidateFilename(name); err == nil {
			t.Fatalf("expected %q to fail validation", name)
		}
	}
}
