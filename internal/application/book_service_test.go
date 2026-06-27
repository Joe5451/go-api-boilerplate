package application

import (
	"context"
	"errors"
	"go-api-boilerplate/internal/domain"
	"go-api-boilerplate/mocks"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestBookService_CreateBook(t *testing.T) {
	type args struct {
		ctx  context.Context
		book domain.Book
	}
	tests := []struct {
		name    string
		args    args
		setup   func(*mocks.MockBookRepository)
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:  context.Background(),
				book: domain.Book{Title: "Test Book", Author: "Test Author"},
			},
			setup: func(m *mocks.MockBookRepository) {
				m.EXPECT().CreateBook(gomock.Any(), domain.Book{Title: "Test Book", Author: "Test Author"}).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "validation error - empty title",
			args: args{
				ctx:  context.Background(),
				book: domain.Book{Title: "", Author: "Test Author"},
			},
			setup:   func(m *mocks.MockBookRepository) {},
			wantErr: true,
		},
		{
			name: "validation error - empty author",
			args: args{
				ctx:  context.Background(),
				book: domain.Book{Title: "Test Book", Author: ""},
			},
			setup:   func(m *mocks.MockBookRepository) {},
			wantErr: true,
		},
		{
			name: "repository error",
			args: args{
				ctx:  context.Background(),
				book: domain.Book{Title: "Test Book", Author: "Test Author"},
			},
			setup: func(m *mocks.MockBookRepository) {
				m.EXPECT().CreateBook(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockBookRepository(ctrl)
			tt.setup(mockRepo)

			s := NewBookService(mockRepo)
			if err := s.CreateBook(tt.args.ctx, tt.args.book); (err != nil) != tt.wantErr {
				t.Errorf("BookService.CreateBook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBookService_GetBook(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name    string
		args    args
		setup   func(*mocks.MockBookRepository)
		want    domain.Book
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			setup: func(m *mocks.MockBookRepository) {
				m.EXPECT().GetBook(gomock.Any(), 1).Return(domain.Book{ID: 1, Title: "Test Book", Author: "Test Author"}, nil)
			},
			want:    domain.Book{ID: 1, Title: "Test Book", Author: "Test Author"},
			wantErr: false,
		},
		{
			name: "not found",
			args: args{
				ctx: context.Background(),
				id:  999,
			},
			setup: func(m *mocks.MockBookRepository) {
				m.EXPECT().GetBook(gomock.Any(), 999).Return(domain.Book{}, domain.ErrBookNotFound)
			},
			want:    domain.Book{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockBookRepository(ctrl)
			tt.setup(mockRepo)

			s := NewBookService(mockRepo)
			got, err := s.GetBook(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("BookService.GetBook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BookService.GetBook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBookService_GetBooks(t *testing.T) {
	type args struct {
		ctx     context.Context
		page    int
		perPage int
	}
	tests := []struct {
		name    string
		args    args
		setup   func(*mocks.MockBookRepository)
		want    []domain.Book
		wantErr bool
	}{
		{
			name: "success - first page",
			args: args{
				ctx:     context.Background(),
				page:    1,
				perPage: 10,
			},
			setup: func(m *mocks.MockBookRepository) {
				m.EXPECT().GetBooks(gomock.Any(), 0, 10).Return([]domain.Book{
					{ID: 1, Title: "Book 1", Author: "Author 1"},
					{ID: 2, Title: "Book 2", Author: "Author 2"},
				}, nil)
			},
			want: []domain.Book{
				{ID: 1, Title: "Book 1", Author: "Author 1"},
				{ID: 2, Title: "Book 2", Author: "Author 2"},
			},
			wantErr: false,
		},
		{
			name: "success - second page",
			args: args{
				ctx:     context.Background(),
				page:    2,
				perPage: 10,
			},
			setup: func(m *mocks.MockBookRepository) {
				m.EXPECT().GetBooks(gomock.Any(), 10, 10).Return([]domain.Book{}, nil)
			},
			want:    []domain.Book{},
			wantErr: false,
		},
		{
			name: "default pagination - invalid page",
			args: args{
				ctx:     context.Background(),
				page:    0,
				perPage: 10,
			},
			setup: func(m *mocks.MockBookRepository) {
				m.EXPECT().GetBooks(gomock.Any(), 0, 10).Return([]domain.Book{}, nil)
			},
			want:    []domain.Book{},
			wantErr: false,
		},
		{
			name: "default pagination - invalid perPage",
			args: args{
				ctx:     context.Background(),
				page:    1,
				perPage: 0,
			},
			setup: func(m *mocks.MockBookRepository) {
				m.EXPECT().GetBooks(gomock.Any(), 0, 10).Return([]domain.Book{}, nil)
			},
			want:    []domain.Book{},
			wantErr: false,
		},
		{
			name: "repository error",
			args: args{
				ctx:     context.Background(),
				page:    1,
				perPage: 10,
			},
			setup: func(m *mocks.MockBookRepository) {
				m.EXPECT().GetBooks(gomock.Any(), 0, 10).Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockBookRepository(ctrl)
			tt.setup(mockRepo)

			s := NewBookService(mockRepo)
			got, err := s.GetBooks(tt.args.ctx, tt.args.page, tt.args.perPage)
			if (err != nil) != tt.wantErr {
				t.Errorf("BookService.GetBooks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BookService.GetBooks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBookService_UpdateBook(t *testing.T) {
	type args struct {
		ctx  context.Context
		book domain.Book
	}
	tests := []struct {
		name    string
		args    args
		setup   func(*mocks.MockBookRepository)
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:  context.Background(),
				book: domain.Book{ID: 1, Title: "Updated Book", Author: "Updated Author"},
			},
			setup: func(m *mocks.MockBookRepository) {
				m.EXPECT().UpdateBook(gomock.Any(), domain.Book{ID: 1, Title: "Updated Book", Author: "Updated Author"}).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "validation error - empty title",
			args: args{
				ctx:  context.Background(),
				book: domain.Book{ID: 1, Title: "", Author: "Updated Author"},
			},
			setup:   func(m *mocks.MockBookRepository) {},
			wantErr: true,
		},
		{
			name: "validation error - empty author",
			args: args{
				ctx:  context.Background(),
				book: domain.Book{ID: 1, Title: "Updated Book", Author: ""},
			},
			setup:   func(m *mocks.MockBookRepository) {},
			wantErr: true,
		},
		{
			name: "repository error",
			args: args{
				ctx:  context.Background(),
				book: domain.Book{ID: 1, Title: "Updated Book", Author: "Updated Author"},
			},
			setup: func(m *mocks.MockBookRepository) {
				m.EXPECT().UpdateBook(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockBookRepository(ctrl)
			tt.setup(mockRepo)

			s := NewBookService(mockRepo)
			if err := s.UpdateBook(tt.args.ctx, tt.args.book); (err != nil) != tt.wantErr {
				t.Errorf("BookService.UpdateBook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBookService_DeleteBook(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name    string
		args    args
		setup   func(*mocks.MockBookRepository)
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			setup: func(m *mocks.MockBookRepository) {
				m.EXPECT().DeleteBook(gomock.Any(), 1).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "not found",
			args: args{
				ctx: context.Background(),
				id:  999,
			},
			setup: func(m *mocks.MockBookRepository) {
				m.EXPECT().DeleteBook(gomock.Any(), 999).Return(domain.ErrBookNotFound)
			},
			wantErr: true,
		},
		{
			name: "repository error",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			setup: func(m *mocks.MockBookRepository) {
				m.EXPECT().DeleteBook(gomock.Any(), 1).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockBookRepository(ctrl)
			tt.setup(mockRepo)

			s := NewBookService(mockRepo)
			if err := s.DeleteBook(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("BookService.DeleteBook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
