package postgres

import (
	"context"
	"go-api-boilerplate/internal/domain"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
)

func TestPostgresBookRepo_CreateBook(t *testing.T) {
	tests := []struct {
		name    string
		book    domain.Book
		setup   func(pgxmock.PgxPoolIface)
		wantErr bool
	}{
		{
			name: "success",
			book: domain.Book{Title: "Test Book", Author: "Test Author"},
			setup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("INSERT INTO books").
					WithArgs("Test Book", "Test Author").
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			wantErr: false,
		},
		{
			name: "db error",
			book: domain.Book{Title: "Test Book", Author: "Test Author"},
			setup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("INSERT INTO books").
					WithArgs("Test Book", "Test Author").
					WillReturnError(pgx.ErrTxClosed)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mock.Close()

			tt.setup(mock)

			r := NewPostgresBookRepo(mock)
			if err := r.CreateBook(context.Background(), tt.book); (err != nil) != tt.wantErr {
				t.Errorf("PostgresBookRepo.CreateBook() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}

func TestPostgresBookRepo_GetBook(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		setup   func(pgxmock.PgxPoolIface)
		want    domain.Book
		wantErr bool
	}{
		{
			name: "success",
			id:   1,
			setup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "title", "author"}).
					AddRow(1, "Test Book", "Test Author")
				mock.ExpectQuery("SELECT id, title, author FROM books WHERE id").
					WithArgs(1).
					WillReturnRows(rows)
			},
			want:    domain.Book{ID: 1, Title: "Test Book", Author: "Test Author"},
			wantErr: false,
		},
		{
			name: "not found",
			id:   999,
			setup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT id, title, author FROM books WHERE id").
					WithArgs(999).
					WillReturnError(pgx.ErrNoRows)
			},
			want:    domain.Book{},
			wantErr: true,
		},
		{
			name: "scan error - other error",
			id:   1,
			setup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "title", "author"}).
					AddRow(1, "Test Book", "Test Author").
					RowError(0, pgx.ErrTxClosed)
				mock.ExpectQuery("SELECT id, title, author FROM books WHERE id").
					WithArgs(1).
					WillReturnRows(rows)
			},
			want:    domain.Book{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mock.Close()

			tt.setup(mock)

			r := NewPostgresBookRepo(mock)
			got, err := r.GetBook(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostgresBookRepo.GetBook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostgresBookRepo.GetBook() = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}

func TestPostgresBookRepo_GetBooks(t *testing.T) {
	tests := []struct {
		name    string
		offset  int
		limit   int
		setup   func(pgxmock.PgxPoolIface)
		want    []domain.Book
		wantErr bool
	}{
		{
			name:   "success with books",
			offset: 0,
			limit:  10,
			setup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "title", "author"}).
					AddRow(1, "Book 1", "Author 1").
					AddRow(2, "Book 2", "Author 2")
				mock.ExpectQuery("SELECT id, title, author FROM books ORDER BY id ASC").
					WithArgs(10, 0).
					WillReturnRows(rows)
			},
			want: []domain.Book{
				{ID: 1, Title: "Book 1", Author: "Author 1"},
				{ID: 2, Title: "Book 2", Author: "Author 2"},
			},
			wantErr: false,
		},
		{
			name:   "success with empty result",
			offset: 0,
			limit:  10,
			setup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "title", "author"})
				mock.ExpectQuery("SELECT id, title, author FROM books ORDER BY id ASC").
					WithArgs(10, 0).
					WillReturnRows(rows)
			},
			want:    []domain.Book{},
			wantErr: false,
		},
		{
			name:   "query error",
			offset: 0,
			limit:  10,
			setup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT id, title, author FROM books ORDER BY id ASC").
					WithArgs(10, 0).
					WillReturnError(pgx.ErrTxClosed)
			},
			want:    []domain.Book{},
			wantErr: true,
		},
		{
			name:   "scan error in rows",
			offset: 0,
			limit:  10,
			setup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "title", "author"}).
					AddRow(1, "Book 1", "Author 1").
					AddRow(2, "Book 2", "Author 2").
					RowError(1, pgx.ErrTxClosed)
				mock.ExpectQuery("SELECT id, title, author FROM books ORDER BY id ASC").
					WithArgs(10, 0).
					WillReturnRows(rows)
			},
			want:    []domain.Book{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mock.Close()

			tt.setup(mock)

			r := NewPostgresBookRepo(mock)
			got, err := r.GetBooks(context.Background(), tt.offset, tt.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostgresBookRepo.GetBooks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostgresBookRepo.GetBooks() = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}

func TestPostgresBookRepo_UpdateBook(t *testing.T) {
	tests := []struct {
		name    string
		book    domain.Book
		setup   func(pgxmock.PgxPoolIface)
		wantErr bool
	}{
		{
			name: "success",
			book: domain.Book{ID: 1, Title: "Updated Book", Author: "Updated Author"},
			setup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("UPDATE books SET title").
					WithArgs("Updated Book", "Updated Author", 1).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			wantErr: false,
		},
		{
			name: "not found - zero rows affected",
			book: domain.Book{ID: 999, Title: "Updated Book", Author: "Updated Author"},
			setup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("UPDATE books SET title").
					WithArgs("Updated Book", "Updated Author", 999).
					WillReturnResult(pgxmock.NewResult("UPDATE", 0))
			},
			wantErr: true,
		},
		{
			name: "db error",
			book: domain.Book{ID: 1, Title: "Updated Book", Author: "Updated Author"},
			setup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("UPDATE books SET title").
					WithArgs("Updated Book", "Updated Author", 1).
					WillReturnError(pgx.ErrTxClosed)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mock.Close()

			tt.setup(mock)

			r := NewPostgresBookRepo(mock)
			if err := r.UpdateBook(context.Background(), tt.book); (err != nil) != tt.wantErr {
				t.Errorf("PostgresBookRepo.UpdateBook() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}

func TestPostgresBookRepo_DeleteBook(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		setup   func(pgxmock.PgxPoolIface)
		wantErr bool
	}{
		{
			name: "success",
			id:   1,
			setup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("DELETE FROM books WHERE id").
					WithArgs(1).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
			wantErr: false,
		},
		{
			name: "not found - zero rows affected",
			id:   999,
			setup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("DELETE FROM books WHERE id").
					WithArgs(999).
					WillReturnResult(pgxmock.NewResult("DELETE", 0))
			},
			wantErr: true,
		},
		{
			name: "db error",
			id:   1,
			setup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("DELETE FROM books WHERE id").
					WithArgs(1).
					WillReturnError(pgx.ErrTxClosed)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mock.Close()

			tt.setup(mock)

			r := NewPostgresBookRepo(mock)
			if err := r.DeleteBook(context.Background(), tt.id); (err != nil) != tt.wantErr {
				t.Errorf("PostgresBookRepo.DeleteBook() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}
