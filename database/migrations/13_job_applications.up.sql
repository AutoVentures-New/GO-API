CREATE TABLE IF NOT EXISTS job_applications (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    company_id BIGINT NOT NULL,
    job_id BIGINT NOT NULL,
    candidate_id BIGINT NOT NULL,
    steps JSON NOT NULL,
    current_step VARCHAR(50) NOT NULL,
    status ENUM('FILLING', 'WAITING_EVALUATION', 'REPROVED', 'APPROVED', 'CANCELED') NOT NULL,
    created_at DATETIME(6) NOT NULL,
    updated_at DATETIME(6) NOT NULL,

    FOREIGN KEY (job_id) REFERENCES jobs(id),
    FOREIGN KEY (company_id) REFERENCES companies(id),
    FOREIGN KEY (candidate_id) REFERENCES candidates(id)
);