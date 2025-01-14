CREATE TABLE IF NOT EXISTS jobs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    company_id BIGINT NOT NULL,
    is_talent_bank TINYINT(1) NOT NULL,
    is_special_needs TINYINT(1) NOT NULL,
    description TEXT NOT NULL,
    job_mode ENUM('IN_PERSON', 'HYBRID', 'SEMI_IN_PERSON', 'HOME_OFFICE') NOT NULL,
    contracting_modality ENUM('PJ', 'CLT') NOT NULL,
    state VARCHAR(255) NOT NULL,
    city VARCHAR(255) NOT NULL,
    responsibilities TEXT NOT NULL,
    questionnaire ENUM('PROFESSIONAL', 'BEHAVIORAL') NOT NULL,
    video_link VARCHAR(255) NOT NULL,
    status ENUM('ACTIVE', 'INACTIVE') NOT NULL,
    publish_at DATETIME(6) NOT NULL,
    finish_at DATETIME(6) NOT NULL,
    created_at DATETIME(6) NOT NULL,
    updated_at DATETIME(6) NOT NULL,

    FOREIGN KEY (company_id) REFERENCES companies(id)
);