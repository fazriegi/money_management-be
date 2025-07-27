CREATE TABLE IF NOT EXISTS incomes (
    id SERIAL PRIMARY KEY,
    period_code VARCHAR(50) NOT NULL,
    type VARCHAR(10) NOT NULL,
    name VARCHAR(255) NOT NULL,
    value VARCHAR(255) NOT NULL,
    order_no INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_incomes_period_code ON incomes(period_code);

CREATE TABLE IF NOT EXISTS assets (
    id SERIAL PRIMARY KEY,
    period_code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    amount VARCHAR(255) NOT NULL,
    value VARCHAR(255) NOT NULL,
    order_no INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_assets_period_code ON assets(period_code);

CREATE TABLE IF NOT EXISTS liabilities (
    id SERIAL PRIMARY KEY,
    period_code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    value VARCHAR(255) NOT NULL,
    order_no INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_liabilities_period_code ON liabilities(period_code);

CREATE TABLE IF NOT EXISTS expenses (
    id SERIAL PRIMARY KEY,
    period_code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    value VARCHAR(255) NOT NULL,
    order_no INT NOT NULL,
    liability_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (liability_id) REFERENCES liabilities(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_expenses_period_code ON expenses(period_code);
CREATE INDEX IF NOT EXISTS idx_expenses_liability_id ON expenses(liability_id);