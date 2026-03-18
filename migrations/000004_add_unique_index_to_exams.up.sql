CREATE UNIQUE INDEX IF NOT EXISTS ux_exams_source_external_id
ON exams (source, external_id)
WHERE external_id IS NOT NULL;