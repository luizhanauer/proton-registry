package domain_test

import (
	"reflect"
	"testing"

	"github.com/luizhanauer/proton-registry/internal/domain"
)

func TestSmartFilter_Apply(t *testing.T) {
	filter := domain.NewSmartFilter()

	tests := []struct {
		name     string
		input    domain.ReleaseCollection
		expected domain.ReleaseCollection
	}{
		{
			name:     "Should return empty collection when input is empty",
			input:    domain.ReleaseCollection{Releases: []domain.Release{}},
			expected: domain.ReleaseCollection{Releases: []domain.Release{}},
		},
		{
			name: "Should keep all when less than 10 releases",
			input: domain.ReleaseCollection{Releases: []domain.Release{
				{Version: "GE-Proton10-3", Major: "Proton10"},
				{Version: "GE-Proton10-2", Major: "Proton10"},
			}},
			expected: domain.ReleaseCollection{Releases: []domain.Release{
				{Version: "GE-Proton10-3", Major: "Proton10"},
				{Version: "GE-Proton10-2", Major: "Proton10"},
			}},
		},
		{
			name: "Should keep top 10 bleeding edge and only one of each legacy major",
			input: domain.ReleaseCollection{Releases: []domain.Release{
				// Top 10 (Bleeding Edge) - Todos devem ser mantidos
				{Version: "GE-Proton10-10", Major: "Proton10"},
				{Version: "GE-Proton10-9", Major: "Proton10"},
				{Version: "GE-Proton10-8", Major: "Proton10"},
				{Version: "GE-Proton10-7", Major: "Proton10"},
				{Version: "GE-Proton10-6", Major: "Proton10"},
				{Version: "GE-Proton10-5", Major: "Proton10"},
				{Version: "GE-Proton10-4", Major: "Proton10"},
				{Version: "GE-Proton10-3", Major: "Proton10"},
				{Version: "GE-Proton10-2", Major: "Proton10"},
				{Version: "GE-Proton10-1", Major: "Proton10"},

				// Legacy (A partir do 11º item)
				// Proton9 não foi visto no top 10, então mantém o primeiro
				{Version: "GE-Proton9-20", Major: "Proton9"},
				// Este Proton9 deve ser descartado pois já pegamos o mais recente
				{Version: "GE-Proton9-19", Major: "Proton9"},
				// Proton8 não foi visto, mantém o primeiro
				{Version: "GE-Proton8-30", Major: "Proton8"},
				// Este Proton8 deve ser descartado
				{Version: "GE-Proton8-29", Major: "Proton8"},
				// "Outros" ou vazio sempre devem ser descartados no legacy
				{Version: "Custom-Build", Major: "Outros"},
			}},
			expected: domain.ReleaseCollection{Releases: []domain.Release{
				{Version: "GE-Proton10-10", Major: "Proton10"},
				{Version: "GE-Proton10-9", Major: "Proton10"},
				{Version: "GE-Proton10-8", Major: "Proton10"},
				{Version: "GE-Proton10-7", Major: "Proton10"},
				{Version: "GE-Proton10-6", Major: "Proton10"},
				{Version: "GE-Proton10-5", Major: "Proton10"},
				{Version: "GE-Proton10-4", Major: "Proton10"},
				{Version: "GE-Proton10-3", Major: "Proton10"},
				{Version: "GE-Proton10-2", Major: "Proton10"},
				{Version: "GE-Proton10-1", Major: "Proton10"},
				{Version: "GE-Proton9-20", Major: "Proton9"},
				{Version: "GE-Proton8-30", Major: "Proton8"},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filter.Apply(tt.input)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("\nExpected: %+v\nGot: %+v", tt.expected, result)
			}
		})
	}
}
