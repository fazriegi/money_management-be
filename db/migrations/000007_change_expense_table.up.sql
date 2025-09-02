CREATE TABLE user_expense_category (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    user_id BIGINT NOT NULL,
    CONSTRAINT fk_expense_cat_user FOREIGN KEY (user_id) REFERENCES user(id)
        ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX idx_expense_cat_user_id ON user_expense_category(user_id);

CREATE TABLE expense (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    category_id BIGINT NOT NULL,
    date DATETIME NOT NULL,
    value VARCHAR(100) NOT NULL,
    user_id BIGINT NOT NULL,
    notes VARCHAR(255),
    CONSTRAINT fk_expense_user FOREIGN KEY (user_id) REFERENCES user(id)
        ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_expense_category FOREIGN KEY (category_id) REFERENCES user_expense_category(id)
        ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE INDEX idx_expense_category_id ON expense(category_id);
CREATE INDEX idx_expense_user_id ON expense(user_id);