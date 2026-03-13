package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/config"
	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/db"
)

type finding struct {
	Name   string
	Count  int64
	Sample string
}

func main() {
	if _, err := os.Stat(".env"); err == nil {
		_ = config.LoadDotEnv(".env")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("db error: %v", err)
	}

	checks := []struct {
		name      string
		countSQL  string
		sampleSQL string
	}{
		{
			name:      "tasks tenant mismatch with projects",
			countSQL:  "SELECT count(1) FROM tasks t JOIN projects p ON t.project_id = p.id WHERE t.tenant_id <> p.tenant_id AND t.deleted_at IS NULL AND p.deleted_at IS NULL",
			sampleSQL: "SELECT t.id::text FROM tasks t JOIN projects p ON t.project_id = p.id WHERE t.tenant_id <> p.tenant_id AND t.deleted_at IS NULL AND p.deleted_at IS NULL LIMIT 1",
		},
		{
			name:      "tasks assignee tenant mismatch",
			countSQL:  "SELECT count(1) FROM tasks t JOIN users u ON t.assignee_id = u.id WHERE t.assignee_id IS NOT NULL AND t.tenant_id <> u.tenant_id AND t.deleted_at IS NULL AND u.deleted_at IS NULL",
			sampleSQL: "SELECT t.id::text FROM tasks t JOIN users u ON t.assignee_id = u.id WHERE t.assignee_id IS NOT NULL AND t.tenant_id <> u.tenant_id AND t.deleted_at IS NULL AND u.deleted_at IS NULL LIMIT 1",
		},
		{
			name:      "task_comments tenant mismatch",
			countSQL:  "SELECT count(1) FROM task_comments c JOIN tasks t ON c.task_id = t.id WHERE c.tenant_id <> t.tenant_id AND c.deleted_at IS NULL AND t.deleted_at IS NULL",
			sampleSQL: "SELECT c.id::text FROM task_comments c JOIN tasks t ON c.task_id = t.id WHERE c.tenant_id <> t.tenant_id AND c.deleted_at IS NULL AND t.deleted_at IS NULL LIMIT 1",
		},
		{
			name:      "task_approvals tenant mismatch",
			countSQL:  "SELECT count(1) FROM task_approvals a JOIN tasks t ON a.task_id = t.id WHERE a.tenant_id <> t.tenant_id AND a.deleted_at IS NULL AND t.deleted_at IS NULL",
			sampleSQL: "SELECT a.id::text FROM task_approvals a JOIN tasks t ON a.task_id = t.id WHERE a.tenant_id <> t.tenant_id AND a.deleted_at IS NULL AND t.deleted_at IS NULL LIMIT 1",
		},
		{
			name:      "task_attachments tenant mismatch",
			countSQL:  "SELECT count(1) FROM task_attachments a JOIN tasks t ON a.task_id = t.id WHERE a.tenant_id <> t.tenant_id AND a.deleted_at IS NULL AND t.deleted_at IS NULL",
			sampleSQL: "SELECT a.id::text FROM task_attachments a JOIN tasks t ON a.task_id = t.id WHERE a.tenant_id <> t.tenant_id AND a.deleted_at IS NULL AND t.deleted_at IS NULL LIMIT 1",
		},
		{
			name:      "notifications tenant mismatch",
			countSQL:  "SELECT count(1) FROM notifications n JOIN users u ON n.user_id = u.id WHERE n.tenant_id <> u.tenant_id AND n.deleted_at IS NULL AND u.deleted_at IS NULL",
			sampleSQL: "SELECT n.id::text FROM notifications n JOIN users u ON n.user_id = u.id WHERE n.tenant_id <> u.tenant_id AND n.deleted_at IS NULL AND u.deleted_at IS NULL LIMIT 1",
		},
		{
			name:      "refresh_tokens tenant mismatch",
			countSQL:  "SELECT count(1) FROM refresh_tokens r JOIN users u ON r.user_id = u.id WHERE r.tenant_id <> u.tenant_id",
			sampleSQL: "SELECT r.id::text FROM refresh_tokens r JOIN users u ON r.user_id = u.id WHERE r.tenant_id <> u.tenant_id LIMIT 1",
		},
		{
			name:      "projects org tenant mismatch",
			countSQL:  "SELECT count(1) FROM projects p JOIN orgs o ON p.org_id = o.id WHERE p.org_id IS NOT NULL AND p.tenant_id <> o.tenant_id AND p.deleted_at IS NULL AND o.deleted_at IS NULL",
			sampleSQL: "SELECT p.id::text FROM projects p JOIN orgs o ON p.org_id = o.id WHERE p.org_id IS NOT NULL AND p.tenant_id <> o.tenant_id AND p.deleted_at IS NULL AND o.deleted_at IS NULL LIMIT 1",
		},
		{
			name:      "users org tenant mismatch",
			countSQL:  "SELECT count(1) FROM users u JOIN orgs o ON u.org_id = o.id WHERE u.org_id IS NOT NULL AND u.tenant_id <> o.tenant_id AND u.deleted_at IS NULL AND o.deleted_at IS NULL",
			sampleSQL: "SELECT u.id::text FROM users u JOIN orgs o ON u.org_id = o.id WHERE u.org_id IS NOT NULL AND u.tenant_id <> o.tenant_id AND u.deleted_at IS NULL AND o.deleted_at IS NULL LIMIT 1",
		},
	}

	results := make([]finding, 0, len(checks))
	for _, check := range checks {
		var count int64
		if err := database.Raw(check.countSQL).Scan(&count).Error; err != nil {
			results = append(results, finding{Name: check.name, Count: -1, Sample: fmt.Sprintf("error: %v", err)})
			continue
		}
		sample := ""
		if count > 0 {
			_ = database.Raw(check.sampleSQL).Scan(&sample).Error
		}
		results = append(results, finding{Name: check.name, Count: count, Sample: sample})
	}

	fmt.Printf("Permission scan at %s\n", time.Now().Format(time.RFC3339))
	for _, r := range results {
		status := "OK"
		if r.Count > 0 {
			status = "FAIL"
		}
		if r.Count < 0 {
			status = "ERROR"
		}
		if r.Sample != "" {
			fmt.Printf("- %s: %s (count=%d, sample=%s)\n", r.Name, status, r.Count, r.Sample)
		} else {
			fmt.Printf("- %s: %s (count=%d)\n", r.Name, status, r.Count)
		}
	}
}
