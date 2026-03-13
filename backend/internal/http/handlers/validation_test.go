package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestValidatePassword(t *testing.T) {
	cases := []struct {
		name     string
		password string
		valid    bool
	}{
		{"valid", "Admin123", true},
		{"noUpper", "admin123", false},
		{"noLower", "ADMIN123", false},
		{"noDigit", "AdminPass", false},
		{"tooShort", "A1b", false},
	}

	for _, c := range cases {
		err := validatePassword(c.password)
		if c.valid && err != nil {
			t.Fatalf("%s expected valid, got %v", c.name, err)
		}
		if !c.valid && err == nil {
			t.Fatalf("%s expected error", c.name)
		}
	}
}

func TestRequireConfirm(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("confirm true", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("DELETE", "/?confirm=true", nil)
		if !requireConfirm(c) {
			t.Fatalf("expected confirm true")
		}
	})

	t.Run("confirm missing", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("DELETE", "/", nil)
		if requireConfirm(c) {
			t.Fatalf("expected confirm false")
		}
	})
}
