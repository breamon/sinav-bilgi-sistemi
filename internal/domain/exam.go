package domain

import "time"

type Exam struct {
	ID                   int64      `db:"id" json:"id"`
	Source               string     `db:"source" json:"source"`
	ExternalID           *string    `db:"external_id" json:"external_id,omitempty"`
	Title                string     `db:"title" json:"title"`
	Description          *string    `db:"description" json:"description,omitempty"`
	Category             *string    `db:"category" json:"category,omitempty"`
	ApplicationStartDate *time.Time `db:"application_start_date" json:"application_start_date,omitempty"`
	ApplicationEndDate   *time.Time `db:"application_end_date" json:"application_end_date,omitempty"`
	ExamDate             *time.Time `db:"exam_date" json:"exam_date,omitempty"`
	ResultDate           *time.Time `db:"result_date" json:"result_date,omitempty"`
	Status               string     `db:"status" json:"status"`
	CreatedAt            time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time  `db:"updated_at" json:"updated_at"`
}
