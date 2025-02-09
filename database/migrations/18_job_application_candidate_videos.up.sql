CREATE TABLE IF NOT EXISTS job_application_candidate_videos (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    application_id BIGINT NOT NULL,
    bucket_name VARCHAR(50) NOT NULL,
    video_path VARCHAR(255) NOT NULL,
    created_at DATETIME(6) NOT NULL,
    updated_at DATETIME(6) NOT NULL,

    FOREIGN KEY (application_id) REFERENCES job_applications(id)
);