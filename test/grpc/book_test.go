package grpc

import (
	"context"
	"os"
	"testing"

	"go-api-boilerplate/proto"
	"go-api-boilerplate/test"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMain(m *testing.M) {
	test.Setup()
	code := m.Run()
	test.Teardown()
	os.Exit(code)
}

func setupTestApp(t *testing.T) *test.GrpcApp {
	app, err := test.InitGrpcApp(context.Background())
	if err != nil {
		t.Fatalf("failed to initialize app: %v", err)
	}
	t.Cleanup(app.Close)
	return app
}

func TestBookAPI_CreateBook(t *testing.T) {
	app := setupTestApp(t)
	defer test.CleanupDatabase(t)

	tests := []struct {
		name     string
		title    string
		author   string
		wantCode codes.Code
		checkDB  bool
	}{
		{name: "success", title: "1984", author: "George", wantCode: codes.OK, checkDB: true},
		{name: "missing_title", title: "", author: "Author", wantCode: codes.InvalidArgument},
		{name: "missing_author", title: "Title", author: "", wantCode: codes.InvalidArgument},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := app.Client.CreateBook(context.Background(), &proto.CreateBookRequest{
				Title:  tt.title,
				Author: tt.author,
			})
			if status.Code(err) != tt.wantCode {
				t.Fatalf("expected %s, got %s: %v", tt.wantCode, status.Code(err), err)
			}

			if tt.checkDB {
				var title, author string
				err := test.DB().QueryRow(
					context.Background(),
					"SELECT title, author FROM books WHERE id = 1",
				).Scan(&title, &author)
				if err != nil {
					t.Fatalf("book not found in database: %v", err)
				}
				if title != tt.title || author != tt.author {
					t.Fatalf("database mismatch: got (%s, %s), want (%s, %s)", title, author, tt.title, tt.author)
				}
			}
		})
	}
}

func TestBookAPI_GetBooks(t *testing.T) {
	app := setupTestApp(t)
	defer test.CleanupDatabase(t)

	t.Run("basic_list", func(t *testing.T) {
		createBook(t, "Animal Farm", "George Orwell")
		createBook(t, "1984", "George Orwell")

		resp, err := app.Client.GetBooks(context.Background(), &proto.GetBooksRequest{
			Page:    1,
			PerPage: 10,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resp.GetBooks()) != 2 {
			t.Fatalf("expected 2 books, got %d", len(resp.GetBooks()))
		}
	})

	t.Run("pagination", func(t *testing.T) {
		test.CleanupDatabase(t)
		for i := 1; i <= 5; i++ {
			createBook(t, "Book "+string(rune('0'+i)), "Author")
		}

		tests := []struct {
			name      string
			page      int32
			perPage   int32
			wantCount int
		}{
			{"page_1_size_2", 1, 2, 2},
			{"page_2_size_2", 2, 2, 2},
			{"page_3_size_2", 3, 2, 1},
			{"default_params", 0, 0, 5},
			{"large_per_page", 1, 100, 5},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				resp, err := app.Client.GetBooks(context.Background(), &proto.GetBooksRequest{
					Page:    tt.page,
					PerPage: tt.perPage,
				})
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if len(resp.GetBooks()) != tt.wantCount {
					t.Fatalf("expected %d books, got %d", tt.wantCount, len(resp.GetBooks()))
				}
			})
		}
	})
}

func TestBookAPI_GetBook(t *testing.T) {
	app := setupTestApp(t)
	defer test.CleanupDatabase(t)

	t.Run("success", func(t *testing.T) {
		createBook(t, "1984", "George Orwell")

		resp, err := app.Client.GetBook(context.Background(), &proto.GetBookRequest{Id: 1})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.GetBook().GetTitle() != "1984" {
			t.Fatalf("expected title '1984', got '%s'", resp.GetBook().GetTitle())
		}
	})

	t.Run("not_found", func(t *testing.T) {
		_, err := app.Client.GetBook(context.Background(), &proto.GetBookRequest{Id: 999})
		if status.Code(err) != codes.NotFound {
			t.Fatalf("expected NotFound, got %s", status.Code(err))
		}
	})
}

func TestBookAPI_UpdateBook(t *testing.T) {
	app := setupTestApp(t)
	defer test.CleanupDatabase(t)

	t.Run("success", func(t *testing.T) {
		createBook(t, "Original Title", "Original Author")

		_, err := app.Client.UpdateBook(context.Background(), &proto.UpdateBookRequest{
			Id:     1,
			Title:  "Updated Title",
			Author: "Updated Author",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var title, author string
		err = test.DB().QueryRow(
			context.Background(),
			"SELECT title, author FROM books WHERE id = 1",
		).Scan(&title, &author)
		if err != nil {
			t.Fatalf("book not found in database: %v", err)
		}
		if title != "Updated Title" || author != "Updated Author" {
			t.Fatalf("expected (Updated Title, Updated Author), got (%s, %s)", title, author)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		_, err := app.Client.UpdateBook(context.Background(), &proto.UpdateBookRequest{
			Id:     999,
			Title:  "Title",
			Author: "Author",
		})
		if status.Code(err) != codes.NotFound {
			t.Fatalf("expected NotFound, got %s", status.Code(err))
		}
	})
}

func TestBookAPI_DeleteBook(t *testing.T) {
	app := setupTestApp(t)
	defer test.CleanupDatabase(t)

	t.Run("success", func(t *testing.T) {
		createBook(t, "To Delete", "Author")

		_, err := app.Client.DeleteBook(context.Background(), &proto.DeleteBookRequest{Id: 1})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var count int
		err = test.DB().QueryRow(
			context.Background(),
			"SELECT COUNT(*) FROM books WHERE id = 1",
		).Scan(&count)
		if err != nil {
			t.Fatalf("failed to query database: %v", err)
		}
		if count != 0 {
			t.Fatal("book still exists in database after deletion")
		}
	})

	t.Run("not_found", func(t *testing.T) {
		_, err := app.Client.DeleteBook(context.Background(), &proto.DeleteBookRequest{Id: 999})
		if status.Code(err) != codes.NotFound {
			t.Fatalf("expected NotFound, got %s", status.Code(err))
		}
	})
}

func createBook(t *testing.T, title, author string) {
	t.Helper()
	_, err := test.DB().Exec(context.Background(),
		"INSERT INTO books (title, author) VALUES ($1, $2)",
		title, author,
	)
	if err != nil {
		t.Fatalf("failed to create book: %v", err)
	}
}
