CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    receipt_name VARCHAR(255),
    pretty_name VARCHAR(255),
    product_id VARCHAR(255),
    description TEXT
)