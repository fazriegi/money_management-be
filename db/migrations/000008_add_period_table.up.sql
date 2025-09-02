CREATE TABLE monthly_period (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    day_of_month TINYINT NOT NULL,
    user_id BIGINT NOT NULL,
    CONSTRAINT fk_monthly_period_user FOREIGN KEY (user_id) REFERENCES user(id)
        ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX idx_monthly_period_user_id ON monthly_period(user_id);