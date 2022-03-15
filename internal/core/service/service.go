package service

import (
	"context"
	"fmt"
	"learn-typesense/internal/core/entity"

	"github.com/google/uuid"
	"gopkg.in/validator.v2"
)

type Service interface {
	CreateBook(ctx context.Context, input entity.CreateInput) (*entity.Book, error)
	GetBookById(ctx context.Context, id string) (*entity.Book, error)
	SearchBooks(ctx context.Context) (*[]entity.Book, error)
	// UpdateBook(ctx context.Context, id int)
	// DeleteBook(ctx context.Context, id int)
}

type ServiceConfig struct {
	BookStorage BookStorage `validate:"nonnil"`
}

func (c ServiceConfig) Validate() error {
	return validator.Validate(c)
}

func NewService(cfg ServiceConfig) (Service, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	s := &service{
		bookStorage: cfg.BookStorage,
	}
	return s, nil
}

type service struct {
	bookStorage BookStorage
}

func (s *service) CreateBook(ctx context.Context, input entity.CreateInput) (*entity.Book, error) {
	book, err := s.bookStorage.IndexBook(ctx, entity.Book{
		ID:      uuid.New().String(),
		Title:   input.Title,
		Author:  input.Author,
		Year:    input.Year,
		Summary: input.Summary,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create book due: %w", err)
	}

	return book, nil
}

func (s *service) GetBookById(ctx context.Context, id string) (*entity.Book, error) {
	book, err := s.bookStorage.GetBookById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("unable to get book due: %w", err)
	}

	return book, nil
}

func (s *service) SearchBooks(ctx context.Context) (*[]entity.Book, error) {
	books, err := s.bookStorage.SearchBooks(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to search books due: %w", err)
	}

	return books, nil
}
