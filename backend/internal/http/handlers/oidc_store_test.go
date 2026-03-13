package handlers

import "testing"

func TestOIDCStorePutPop(t *testing.T) {
	store := newOIDCStore()
	store.put("code1", oidcCode{Role: "admin"})
	if _, ok := store.pop("missing"); ok {
		t.Fatalf("expected missing code")
	}
	data, ok := store.pop("code1")
	if !ok || data.Role != "admin" {
		t.Fatalf("unexpected code data")
	}
	if _, ok := store.pop("code1"); ok {
		t.Fatalf("expected code removed")
	}
}

func TestGenerateCode(t *testing.T) {
	c1, err := generateCode()
	if err != nil || c1 == "" {
		t.Fatalf("generateCode error")
	}
	c2, _ := generateCode()
	if c1 == c2 {
		t.Fatalf("expected different codes")
	}
}
