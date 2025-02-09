CREATE TABLE IF NOT EXISTS candidate_questionnaires (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    candidate_id BIGINT NOT NULL,
    type ENUM('PROFESSIONAL', 'BEHAVIORAL') NOT NULL,
    answers JSON NOT NULL,
    expired_at DATETIME(6) NOT NULL,
    created_at DATETIME(6) NOT NULL,
    updated_at DATETIME(6) NOT NULL,

    FOREIGN KEY (candidate_id) REFERENCES candidates(id)
);