package postgres

import (
	"context"
	"go-api-boilerplate/internal/application/port/out"
	"go-api-boilerplate/internal/domain"

	pgx "github.com/jackc/pgx/v5"
)

type PostgresBookRepo struct {
	db PgxIface
}

var _ out.BookRepository = &PostgresBookRepo{}

func NewPostgresBookRepo(db PgxIface) *PostgresBookRepo {
	return &PostgresBookRepo{db: db}
}

func (r *PostgresBookRepo) CreateBook(ctx context.Context, book domain.Book) error {
	_, err := r.db.Exec(
		ctx,
		"INSERT INTO books (title, author) VALUES ($1, $2)",
		book.Title,
		book.Author,
	)
	return err
}

func (r *PostgresBookRepo) GetBook(ctx context.Context, id int) (domain.Book, error) {
	var book domain.Book
	err := r.db.QueryRow(
		ctx,
		"SELECT id, title, author FROM books WHERE id = $1",
		id,
	).Scan(&book.ID, &book.Title, &book.Author)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.Book{}, domain.ErrBookNotFound
		}
		return domain.Book{}, err
	}
	return book, nil
}

func (r *PostgresBookRepo) GetBooks(ctx context.Context, offset, limit int) ([]domain.Book, error) {
	rows, err := r.db.Query(
		ctx,
		"SELECT id, title, author FROM books ORDER BY id ASC LIMIT $1 OFFSET $2",
		limit,
		offset,
	)
	if err != nil {
		return []domain.Book{}, err
	}
	defer rows.Close()

	books := []domain.Book{}

	for rows.Next() {
		var book domain.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author)
		if err != nil {
			return []domain.Book{}, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (r *PostgresBookRepo) UpdateBook(ctx context.Context, book domain.Book) error {
	cmdTag, err := r.db.Exec(
		ctx,
		"UPDATE books SET title = $1, author = $2 WHERE id = $3",
		book.Title,
		book.Author,
		book.ID,
	)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return domain.ErrBookNotFound
	}
	return nil
}

func (r *PostgresBookRepo) DeleteBook(ctx context.Context, id int) error {
	cmdTag, err := r.db.Exec(
		ctx,
		"DELETE FROM books WHERE id = $1",
		id,
	)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return domain.ErrBookNotFound
	}
	return nil
}
