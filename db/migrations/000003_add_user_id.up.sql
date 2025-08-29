ALTER TABLE incomes
ADD COLUMN user_id INT NOT NULL,
ADD CONSTRAINT fk_incomes_user
    FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE assets
ADD COLUMN user_id INT NOT NULL,
ADD CONSTRAINT fk_assets_user
    FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE expenses
ADD COLUMN user_id INT NOT NULL,
ADD CONSTRAINT fk_expenses_user
    FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE liabilities
ADD COLUMN user_id INT NOT NULL,
ADD CONSTRAINT fk_liabilities_user
    FOREIGN KEY (user_id) REFERENCES users(id);
