package domain_test

import (
	"testing"

	"github.com/luizhanauer/proton-registry/internal/domain"
)

func TestReleaseCollection_IsEmpty(t *testing.T) {
	emptyCol := domain.ReleaseCollection{}
	if !emptyCol.IsEmpty() {
		t.Errorf("Expected collection to be empty")
	}

	filledCol := domain.ReleaseCollection{Releases: []domain.Release{{Version: "1.0"}}}
	if filledCol.IsEmpty() {
		t.Errorf("Expected collection to not be empty")
	}
}

func TestReleaseCollection_First(t *testing.T) {
	emptyCol := domain.ReleaseCollection{}
	firstEmpty := emptyCol.First()
	if firstEmpty.Version != "" {
		t.Errorf("Expected empty release, got %v", firstEmpty.Version)
	}

	filledCol := domain.ReleaseCollection{Releases: []domain.Release{{Version: "GE-Proton10-1"}}}
	firstFilled := filledCol.First()
	if firstFilled.Version != "GE-Proton10-1" {
		t.Errorf("Expected GE-Proton10-1, got %v", firstFilled.Version)
	}
}
