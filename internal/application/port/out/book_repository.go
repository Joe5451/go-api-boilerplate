package out

import (
	"context"
	"go-api-boilerplate/internal/domain"
)

type BookRepository interface {
	CreateBook(ctx context.Context, book domain.Book) error
	GetBook(ctx context.Context, id int) (domain.Book, error)
	GetBooks(ctx context.Context, offset, limit int) ([]domain.Book, error)
	UpdateBook(ctx context.Context, book domain.Book) error
	DeleteBook(ctx context.Context, id int) error
}
