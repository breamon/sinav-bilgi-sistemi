package postgres

import (
	"github.com/breamon/sinav-bilgi-sistemi/internal/domain"
	"github.com/jmoiron/sqlx"
)

type ExamRepository struct {
	db *sqlx.DB
}

func NewExamRepository(db *sqlx.DB) *ExamRepository {
	return &ExamRepository{db: db}
}

func (r *ExamRepository) Create(exam *domain.Exam) error {
	query := `
		INSERT INTO exams (source, title, status)
		VALUES ($1,$2,$3)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowx(
		query,
		exam.Source,
		exam.Title,
		exam.Status,
	).Scan(&exam.ID, &exam.CreatedAt, &exam.UpdatedAt)
}

func (r *ExamRepository) List() ([]domain.Exam, error) {
	var exams []domain.Exam

	query := `
		SELECT id, source, title, status, created_at, updated_at
		FROM exams
		ORDER BY id DESC
	`

	err := r.db.Select(&exams, query)
	return exams, err
}

func (r *ExamRepository) GetByID(id int64) (*domain.Exam, error) {
	var exam domain.Exam

	query := `
		SELECT id, source, title, status, created_at, updated_at
		FROM exams
		WHERE id = $1
	`

	err := r.db.Get(&exam, query, id)
	if err != nil {
		return nil, err
	}

	return &exam, nil
}

func (r *ExamRepository) Update(exam *domain.Exam) error {
	query := `
		UPDATE exams
		SET source = $1,
		    title = $2,
		    status = $3,
		    updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at
	`

	return r.db.QueryRowx(
		query,
		exam.Source,
		exam.Title,
		exam.Status,
		exam.ID,
	).Scan(&exam.UpdatedAt)
}

func (r *ExamRepository) Delete(id int64) error {
	query := `DELETE FROM exams WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
