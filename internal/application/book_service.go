package application

import (
	"context"
	"go-api-boilerplate/internal/application/port/in"
	"go-api-boilerplate/internal/application/port/out"
	"go-api-boilerplate/internal/domain"
)

type BookService struct {
	bookRepo out.BookRepository
}

var _ in.BookUseCase = &BookService{}

func NewBookService(bookRepo out.BookRepository) *BookService {
	return &BookService{bookRepo: bookRepo}
}

func (s *BookService) CreateBook(ctx context.Context, book domain.Book) error {
	if err := book.Validate(); err != nil {
		return err
	}
	return s.bookRepo.CreateBook(ctx, book)
}

func (s *BookService) GetBook(ctx context.Context, id int) (domain.Book, error) {
	return s.bookRepo.GetBook(ctx, id)
}

func (s *BookService) GetBooks(ctx context.Context, page, perPage int) ([]domain.Book, error) {
	if page < 1 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}

	offset := (page - 1) * perPage
	limit := perPage

	return s.bookRepo.GetBooks(ctx, offset, limit)
}

func (s *BookService) UpdateBook(ctx context.Context, book domain.Book) error {
	if err := book.Validate(); err != nil {
		return err
	}
	return s.bookRepo.UpdateBook(ctx, book)
}

func (s *BookService) DeleteBook(ctx context.Context, id int) error {
	return s.bookRepo.DeleteBook(ctx, id)
}
