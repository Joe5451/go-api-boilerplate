package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"go-api-boilerplate/test"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	test.Setup()
	code := m.Run()
	test.Teardown()
	os.Exit(code)
}

func setupTestApp(t *testing.T) *test.RestApp {
	app, err := test.InitRestApp(context.Background())
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
		name       string
		body       map[string]string
		wantStatus int
		checkDB    bool
	}{
		{
			name:       "success",
			body:       map[string]string{"title": "1984", "author": "George"},
			wantStatus: http.StatusNoContent,
			checkDB:    true,
		},
		{
			name:       "missing_title",
			body:       map[string]string{"author": "Author"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing_author",
			body:       map[string]string{"title": "Title"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty_title",
			body:       map[string]string{"title": "", "author": "Author"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty_author",
			body:       map[string]string{"title": "Title", "author": ""},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.body)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(
				"POST",
				"/books",
				bytes.NewBuffer(jsonBody),
			)
			req.Header.Set("Content-Type", "application/json")

			app.Router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf(
					"expected %d, got %d: %s",
					tt.wantStatus,
					w.Code,
					w.Body.String(),
				)
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

				wantTitle := tt.body["title"]
				wantAuthor := tt.body["author"]
				if title != wantTitle || author != wantAuthor {
					t.Errorf(
						"database mismatch: got (%s, %s), want (%s, %s)",
						title,
						author,
						wantTitle,
						wantAuthor,
					)
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

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/books?page=1&per_page=10", nil)
		app.Router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf(
				"expected 200, got %d: %s",
				w.Code,
				w.Body.String(),
			)
		}

		var books []map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &books); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if len(books) != 2 {
			t.Errorf("expected 2 books, got %d", len(books))
		}

		var count int
		test.DB().QueryRow(
			context.Background(),
			"SELECT COUNT(*) FROM books",
		).Scan(&count)

		if count != 2 {
			t.Errorf("expected 2 books in database, got %d", count)
		}
	})

	t.Run("pagination", func(t *testing.T) {
		test.CleanupDatabase(t)
		for i := 1; i <= 5; i++ {
			createBook(t, "Book "+string(rune('0'+i)), "Author")
		}

		tests := []struct {
			name      string
			query     string
			wantCount int
		}{
			{"page_1_size_2", "?page=1&per_page=2", 2},
			{"page_2_size_2", "?page=2&per_page=2", 2},
			{"page_3_size_2", "?page=3&per_page=2", 1},
			{"default_params", "", 5},
			{"large_per_page", "?page=1&per_page=100", 5},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/books"+tt.query, nil)
				app.Router.ServeHTTP(w, req)

				if w.Code != http.StatusOK {
					t.Errorf("expected 200, got %d", w.Code)
				}

				var books []map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &books)

				if len(books) != tt.wantCount {
					t.Errorf(
						"expected %d books, got %d",
						tt.wantCount,
						len(books),
					)
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

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/books/1", nil)
		app.Router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
		}

		var book map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &book); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if book["title"] != "1984" {
			t.Errorf("expected title '1984', got '%v'", book["title"])
		}
	})

	t.Run("not_found", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/books/999", nil)
		app.Router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})
}

func TestBookAPI_UpdateBook(t *testing.T) {
	app := setupTestApp(t)
	defer test.CleanupDatabase(t)

	t.Run("success", func(t *testing.T) {
		createBook(t, "Original Title", "Original Author")

		body := map[string]string{
			"title":  "Updated Title",
			"author": "Updated Author",
		}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"PUT",
			"/books/1",
			bytes.NewBuffer(jsonBody),
		)
		req.Header.Set("Content-Type", "application/json")

		app.Router.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d: %s", w.Code, w.Body.String())
		}

		var title, author string
		err := test.DB().QueryRow(
			context.Background(),
			"SELECT title, author FROM books WHERE id = 1",
		).Scan(&title, &author)

		if err != nil {
			t.Fatalf("book not found in database: %v", err)
		}

		if title != "Updated Title" || author != "Updated Author" {
			t.Errorf(
				"expected (Updated Title, Updated Author), got (%s, %s)",
				title,
				author,
			)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		body := map[string]string{
			"title":  "Title",
			"author": "Author",
		}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/books/999", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		app.Router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})
}

func TestBookAPI_DeleteBook(t *testing.T) {
	app := setupTestApp(t)
	defer test.CleanupDatabase(t)

	t.Run("success", func(t *testing.T) {
		createBook(t, "To Delete", "Author")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/books/1", nil)
		app.Router.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d: %s", w.Code, w.Body.String())
		}

		var count int
		err := test.DB().QueryRow(
			context.Background(),
			"SELECT COUNT(*) FROM books WHERE id = 1",
		).Scan(&count)

		if err != nil {
			t.Fatalf("failed to query database: %v", err)
		}

		if count != 0 {
			t.Errorf("book still exists in database after deletion")
		}
	})

	t.Run("not_found", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/books/999", nil)
		app.Router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})
}

// Helper function to create a book
func createBook(t *testing.T, title, author string) {
	_, err := test.DB().Exec(context.Background(),
		"INSERT INTO books (title, author) VALUES ($1, $2)",
		title, author,
	)
	if err != nil {
		t.Fatalf("failed to create book: %v", err)
	}
}
