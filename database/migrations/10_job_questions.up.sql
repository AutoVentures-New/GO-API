CREATE TABLE IF NOT EXISTS job_questions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    company_id BIGINT NOT NULL,
    job_id BIGINT NOT NULL,
    question_id BIGINT NOT NULL,
    created_at DATETIME(6) NOT NULL,
    updated_at DATETIME(6) NOT NULL,

    FOREIGN KEY (job_id) REFERENCES jobs(id),
    FOREIGN KEY (company_id) REFERENCES companies(id),
    FOREIGN KEY (question_id) REFERENCES questionnaire_questions(id)
);