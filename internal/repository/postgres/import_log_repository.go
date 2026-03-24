package postgres

import (
	"github.com/breamon/sinav-bilgi-sistemi/internal/domain"
	"github.com/jmoiron/sqlx"
)

type ImportLogRepository struct {
	db *sqlx.DB
}

func NewImportLogRepository(db *sqlx.DB) *ImportLogRepository {
	return &ImportLogRepository{db: db}
}

func (r *ImportLogRepository) Create(logEntry *domain.ImportLog) error {
	query := `
		INSERT INTO import_logs (provider, status, imported_count, error_message)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	return r.db.QueryRowx(
		query,
		logEntry.Provider,
		logEntry.Status,
		logEntry.ImportedCount,
		logEntry.ErrorMessage,
	).Scan(&logEntry.ID, &logEntry.CreatedAt)
}

func (r *ImportLogRepository) List(limit int) ([]domain.ImportLog, error) {
	var logs []domain.ImportLog

	query := `
		SELECT id, provider, status, imported_count, error_message, created_at
		FROM import_logs
		ORDER BY id DESC
		LIMIT $1
	`

	err := r.db.Select(&logs, query, limit)
	return logs, err
}
