package storage_test

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/luizhanauer/proton-registry/internal/domain"
	"github.com/luizhanauer/proton-registry/internal/infrastructure/storage"
)

func TestFileStorage_SaveAndReadIndex(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test_index.json")

	store := storage.NewFileStorage()

	original := domain.ReleaseCollection{
		Releases: []domain.Release{
			{Version: "GE-Proton10-32", DownloadURL: "http://example.com/file.tar.gz", Size: 1024, Date: "2026-02-16", Major: "Proton10"},
		},
	}

	// Testa a Escrita
	err := store.SaveIndex(filePath, original)
	if err != nil {
		t.Fatalf("Failed to save index: %v", err)
	}

	// Testa a Leitura
	readCollection, err := store.ReadIndex(filePath)
	if err != nil {
		t.Fatalf("Failed to read index: %v", err)
	}

	if !reflect.DeepEqual(original, readCollection) {
		t.Errorf("Expected %+v, got %+v", original, readCollection)
	}
}

func TestFileStorage_ReadIndex_FileNotFound(t *testing.T) {
	store := storage.NewFileStorage()
	_, err := store.ReadIndex("non_existent_file.json")
	if err == nil {
		t.Errorf("Expected error for non existent file, got nil")
	}
}
