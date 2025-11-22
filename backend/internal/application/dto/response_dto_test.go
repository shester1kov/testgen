package dto

import "testing"

func TestResponseHelpers(t *testing.T) {
	errResp := NewErrorResponse(ErrCodeInvalidInput, "bad payload")
	if errResp.Error.Code != ErrCodeInvalidInput || errResp.Error.Message != "bad payload" {
		t.Fatalf("unexpected error response: %+v", errResp)
	}

	success := NewSuccessResponse("ok")
	if !success.Success || success.Message != "ok" {
		t.Fatalf("unexpected success response: %+v", success)
	}

	message := NewMessageResponse("hello")
	if message.Message != "hello" {
		t.Fatalf("unexpected message response: %+v", message)
	}
}
