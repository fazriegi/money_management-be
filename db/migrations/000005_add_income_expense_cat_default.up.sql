CREATE TABLE expense_category_default (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);

CREATE TABLE income_category_default (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);

CREATE TABLE asset_category_default (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);

INSERT INTO expense_category_default (name) VALUES 
('Food'),
('Social Life'),
('Transport'),
('Health'),
('Entertainment'),
('Savings'),
('Education'),
('Internet'),
('Debt'),
('Gift'),
('Other');

INSERT INTO income_category_default (name) VALUES 
('Allowance'),
('Salary'),
('Bonus'),
('Dividend'),
('Other');

INSERT INTO asset_category_default (name) VALUES 
('Property'),
('Gold'),
('Cryptocurrency'),
('Stock'),
('Deposito'),
('Other');