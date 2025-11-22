package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Positive test - successful token generation and validation
func TestJWTManager_GenerateAndValidate_Success(t *testing.T) {
	manager, err := NewJWTManager("test-secret-key", "1h")
	if err != nil {
		t.Fatalf("Failed to create JWT manager: %v", err)
	}

	userID := uuid.New()
	email := "test@example.com"
	role := "teacher"

	token, err := manager.GenerateToken(userID, email, role)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	claims, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected userID %s, got %s", userID, claims.UserID)
	}
	if claims.Email != email {
		t.Errorf("Expected email %s, got %s", email, claims.Email)
	}
	if claims.Role != role {
		t.Errorf("Expected role %s, got %s", role, claims.Role)
	}
}

// NEGATIVE TEST: Invalid expiration format
func TestJWTManager_NewJWTManager_InvalidExpiration(t *testing.T) {
	testCases := []struct {
		name       string
		expiration string
	}{
		{"empty string", ""},
		{"invalid format", "invalid"},
		{"no unit", "123"},
		{"invalid unit", "1x"},
		// Note: negative durations are valid in Go, so this won't error
		// {"negative", "-1h"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewJWTManager("secret", tc.expiration)
			if err == nil {
				t.Errorf("Expected error for expiration '%s', got nil", tc.expiration)
			}
		})
	}
}

// NEGATIVE TEST: Validate malformed token
func TestJWTManager_ValidateToken_MalformedToken(t *testing.T) {
	manager, _ := NewJWTManager("test-secret", "1h")

	testCases := []struct {
		name  string
		token string
	}{
		{"empty token", ""},
		{"invalid format", "not.a.token"},
		{"random string", "randomstring"},
		{"only two parts", "header.payload"},
		{"four parts", "a.b.c.d"},
		{"invalid base64", "!!!.!!!.!!!"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := manager.ValidateToken(tc.token)
			if err == nil {
				t.Errorf("Expected error for token '%s', got nil", tc.token)
			}
		})
	}
}

// NEGATIVE TEST: Validate token with wrong secret
func TestJWTManager_ValidateToken_WrongSecret(t *testing.T) {
	manager1, _ := NewJWTManager("secret1", "1h")
	manager2, _ := NewJWTManager("secret2", "1h")

	userID := uuid.New()
	token, err := manager1.GenerateToken(userID, "test@example.com", "teacher")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Try to validate with different secret - should fail
	_, err = manager2.ValidateToken(token)
	if err == nil {
		t.Error("Expected error when validating token with wrong secret, got nil")
	}
}

// NEGATIVE TEST: Expired token
func TestJWTManager_ValidateToken_ExpiredToken(t *testing.T) {
	// Create manager with very short expiration
	manager, _ := NewJWTManager("test-secret", "1ns") // 1 nanosecond

	userID := uuid.New()
	token, err := manager.GenerateToken(userID, "test@example.com", "teacher")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Wait a bit to ensure token expires
	time.Sleep(10 * time.Millisecond)

	// Try to validate expired token - should fail
	_, err = manager.ValidateToken(token)
	if err == nil {
		t.Error("Expected error for expired token, got nil")
	}
}

// NEGATIVE TEST: Token with wrong signing method
func TestJWTManager_ValidateToken_WrongSigningMethod(t *testing.T) {
	manager, _ := NewJWTManager("test-secret", "1h")

	// Create token with RS256 instead of HS256
	claims := &JWTClaims{
		UserID: uuid.New(),
		Email:  "test@example.com",
		Role:   "teacher",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Use None signing method (insecure)
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenString, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)

	// Try to validate - should fail due to wrong signing method
	_, err := manager.ValidateToken(tokenString)
	if err == nil {
		t.Error("Expected error for token with wrong signing method, got nil")
	}
}

// NEGATIVE TEST: Token with tampered payload
func TestJWTManager_ValidateToken_TamperedPayload(t *testing.T) {
	manager, _ := NewJWTManager("test-secret", "1h")

	userID := uuid.New()
	token, _ := manager.GenerateToken(userID, "test@example.com", "student")

	// Tamper with the token by modifying the middle part (payload)
	// This will cause signature verification to fail
	tamperedToken := token[:len(token)/2] + "TAMPERED" + token[len(token)/2:]

	_, err := manager.ValidateToken(tamperedToken)
	if err == nil {
		t.Error("Expected error for tampered token, got nil")
	}
}

// NEGATIVE TEST: Token without required claims
// Note: JWT library will parse it and fill zero values for missing fields
// This is expected behavior - validation of required fields should be done at business logic level
func TestJWTManager_ValidateToken_MissingClaims(t *testing.T) {
	manager, _ := NewJWTManager("test-secret", "1h")

	// Create token with standard claims only (no custom claims)
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("test-secret"))

	// JWT library will accept this and fill zero values
	parsedClaims, err := manager.ValidateToken(tokenString)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify that custom fields have zero values
	if parsedClaims.UserID != uuid.Nil {
		t.Error("Expected zero UUID for missing UserID claim")
	}
	if parsedClaims.Email != "" {
		t.Error("Expected empty string for missing Email claim")
	}
	if parsedClaims.Role != "" {
		t.Error("Expected empty string for missing Role claim")
	}
}

// NEGATIVE TEST: Generate token with empty/invalid values
func TestJWTManager_GenerateToken_InvalidInputs(t *testing.T) {
	manager, _ := NewJWTManager("test-secret", "1h")

	testCases := []struct {
		name   string
		userID uuid.UUID
		email  string
		role   string
	}{
		{"empty email", uuid.New(), "", "teacher"},
		{"empty role", uuid.New(), "test@example.com", ""},
		{"zero UUID", uuid.UUID{}, "test@example.com", "teacher"},
		{"all empty", uuid.UUID{}, "", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Token generation should succeed even with empty values
			// (validation of business rules should happen at handler level)
			token, err := manager.GenerateToken(tc.userID, tc.email, tc.role)
			if err != nil {
				t.Errorf("Token generation failed unexpectedly: %v", err)
			}

			// But we can verify the claims are what we set
			claims, err := manager.ValidateToken(token)
			if err != nil {
				t.Errorf("Token validation failed: %v", err)
			}

			if claims.Email != tc.email {
				t.Errorf("Expected email '%s', got '%s'", tc.email, claims.Email)
			}
		})
	}
}

// NEGATIVE TEST: Token issued in the future (NotBefore)
func TestJWTManager_ValidateToken_NotYetValid(t *testing.T) {
	manager, _ := NewJWTManager("test-secret", "1h")

	// Create token that's not valid yet
	claims := &JWTClaims{
		UserID: uuid.New(),
		Email:  "test@example.com",
		Role:   "teacher",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)), // Valid only after 1 hour
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("test-secret"))

	// Should fail because NotBefore is in the future
	_, err := manager.ValidateToken(tokenString)
	if err == nil {
		t.Error("Expected error for token not yet valid (NotBefore), got nil")
	}
}
