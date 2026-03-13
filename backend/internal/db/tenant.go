package db

import "gorm.io/gorm"

// WithTenant scopes queries to a tenant_id column.
func WithTenant(db *gorm.DB, tenantID string) *gorm.DB {
	if tenantID == "" {
		return db
	}
	return db.Where("tenant_id = ?", tenantID)
}
