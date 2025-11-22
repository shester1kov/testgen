package moodle

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUploadQuizSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if contentType := r.Header.Get("Content-Type"); contentType == "" {
			t.Fatalf("expected multipart content type, got empty")
		}

		resp := UploadQuizResponse{QuizID: "1", CourseID: "2", Success: true, Message: "ok"}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	t.Cleanup(server.Close)

	client := NewClient(server.URL, "token")
	client.httpClient = server.Client()

	resp, err := client.UploadQuiz(context.Background(), UploadQuizRequest{CourseName: "c", QuizName: "q", XMLContent: "<xml/>"})
	if err != nil {
		t.Fatalf("expected upload to succeed, got %v", err)
	}

	if resp == nil || resp.QuizID != "1" || resp.CourseID != "2" || !resp.Success {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestUploadQuizErrors(t *testing.T) {
	t.Run("non-200 status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("bad request"))
		}))
		t.Cleanup(server.Close)

		client := NewClient(server.URL, "token")
		client.httpClient = server.Client()

		if _, err := client.UploadQuiz(context.Background(), UploadQuizRequest{XMLContent: "<xml/>"}); err == nil {
			t.Fatalf("expected error for non-OK status")
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("not json"))
		}))
		t.Cleanup(server.Close)

		client := NewClient(server.URL, "token")
		client.httpClient = server.Client()

		if _, err := client.UploadQuiz(context.Background(), UploadQuizRequest{XMLContent: "<xml/>"}); err == nil {
			t.Fatalf("expected decode error")
		}
	})
}

func TestGetCourses(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			courses := []Course{{ID: 1, ShortName: "short", FullName: "full", Category: 10}}
			_ = json.NewEncoder(w).Encode(courses)
		}))
		t.Cleanup(server.Close)

		client := NewClient(server.URL, "token")
		client.httpClient = server.Client()

		courses, err := client.GetCourses(context.Background())
		if err != nil {
			t.Fatalf("expected courses, got error: %v", err)
		}
		if len(courses) != 1 || courses[0].ID != 1 {
			t.Fatalf("unexpected courses: %+v", courses)
		}
	})

	t.Run("status error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		t.Cleanup(server.Close)

		client := NewClient(server.URL, "token")
		client.httpClient = server.Client()

		if _, err := client.GetCourses(context.Background()); err == nil {
			t.Fatalf("expected error for non-OK status")
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("invalid json"))
		}))
		t.Cleanup(server.Close)

		client := NewClient(server.URL, "token")
		client.httpClient = server.Client()

		if _, err := client.GetCourses(context.Background()); err == nil {
			t.Fatalf("expected decode error")
		}
	})
}

func TestValidateConnection(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{"status":"ok"}`))
		}))
		t.Cleanup(server.Close)

		client := NewClient(server.URL, "token")
		client.httpClient = server.Client()

		if err := client.ValidateConnection(context.Background()); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("moodle error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{"error":"token invalid"}`))
		}))
		t.Cleanup(server.Close)

		client := NewClient(server.URL, "token")
		client.httpClient = server.Client()

		if err := client.ValidateConnection(context.Background()); err == nil {
			t.Fatalf("expected error from moodle response")
		}
	})

	t.Run("decode error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("not json"))
		}))
		t.Cleanup(server.Close)

		client := NewClient(server.URL, "token")
		client.httpClient = server.Client()

		if err := client.ValidateConnection(context.Background()); err == nil {
			t.Fatalf("expected decode error")
		}
	})
}
