CREATE TABLE IF NOT EXISTS candidate_photos (
    candidate_id BIGINT PRIMARY KEY,
    bucket_name VARCHAR(50) NOT NULL,
    photo_path VARCHAR(255) NOT NULL,
    created_at DATETIME(6) NOT NULL,
    updated_at DATETIME(6) NOT NULL,

    FOREIGN KEY (candidate_id) REFERENCES candidates(id)
);