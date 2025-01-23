CREATE TABLE IF NOT EXISTS areas (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    created_at DATETIME(6) NOT NULL,
    updated_at DATETIME(6) NOT NULL
);

INSERT INTO areas(title, created_At, updated_at) VALUES ('Tecnologia', NOW(), NOW());
INSERT INTO areas(title, created_At, updated_at) VALUES ('Recursos Humanos', NOW(), NOW());
INSERT INTO areas(title, created_At, updated_at) VALUES ('Engenharia Civil', NOW(), NOW());
