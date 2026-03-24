package domain

import "time"

type ImportLog struct {
	ID            int64     `db:"id" json:"id"`
	Provider      string    `db:"provider" json:"provider"`
	Status        string    `db:"status" json:"status"`
	ImportedCount int       `db:"imported_count" json:"imported_count"`
	ErrorMessage  *string   `db:"error_message" json:"error_message,omitempty"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}
