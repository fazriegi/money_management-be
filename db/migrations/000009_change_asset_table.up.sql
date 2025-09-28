CREATE TABLE asset_category (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    user_id BIGINT NOT NULL,
    CONSTRAINT fk_asset_cat_user FOREIGN KEY (user_id) REFERENCES user(id)
        ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX idx_asset_cat_user_id ON asset_category(user_id);

CREATE TABLE asset (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    category_id BIGINT NOT NULL,
    value VARCHAR(100) NOT NULL,
    notes VARCHAR(255) NOT NULL,
    amount VARCHAR(100) NOT NULL,
    user_id BIGINT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_asset_user FOREIGN KEY (user_id) REFERENCES user(id)
        ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_asset_category FOREIGN KEY (category_id) REFERENCES asset_category(id)
        ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE INDEX idx_asset_category_id ON asset(category_id);
CREATE INDEX idx_asset_user_id ON asset(user_id);