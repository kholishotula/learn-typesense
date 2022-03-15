package storage

import (
	"context"
	"fmt"
	"learn-typesense/internal/core/entity"
	"learn-typesense/internal/core/ts"

	"gopkg.in/validator.v2"
)

type Storage struct {
	tsStorage        ts.Storage
	tsCollectionName string
}

type StorageConfig struct {
	TSStorage        ts.Storage `validate:"nonnil"`
	TSCollectionName string     `validate:"nonnil"`
}

func (c StorageConfig) Validate() error {
	return validator.Validate(c)
}

func NewStorage(cfg StorageConfig) (*Storage, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	return &Storage{
		tsStorage:        cfg.TSStorage,
		tsCollectionName: cfg.TSCollectionName,
	}, nil
}

func (s *Storage) IndexBook(ctx context.Context, book entity.Book) (*entity.Book, error) {
	err := s.tsStorage.Index(ctx, ts.IndexInput{
		CollectionName: s.tsCollectionName,
		Doc:            &book,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to index book due: %w", err)
	}
	return &book, nil
}

func (s *Storage) GetBookById(ctx context.Context, id string) (*entity.Book, error) {
	res, err := s.tsStorage.RetrieveById(ctx, ts.RetrieveInput{
		CollectionName: s.tsCollectionName,
		ID:             id,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to get book due: %w", err)
	}

	var book entity.Book
	bookStr := res.Doc.EncodeJSON()
	err = book.DecodeJSON(bookStr)
	if err != nil {
		return nil, fmt.Errorf("unable to decode book due: %w", err)
	}
	return &book, nil
}

func (s *Storage) SearchBooks(ctx context.Context) (*[]entity.Book, error) {
	res, err := s.tsStorage.Search(ctx, ts.SearchInput{
		CollectionName: s.tsCollectionName,
		Query:          "*",
		QueryBy:        "title",
		Page:           1,
		Limit:          10,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to search books due: %w", err)
	}

	var books []entity.Book
	if res.Found > 0 {
		for _, doc := range res.Hits {
			var book entity.Book
			bookStr := doc.EncodeJSON()
			err = book.DecodeJSON(bookStr)
			if err != nil {
				continue
			}
			books = append(books, book)
		}
	}

	return &books, nil
}
