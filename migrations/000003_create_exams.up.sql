CREATE TABLE IF NOT EXISTS exams (
    id BIGSERIAL PRIMARY KEY,
    source VARCHAR(100) NOT NULL,
    external_id VARCHAR(255),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100),
    application_start_date TIMESTAMP,
    application_end_date TIMESTAMP,
    exam_date TIMESTAMP,
    result_date TIMESTAMP,
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);