CREATE TABLE IF NOT EXISTS candidate_curriculum (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    candidate_id BIGINT NOT NULL,
    gender VARCHAR(255) NOT NULL,
    gender_identifier VARCHAR(255) NOT NULL,
    color VARCHAR(255) NOT NULL,
    is_special_needs TINYINT(1) NOT NULL,
    languages JSON NOT NULL,
    created_at DATETIME(6) NOT NULL,
    updated_at DATETIME(6) NOT NULL,

    FOREIGN KEY (candidate_id) REFERENCES candidates(id)
);

CREATE TABLE IF NOT EXISTS candidate_curriculum_professional_experience (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    candidate_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    company VARCHAR(255) NOT NULL,
    area_id BIGINT NOT NULL,
    city VARCHAR(255) NOT NULL,
    state VARCHAR(255) NOT NULL,
    job_mode ENUM('IN_PERSON', 'HYBRID', 'SEMI_IN_PERSON', 'HOME_OFFICE') NOT NULL,
    current_job TINYINT(1) NOT NULL,
    start_date DATETIME(6) NOT NULL,
    end_date DATETIME(6) NULL,
    created_at DATETIME(6) NOT NULL,
    updated_at DATETIME(6) NOT NULL,

    FOREIGN KEY (candidate_id) REFERENCES candidates(id),
    FOREIGN KEY (area_id) REFERENCES areas(id)
);

CREATE TABLE IF NOT EXISTS candidate_curriculum_academic_experience (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    candidate_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    company VARCHAR(255) NOT NULL,
    area_id BIGINT NOT NULL,
    status ENUM('IN_PROCESS', 'FINISHED') NOT NULL,
    level ENUM('INTERMEDIATE','HIGHER_EDUCATION','POSTGRADUATE','MASTER','PHD') NOT NULL,
    start_date DATETIME(6) NOT NULL,
    end_date DATETIME(6) NULL,
    created_at DATETIME(6) NOT NULL,
    updated_at DATETIME(6) NOT NULL,

    FOREIGN KEY (candidate_id) REFERENCES candidates(id),
    FOREIGN KEY (area_id) REFERENCES areas(id)
);