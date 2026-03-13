package middleware

import "testing"

func TestDeriveAuditAction(t *testing.T) {
	cases := []struct {
		method   string
		path     string
		expected string
		resource string
	}{
		{"POST", "/api/projects", "create", "project"},
		{"GET", "/api/tasks/123", "read", "task"},
		{"POST", "/api/tasks/123/approve", "approve", "task"},
		{"POST", "/api/tasks/123/status", "status_change", "task"},
		{"POST", "/api/notifications/1/read", "mark_read", "notification"},
		{"DELETE", "/api/admin/users/1", "delete", "admin"},
	}

	for _, c := range cases {
		action, resource := deriveAuditAction(c.method, c.path)
		if action != c.expected || resource != c.resource {
			t.Fatalf("%s %s => %s/%s, got %s/%s", c.method, c.path, c.expected, c.resource, action, resource)
		}
	}
}
