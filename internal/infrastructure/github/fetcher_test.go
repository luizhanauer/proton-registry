package github_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/luizhanauer/proton-registry/internal/infrastructure/github"
)

func TestFetcher_GetLatestTagName(t *testing.T) {
	// Cria um servidor HTTP falso
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/releases/latest" {
			w.WriteHeader(http.StatusOK)
			// Retorna um JSON simulando a API do GitHub
			fmt.Fprintln(w, `{"tag_name": "GE-Proton10-32"}`)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	// Injeta a URL do servidor falso no nosso Fetcher
	fetcher := github.NewFetcher(server.URL)

	tag, err := fetcher.GetLatestTagName()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if tag != "GE-Proton10-32" {
		t.Errorf("Expected tag GE-Proton10-32, got %v", tag)
	}
}

func TestFetcher_FetchAll(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simula paginação: Página 1 tem dados, Página 2 é vazia (fim do loop)
		if r.URL.Query().Get("page") == "1" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `[
				{
					"tag_name": "GE-Proton10-32",
					"published_at": "2026-02-16T15:00:00Z",
					"assets": [
						{"name": "GE-Proton10-32.tar.gz", "browser_download_url": "http://url", "size": 100}
					]
				}
			]`)
			return
		}

		// Página 2 (vazia) encerra o loop de paginação no nosso Fetcher
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `[]`)
	}))
	defer server.Close()

	fetcher := github.NewFetcher(server.URL)

	collection, err := fetcher.FetchAll()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if collection.IsEmpty() {
		t.Fatalf("Expected collection to have items")
	}

	first := collection.First()
	if first.Version != "GE-Proton10-32" || first.Major != "Proton10" {
		t.Errorf("Unexpected release parsing result: %+v", first)
	}
}
