package usecase_test

import (
	"testing"

	"github.com/luizhanauer/proton-registry/internal/domain"
	"github.com/luizhanauer/proton-registry/internal/usecase"
)

// --- MOCKS ---
type MockFetcher struct {
	LatestTag string
	Releases  domain.ReleaseCollection
	Err       error
}

func (m *MockFetcher) GetLatestTagName() (string, error)           { return m.LatestTag, m.Err }
func (m *MockFetcher) FetchAll() (domain.ReleaseCollection, error) { return m.Releases, m.Err }

type MockStorage struct {
	LocalCollection domain.ReleaseCollection
	SaveError       error
	ReadError       error
}

func (m *MockStorage) ReadIndex(name string) (domain.ReleaseCollection, error) {
	return m.LocalCollection, m.ReadError
}
func (m *MockStorage) SaveIndex(name string, col domain.ReleaseCollection) error { return m.SaveError }

type MockFilter struct{}

func (m *MockFilter) Apply(col domain.ReleaseCollection) domain.ReleaseCollection { return col }

// --- TESTES ---
func TestUpdater_Execute_NoUpdateNeeded(t *testing.T) {
	col := domain.ReleaseCollection{Releases: []domain.Release{{Version: "GE-Proton10-1"}}}

	fetcher := &MockFetcher{LatestTag: "GE-Proton10-1"}
	storage := &MockStorage{LocalCollection: col}

	updater := usecase.NewUpdater(fetcher, storage, &MockFilter{})

	err := updater.Execute("full.json", "smart.json")
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
}

func TestUpdater_Execute_RequiresUpdate(t *testing.T) {
	remoteCol := domain.ReleaseCollection{Releases: []domain.Release{{Version: "GE-Proton10-2"}}}
	localCol := domain.ReleaseCollection{Releases: []domain.Release{{Version: "GE-Proton10-1"}}}

	fetcher := &MockFetcher{LatestTag: "GE-Proton10-2", Releases: remoteCol}
	storage := &MockStorage{LocalCollection: localCol}

	updater := usecase.NewUpdater(fetcher, storage, &MockFilter{})

	err := updater.Execute("full.json", "smart.json")
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
}
