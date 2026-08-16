package book

import (
	"context"

	"go-api-boilerplate/internal/adapter/grpc/grpckit"
	"go-api-boilerplate/internal/application/port/in"
	"go-api-boilerplate/internal/domain"
	"go-api-boilerplate/proto"
)

type BookServer struct {
	proto.UnimplementedBookServiceServer
	bookService in.BookUseCase
}

var _ proto.BookServiceServer = (*BookServer)(nil)

func NewBookServer(bookService in.BookUseCase) *BookServer {
	return &BookServer{bookService: bookService}
}

func (s *BookServer) CreateBook(ctx context.Context, req *proto.CreateBookRequest) (*proto.CreateBookResponse, error) {
	err := s.bookService.CreateBook(ctx, domain.Book{
		Title:  req.GetTitle(),
		Author: req.GetAuthor(),
	})
	if err != nil {
		return nil, err
	}
	return &proto.CreateBookResponse{}, nil
}

func (s *BookServer) GetBook(ctx context.Context, req *proto.GetBookRequest) (*proto.GetBookResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, grpckit.NewValidationError("id must be greater than 0")
	}

	book, err := s.bookService.GetBook(ctx, id)
	if err != nil {
		return nil, err
	}

	return &proto.GetBookResponse{Book: bookFromDomain(book)}, nil
}

func (s *BookServer) GetBooks(ctx context.Context, req *proto.GetBooksRequest) (*proto.GetBooksResponse, error) {
	page := int(req.GetPage())
	perPage := int(req.GetPerPage())

	if page < 0 {
		return nil, grpckit.NewValidationError("page must be greater than or equal to 1")
	}
	if perPage < 0 || perPage > 100 {
		return nil, grpckit.NewValidationError("per_page must be between 1 and 100")
	}

	books, err := s.bookService.GetBooks(ctx, page, perPage)
	if err != nil {
		return nil, err
	}

	items := make([]*proto.Book, 0, len(books))
	for _, book := range books {
		items = append(items, bookFromDomain(book))
	}

	return &proto.GetBooksResponse{Books: items}, nil
}

func (s *BookServer) UpdateBook(ctx context.Context, req *proto.UpdateBookRequest) (*proto.UpdateBookResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, grpckit.NewValidationError("id must be greater than 0")
	}

	err := s.bookService.UpdateBook(ctx, domain.Book{
		ID:     id,
		Title:  req.GetTitle(),
		Author: req.GetAuthor(),
	})
	if err != nil {
		return nil, err
	}

	return &proto.UpdateBookResponse{}, nil
}

func (s *BookServer) DeleteBook(ctx context.Context, req *proto.DeleteBookRequest) (*proto.DeleteBookResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, grpckit.NewValidationError("id must be greater than 0")
	}

	err := s.bookService.DeleteBook(ctx, id)
	if err != nil {
		return nil, err
	}

	return &proto.DeleteBookResponse{}, nil
}

func bookFromDomain(book domain.Book) *proto.Book {
	return &proto.Book{
		Id:     int32(book.ID),
		Title:  book.Title,
		Author: book.Author,
	}
}
