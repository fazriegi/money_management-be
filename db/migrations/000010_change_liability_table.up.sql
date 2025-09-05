CREATE TABLE liability (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    date DATETIME NOT NULL,
    value VARCHAR(100) NOT NULL,
    user_id BIGINT NOT NULL,
    CONSTRAINT fk_liability_user FOREIGN KEY (user_id) REFERENCES user(id)
        ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX idx_liability_user_id ON liability(user_id);