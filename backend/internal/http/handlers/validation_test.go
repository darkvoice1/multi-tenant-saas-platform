package handlers

import "testing"

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
