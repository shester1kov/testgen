package errors

import (
	goerrors "errors"
	"testing"
)

func TestAppErrorErrorFormatting(t *testing.T) {
	t.Run("without wrapped error", func(t *testing.T) {
		appErr := NewAppError(400, "bad request", nil)

		if got := appErr.Error(); got != "bad request" {
			t.Fatalf("expected message without suffix, got %q", got)
		}
	})

	t.Run("with wrapped error", func(t *testing.T) {
		inner := goerrors.New("db down")
		appErr := NewAppError(500, "internal", inner)

		if got := appErr.Error(); got != "internal: db down" {
			t.Fatalf("expected wrapped message, got %q", got)
		}
	})
}

func TestAppErrorConstructors(t *testing.T) {
	cases := []struct {
		name    string
		builder func() *AppError
		want    AppError
	}{
		{
			name:    "bad request",
			builder: func() *AppError { return BadRequest("oops") },
			want:    AppError{Code: 400, Message: "oops"},
		},
		{
			name:    "unauthorized",
			builder: func() *AppError { return Unauthorized("no token") },
			want:    AppError{Code: 401, Message: "no token"},
		},
		{
			name:    "forbidden",
			builder: func() *AppError { return Forbidden("denied") },
			want:    AppError{Code: 403, Message: "denied"},
		},
		{
			name:    "not found",
			builder: func() *AppError { return NotFound("missing") },
			want:    AppError{Code: 404, Message: "missing"},
		},
		{
			name:    "conflict",
			builder: func() *AppError { return Conflict("dup") },
			want:    AppError{Code: 409, Message: "dup"},
		},
		{
			name: "internal",
			builder: func() *AppError {
				return InternalServerError("boom", goerrors.New("root cause"))
			},
			want: AppError{Code: 500, Message: "boom"},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.builder()

			if err.Code != tt.want.Code || err.Message != tt.want.Message {
				t.Fatalf("unexpected error payload: %#v", err)
			}

			if tt.want.Err != nil && err.Err == nil {
				t.Fatalf("expected wrapped error to be present")
			}
		})
	}
}
