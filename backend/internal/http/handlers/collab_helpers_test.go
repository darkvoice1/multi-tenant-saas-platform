package handlers

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestNormalizeStatusAndPriority(t *testing.T) {
	if normalizeStatus("", "todo") != "todo" {
		t.Fatalf("expected fallback")
	}
	if normalizeStatus("In_Progress", "todo") != "in_progress" {
		t.Fatalf("expected normalized status")
	}
	if normalizeStatus("invalid", "todo") != "todo" {
		t.Fatalf("expected fallback on invalid status")
	}

	if normalizePriority("") != "medium" {
		t.Fatalf("expected default priority")
	}
	if normalizePriority("HIGH") != "high" {
		t.Fatalf("expected normalized priority")
	}
	if normalizePriority("bad") != "medium" {
		t.Fatalf("expected fallback priority")
	}
}

func TestBuildStorageKeyAndFileExt(t *testing.T) {
	tenantID := uuid.New()
	taskID := uuid.New()
	key := buildStorageKey(tenantID, taskID, "photo.PNG", "image/png")
	prefix := tenantID.String() + "/" + taskID.String() + "/"
	if !strings.HasPrefix(key, prefix) {
		t.Fatalf("unexpected key prefix: %s", key)
	}
	if !strings.HasSuffix(key, ".png") {
		t.Fatalf("expected .png suffix: %s", key)
	}

	if ext := fileExtFrom("", "image/jpeg"); ext == "" {
		t.Fatalf("expected ext from content type")
	}
	if ext := fileExtFrom("file.txt", ""); ext != ".txt" {
		t.Fatalf("expected .txt, got %s", ext)
	}
}

func TestSanitizeFileName(t *testing.T) {
	name := sanitizeFileName("..\\..\\evil.txt")
	if strings.Contains(name, "..") || strings.Contains(name, "\\") || name == "" {
		t.Fatalf("unexpected sanitized name: %s", name)
	}
}

func TestIsImageContentType(t *testing.T) {
	if !isImageContentType("image/png") {
		t.Fatalf("expected image content type")
	}
	if isImageContentType("text/plain") {
		t.Fatalf("expected non-image content type")
	}
}

func TestStatusTransition(t *testing.T) {
	if !isStatusTransitionAllowed("todo", "in_progress") {
		t.Fatalf("expected allowed transition")
	}
	if isStatusTransitionAllowed("done", "todo") {
		t.Fatalf("expected disallowed transition")
	}
	if !isStatusTransitionAllowed("review", "rejected") {
		t.Fatalf("expected allowed transition")
	}
}

func TestParseMentions(t *testing.T) {
	mentions := parseMentions("Hi @a@example.com and @b@example.com and @a@example.com")
	if len(mentions) != 2 {
		t.Fatalf("expected 2 mentions, got %d", len(mentions))
	}
}

func TestValidateLength(t *testing.T) {
	if !validateLength("abc", 1, 3) {
		t.Fatalf("expected valid length")
	}
	if validateLength("", 1, 3) {
		t.Fatalf("expected invalid length")
	}
}

func TestRequireConfirm(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/?confirm=true", nil)
	if !requireConfirm(c) {
		t.Fatalf("expected confirm true")
	}

	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request = httptest.NewRequest(http.MethodDelete, "/", nil)
	if requireConfirm(c2) {
		t.Fatalf("expected confirm false")
	}
	if w2.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w2.Code)
	}
}

func TestGenerateThumbnail(t *testing.T) {
	small := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			small.Set(x, y, color.RGBA{R: 200, G: 100, B: 50, A: 255})
		}
	}
	var smallBuf bytes.Buffer
	if err := png.Encode(&smallBuf, small); err != nil {
		t.Fatalf("png encode error: %v", err)
	}
	out, err := generateThumbnail(smallBuf.Bytes())
	if err != nil || len(out) == 0 {
		t.Fatalf("expected thumbnail for small image")
	}

	large := image.NewRGBA(image.Rect(0, 0, 512, 512))
	var largeBuf bytes.Buffer
	if err := png.Encode(&largeBuf, large); err != nil {
		t.Fatalf("png encode error: %v", err)
	}
	out2, err := generateThumbnail(largeBuf.Bytes())
	if err != nil || len(out2) == 0 {
		t.Fatalf("expected thumbnail for large image")
	}
}
