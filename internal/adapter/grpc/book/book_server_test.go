package book

import (
	"context"
	"errors"
	"testing"

	"go-api-boilerplate/internal/adapter/grpc/grpckit"
	"go-api-boilerplate/internal/domain"
	"go-api-boilerplate/mocks"
	"go-api-boilerplate/proto"

	"go.uber.org/mock/gomock"
)

func TestBookServer_CreateBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	uc := mocks.NewMockBookUseCase(ctrl)
	s := NewBookServer(uc)

	uc.EXPECT().CreateBook(gomock.Any(), domain.Book{Title: "1984", Author: "George"}).Return(nil)

	_, err := s.CreateBook(context.Background(), &proto.CreateBookRequest{
		Title:  "1984",
		Author: "George",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBookServer_CreateBook_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	uc := mocks.NewMockBookUseCase(ctrl)
	s := NewBookServer(uc)

	uc.EXPECT().CreateBook(gomock.Any(), gomock.Any()).Return(domain.ErrTitleRequired)

	_, err := s.CreateBook(context.Background(), &proto.CreateBookRequest{
		Title:  "",
		Author: "George",
	})
	if !errors.Is(err, domain.ErrTitleRequired) {
		t.Fatalf("expected ErrTitleRequired, got %v", err)
	}
}

func TestBookServer_GetBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	uc := mocks.NewMockBookUseCase(ctrl)
	s := NewBookServer(uc)

	uc.EXPECT().GetBook(gomock.Any(), 1).Return(domain.Book{ID: 1, Title: "1984", Author: "George"}, nil)

	resp, err := s.GetBook(context.Background(), &proto.GetBookRequest{Id: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.GetBook().GetTitle() != "1984" {
		t.Fatalf("unexpected title: %s", resp.GetBook().GetTitle())
	}
}

func TestBookServer_GetBook_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	s := NewBookServer(mocks.NewMockBookUseCase(ctrl))

	_, err := s.GetBook(context.Background(), &proto.GetBookRequest{Id: 0})
	var grpcErr *grpckit.Error
	if !errors.As(err, &grpcErr) {
		t.Fatalf("expected grpckit.Error, got %v", err)
	}
}

func TestBookServer_GetBook_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	uc := mocks.NewMockBookUseCase(ctrl)
	s := NewBookServer(uc)

	uc.EXPECT().GetBook(gomock.Any(), 99).Return(domain.Book{}, domain.ErrBookNotFound)

	_, err := s.GetBook(context.Background(), &proto.GetBookRequest{Id: 99})
	if !errors.Is(err, domain.ErrBookNotFound) {
		t.Fatalf("expected ErrBookNotFound, got %v", err)
	}
}

func TestBookServer_GetBooks(t *testing.T) {
	ctrl := gomock.NewController(t)
	uc := mocks.NewMockBookUseCase(ctrl)
	s := NewBookServer(uc)

	uc.EXPECT().GetBooks(gomock.Any(), 1, 10).Return([]domain.Book{
		{ID: 1, Title: "1984", Author: "George"},
	}, nil)

	resp, err := s.GetBooks(context.Background(), &proto.GetBooksRequest{Page: 1, PerPage: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.GetBooks()) != 1 {
		t.Fatalf("expected 1 book, got %d", len(resp.GetBooks()))
	}
}

func TestBookServer_GetBooks_InvalidPerPage(t *testing.T) {
	ctrl := gomock.NewController(t)
	s := NewBookServer(mocks.NewMockBookUseCase(ctrl))

	_, err := s.GetBooks(context.Background(), &proto.GetBooksRequest{Page: 1, PerPage: 101})
	var grpcErr *grpckit.Error
	if !errors.As(err, &grpcErr) {
		t.Fatalf("expected grpckit.Error, got %v", err)
	}
}

func TestBookServer_UpdateBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	uc := mocks.NewMockBookUseCase(ctrl)
	s := NewBookServer(uc)

	uc.EXPECT().UpdateBook(gomock.Any(), domain.Book{ID: 1, Title: "New", Author: "Author"}).Return(nil)

	_, err := s.UpdateBook(context.Background(), &proto.UpdateBookRequest{
		Id:     1,
		Title:  "New",
		Author: "Author",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBookServer_UpdateBook_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	uc := mocks.NewMockBookUseCase(ctrl)
	s := NewBookServer(uc)

	uc.EXPECT().UpdateBook(gomock.Any(), gomock.Any()).Return(domain.ErrBookNotFound)

	_, err := s.UpdateBook(context.Background(), &proto.UpdateBookRequest{
		Id:     99,
		Title:  "New",
		Author: "Author",
	})
	if !errors.Is(err, domain.ErrBookNotFound) {
		t.Fatalf("expected ErrBookNotFound, got %v", err)
	}
}

func TestBookServer_DeleteBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	uc := mocks.NewMockBookUseCase(ctrl)
	s := NewBookServer(uc)

	uc.EXPECT().DeleteBook(gomock.Any(), 1).Return(nil)

	_, err := s.DeleteBook(context.Background(), &proto.DeleteBookRequest{Id: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBookServer_DeleteBook_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	uc := mocks.NewMockBookUseCase(ctrl)
	s := NewBookServer(uc)

	uc.EXPECT().DeleteBook(gomock.Any(), 99).Return(domain.ErrBookNotFound)

	_, err := s.DeleteBook(context.Background(), &proto.DeleteBookRequest{Id: 99})
	if !errors.Is(err, domain.ErrBookNotFound) {
		t.Fatalf("expected ErrBookNotFound, got %v", err)
	}
}
