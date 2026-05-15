package ancillaries

import (
	"os"
	"path/filepath"
	"testing"

	"goweb/ancillaries"
)

func TestGetEndpoint_IndexFile(t *testing.T) {
	tmpDir := t.TempDir()

	indexFile := filepath.Join(tmpDir, "index_templ.go")
	if err := os.WriteFile(indexFile, []byte("package index"), 0644); err != nil {
		t.Fatalf("failed to create index file: %v", err)
	}

	endpoints := ancillaries.GetEndpoint(tmpDir)

	if len(endpoints) != 1 {
		t.Errorf("expected 1 endpoint, got %d: %v", len(endpoints), endpoints)
	}

	if endpoints[0] != "/" {
		t.Errorf("expected index to map to '/', got: %s", endpoints[0])
	}
}

func TestGetEndpoint_NonIndexFile(t *testing.T) {
	tmpDir := t.TempDir()

	aboutFile := filepath.Join(tmpDir, "about_templ.go")
	if err := os.WriteFile(aboutFile, []byte("package about"), 0644); err != nil {
		t.Fatalf("failed to create about file: %v", err)
	}

	endpoints := ancillaries.GetEndpoint(tmpDir)

	if len(endpoints) != 1 {
		t.Errorf("expected 1 endpoint, got %d: %v", len(endpoints), endpoints)
	}

	if endpoints[0] != "about" {
		t.Errorf("expected endpoint 'about', got: %s", endpoints[0])
	}
}

func TestGetEndpoint_Subdirectory(t *testing.T) {
	tmpDir := t.TempDir()

	userDir := filepath.Join(tmpDir, "user")
	if err := os.MkdirAll(userDir, 0755); err != nil {
		t.Fatalf("failed to create user dir: %v", err)
	}

	profileFile := filepath.Join(userDir, "profile_templ.go")
	if err := os.WriteFile(profileFile, []byte("package profile"), 0644); err != nil {
		t.Fatalf("failed to create profile file: %v", err)
	}

	endpoints := ancillaries.GetEndpoint(tmpDir)

	if len(endpoints) != 1 {
		t.Errorf("expected 1 endpoint, got %d: %v", len(endpoints), endpoints)
	}

	if endpoints[0] != "user/profile" {
		t.Errorf("expected endpoint 'user/profile', got: %s", endpoints[0])
	}
}

func TestGetEndpoint_IgnoresNonTemplFiles(t *testing.T) {
	tmpDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(tmpDir, "index_templ.go"), []byte("package index"), 0644); err != nil {
		t.Fatalf("failed to create index file: %v", err)
	}

	if err := os.WriteFile(filepath.Join(tmpDir, "readme.txt"), []byte("readme"), 0644); err != nil {
		t.Fatalf("failed to create readme file: %v", err)
	}

	endpoints := ancillaries.GetEndpoint(tmpDir)

	if len(endpoints) != 1 {
		t.Errorf("expected 1 endpoint, got %d: %v", len(endpoints), endpoints)
	}
}

func TestGetEndpoint_MultipleFiles(t *testing.T) {
	tmpDir := t.TempDir()

	files := []string{"index_templ.go", "about_templ.go", "contact_templ.go"}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, f), []byte("package "+f), 0644); err != nil {
			t.Fatalf("failed to create file: %v", err)
		}
	}

	endpoints := ancillaries.GetEndpoint(tmpDir)

	if len(endpoints) != 3 {
		t.Errorf("expected 3 endpoints, got %d: %v", len(endpoints), endpoints)
	}
}

func TestGetEndpoint_NonExistentDirectory(t *testing.T) {
	endpoints := ancillaries.GetEndpoint("/nonexistent/directory")

	if len(endpoints) != 0 {
		t.Errorf("expected 0 endpoints for non-existent dir, got %d", len(endpoints))
	}
}