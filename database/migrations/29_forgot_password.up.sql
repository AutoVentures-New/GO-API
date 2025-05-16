CREATE TABLE IF NOT EXISTS forgot_password (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    token VARCHAR(255) NOT NULL,
    type ENUM('CANDIDATE', 'COMPANY') NOT NULL,
    expired_at DATETIME(6) NOT NULL,
    created_at DATETIME(6) NOT NULL,

    UNIQUE(user_id, type)
);