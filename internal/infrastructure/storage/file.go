package storage

import (
	"encoding/json"
	"os"

	"github.com/luizhanauer/proton-registry/internal/domain"
)

type FileStorage struct{}

func NewFileStorage() *FileStorage {
	return &FileStorage{}
}

func (s *FileStorage) ReadIndex(filename string) (domain.ReleaseCollection, error) {
	file, err := os.Open(filename)
	if err != nil {
		return domain.ReleaseCollection{}, err
	}
	defer file.Close()

	var releases []domain.Release
	if err := json.NewDecoder(file).Decode(&releases); err != nil {
		return domain.ReleaseCollection{}, err
	}

	return domain.ReleaseCollection{Releases: releases}, nil
}

func (s *FileStorage) SaveIndex(filename string, collection domain.ReleaseCollection) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(collection.Releases)
}
