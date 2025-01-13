CREATE TABLE IF NOT EXISTS questionnaire_questions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    type ENUM('SINGLE_CHOICE', 'MULTIPLE_CHOICE', 'OPEN_FIELD') NOT NULL,
    questionnaire_id BIGINT NOT NULL,
    created_at DATETIME(6) NOT NULL,
    updated_at DATETIME(6) NOT NULL,

    FOREIGN KEY (questionnaire_id) REFERENCES questionnaires(id)
);

CREATE TABLE IF NOT EXISTS questionnaire_question_answers (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    questionnaire_id BIGINT NOT NULL,
    questionnaire_question_id BIGINT NOT NULL,
    is_correct TINYINT(1) NOT NULL,
    created_at DATETIME(6) NOT NULL,
    updated_at DATETIME(6) NOT NULL,

    FOREIGN KEY (questionnaire_id) REFERENCES questionnaires(id),
    FOREIGN KEY (questionnaire_question_id) REFERENCES questionnaire_questions(id)
);