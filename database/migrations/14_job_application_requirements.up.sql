CREATE TABLE IF NOT EXISTS job_application_requirements (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    application_id BIGINT NOT NULL,
    items JSON NOT NULL,
    match_value INT NOT NULL,
    created_at DATETIME(6) NOT NULL,
    updated_at DATETIME(6) NOT NULL,

    FOREIGN KEY (application_id) REFERENCES job_applications(id)
);