DROP TABLE income;
DROP TABLE user_income_category;
DROP TABLE user;

CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255)
);

CREATE INDEX idx_users_username ON users(username);

CREATE TABLE IF NOT EXISTS incomes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    period_code VARCHAR(50) NOT NULL,
    type VARCHAR(10) NOT NULL,
    name VARCHAR(255) NOT NULL,
    value VARCHAR(255) NOT NULL,
    order_no INT NOT NULL,q
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE INDEX idx_incomes_period_code ON incomes(period_code);

CREATE TABLE IF NOT EXISTS assets (
    id INT AUTO_INCREMENT PRIMARY KEY,
    period_code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    amount VARCHAR(255) NOT NULL,
    value VARCHAR(255) NOT NULL,
    order_no INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE INDEX idx_assets_period_code ON assets(period_code);

CREATE TABLE IF NOT EXISTS liabilities (
    id INT AUTO_INCREMENT PRIMARY KEY,
    period_code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    value VARCHAR(255) NOT NULL,
    order_no INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE INDEX idx_liabilities_period_code ON liabilities(period_code);

CREATE TABLE IF NOT EXISTS expenses (
    id INT AUTO_INCREMENT PRIMARY KEY,
    period_code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    value VARCHAR(255) NOT NULL,
    order_no INT NOT NULL,
    liability_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_expenses_liability
      FOREIGN KEY (liability_id) REFERENCES liabilities(id) ON DELETE SET NULL
);

CREATE INDEX idx_expenses_period_code ON expenses(period_code);
CREATE INDEX idx_expenses_liability_id ON expenses(liability_id);
