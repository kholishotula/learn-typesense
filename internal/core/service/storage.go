package service

import (
	"context"
	"learn-typesense/internal/core/entity"
)

type BookStorage interface {
	IndexBook(ctx context.Context, book entity.Book) (*entity.Book, error)
	GetBookById(ctx context.Context, id string) (*entity.Book, error)
	SearchBooks(ctx context.Context) (*[]entity.Book, error)
}
