DROP TABLE incomes;
DROP TABLE assets;
DROP TABLE expenses;
DROP TABLE liabilities;
DROP TABLE users;

CREATE TABLE IF NOT EXISTS user (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255)
);

CREATE INDEX idx_user_username ON user(username);

CREATE TABLE user_income_category (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    user_id BIGINT NOT NULL,
    CONSTRAINT fk_income_cat_user FOREIGN KEY (user_id) REFERENCES user(id)
        ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX idx_income_cat_user_id ON user_income_category(user_id);

CREATE TABLE income (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    category_id BIGINT NOT NULL,
    date DATETIME NOT NULL,
    value VARCHAR(100) NOT NULL,
    user_id BIGINT NOT NULL,
    notes VARCHAR(255),
    CONSTRAINT fk_income_user FOREIGN KEY (user_id) REFERENCES user(id)
        ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_income_category FOREIGN KEY (category_id) REFERENCES user_income_category(id)
        ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE INDEX idx_income_category_id ON income(category_id);
CREATE INDEX idx_income_user_id ON income(user_id);