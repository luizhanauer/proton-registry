package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/luizhanauer/proton-registry/internal/domain"
	"github.com/luizhanauer/proton-registry/internal/infrastructure/github"
	"github.com/luizhanauer/proton-registry/internal/infrastructure/storage"
	"github.com/luizhanauer/proton-registry/internal/usecase"
)

const (
	githubAPIBase  = "https://api.github.com/repos/GloriousEggroll/proton-ge-custom"
	outputDir      = "public/api"
	fullIndexFile  = "full_index.json"
	smartIndexFile = "smart_index.json"
)

func main() {
	fmt.Println("🤖 Iniciando Proton Registry...")

	// Garante que a pasta public/api/ existe antes de rodar
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("❌ Erro ao criar diretório de saída: %v\n", err)
		os.Exit(1)
	}

	fetcher := github.NewFetcher(githubAPIBase)
	store := storage.NewFileStorage()
	filter := domain.NewSmartFilter()

	updater := usecase.NewUpdater(fetcher, store, filter)

	// Monta os caminhos finais
	fullPath := filepath.Join(outputDir, fullIndexFile)
	smartPath := filepath.Join(outputDir, smartIndexFile)

	if err := updater.Execute(fullPath, smartPath); err != nil {
		fmt.Printf("❌ Erro na execução: %v\n", err)
		os.Exit(1)
	}
}
