package storage

import (
	"bytes"
	"context"
	"io"
	"testing"
)

func TestLocalStorageSaveOpen(t *testing.T) {
	dir := t.TempDir()
	store := NewLocal(dir)

	data := []byte("hello")
	if err := store.Save(context.Background(), "a/b.txt", bytes.NewReader(data), int64(len(data)), "text/plain"); err != nil {
		t.Fatalf("Save error: %v", err)
	}

	r, err := store.Open(context.Background(), "a/b.txt")
	if err != nil {
		t.Fatalf("Open error: %v", err)
	}
	defer r.Close()

	got, _ := io.ReadAll(r)
	if string(got) != "hello" {
		t.Fatalf("unexpected content: %s", string(got))
	}
}
