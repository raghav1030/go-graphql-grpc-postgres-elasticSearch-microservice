CREATE TABLE IF NOT EXISTS orders(
    id CHAR(27) PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    account_id CHAR(27) NOT NULL,
    total_price MONEY NOT NULL
)
CREATE TABLE IF NOT EXISTS order_products (
    order_id CHAR(27) REFERENCES orders (id) ON DELETE CASCADE,
    product_id CHAR(27),
    quantity INT NOT NULL
    PRIMARY KEY (order_id, product_id)
)