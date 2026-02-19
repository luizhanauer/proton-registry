package usecase

import (
	"fmt"

	"github.com/luizhanauer/proton-registry/internal/domain"
)

// Interfaces que definem os contratos esperados pelo UseCase
type Fetcher interface {
	GetLatestTagName() (string, error)
	FetchAll() (domain.ReleaseCollection, error)
}

type Storage interface {
	ReadIndex(filename string) (domain.ReleaseCollection, error)
	SaveIndex(filename string, collection domain.ReleaseCollection) error
}

type Filter interface {
	Apply(collection domain.ReleaseCollection) domain.ReleaseCollection
}

type Updater struct {
	fetcher Fetcher
	storage Storage
	filter  Filter
}

func NewUpdater(f Fetcher, s Storage, filter Filter) *Updater {
	return &Updater{fetcher: f, storage: s, filter: filter}
}

func (u *Updater) Execute(fullIndexPath, smartIndexPath string) error {
	if !u.needsUpdate(fullIndexPath) {
		fmt.Println("✅ O registro já está atualizado. Encerrando.")
		return nil
	}

	fmt.Println("🚀 Nova versão detectada! Atualizando índices...")

	fullCollection, err := u.fetcher.FetchAll()
	if err != nil {
		return fmt.Errorf("erro ao buscar releases: %w", err)
	}

	if err := u.storage.SaveIndex(fullIndexPath, fullCollection); err != nil {
		return fmt.Errorf("erro ao salvar full_index: %w", err)
	}

	smartCollection := u.filter.Apply(fullCollection)
	if err := u.storage.SaveIndex(smartIndexPath, smartCollection); err != nil {
		return fmt.Errorf("erro ao salvar smart_index: %w", err)
	}

	fmt.Println("🎉 Tudo pronto! Arquivos atualizados.")
	return nil
}

func (u *Updater) needsUpdate(localPath string) bool {
	remoteTag, err := u.fetcher.GetLatestTagName()
	if err != nil {
		fmt.Printf("⚠️ Erro ao checar latest: %v. Forçando update.\n", err)
		return true
	}

	localCollection, err := u.storage.ReadIndex(localPath)
	if err != nil || localCollection.IsEmpty() {
		return true
	}

	if remoteTag == localCollection.First().Version {
		fmt.Printf("⏸️ Versão atual (%s) é idêntica à remota.\n", remoteTag)
		return false
	}

	fmt.Printf("🆕 Nova versão encontrada: %s\n", remoteTag)
	return true
}
