package in

import (
	"context"
	"go-api-boilerplate/internal/domain"
)

type BookUseCase interface {
	CreateBook(ctx context.Context, book domain.Book) error
	GetBook(ctx context.Context, id int) (domain.Book, error)
	GetBooks(ctx context.Context, page, perPage int) ([]domain.Book, error)
	UpdateBook(ctx context.Context, book domain.Book) error
	DeleteBook(ctx context.Context, id int) error
}
