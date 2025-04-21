ALTER TABLE candidate_questionnaires ADD COLUMN bucket_name VARCHAR(50) NULL AFTER expired_at;
ALTER TABLE candidate_questionnaires ADD COLUMN result_file_path VARCHAR(255) NULL AFTER bucket_name;