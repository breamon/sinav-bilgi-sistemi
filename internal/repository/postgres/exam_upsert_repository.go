package postgres

import (
	"github.com/breamon/sinav-bilgi-sistemi/internal/domain"
)

func (r *ExamRepository) UpsertImportedExam(exam *domain.Exam) (inserted bool, err error) {
	query := `
		INSERT INTO exams (
			source,
			external_id,
			title,
			description,
			category,
			status,
			application_start_date,
			application_end_date,
			exam_date,
			result_date
		)
		VALUES (
			:source,
			:external_id,
			:title,
			:description,
			:category,
			:status,
			:application_start_date,
			:application_end_date,
			:exam_date,
			:result_date
		)
		ON CONFLICT (source, external_id) WHERE external_id IS NOT NULL
		DO UPDATE SET
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			category = EXCLUDED.category,
			status = EXCLUDED.status,
			application_start_date = EXCLUDED.application_start_date,
			application_end_date = EXCLUDED.application_end_date,
			exam_date = EXCLUDED.exam_date,
			result_date = EXCLUDED.result_date,
			updated_at = NOW()
		RETURNING (xmax = 0) AS inserted;
	`

	rows, err := r.db.NamedQuery(query, exam)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&inserted); err != nil {
			return false, err
		}
	}

	return inserted, nil
}
