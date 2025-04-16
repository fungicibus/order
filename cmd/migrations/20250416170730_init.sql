-- +goose Up
-- +goose StatementBegin
CREATE TYPE order_status AS ENUM (
   'pending',
   'rejected',
   'confirmed',
   'cancelled',
   'processing',
   'ready_for_shipping',
   'delivering',
   'delivered',
   'complete'
);

-- Create Orders table
CREATE TABLE
   orders (
      order_id CHAR(36) PRIMARY KEY,
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      "timestamp" TIMESTAMP NOT NULL,
      order_total DECIMAL(10, 2) NOT NULL,
      "status" order_status NOT NULL
   );

-- Create Order_Items table with price column
CREATE TABLE
   order_items (
      order_item_id BIGSERIAL PRIMARY KEY,
      order_id CHAR(36) NOT NULL REFERENCES orders (order_id) ON DELETE NO ACTION,
      product_id CHAR(36) NOT NULL,
      quantity INTEGER NOT NULL CHECK (quantity > 0),
      price DECIMAL(10, 2) NOT NULL CHECK (price >= 0)
   );

-- Create index for faster queries
CREATE INDEX idx_order_items_order_id ON order_items (order_id);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_order_items_order_id;

DROP TABLE IF EXISTS order_items;

DROP TABLE IF EXISTS orders;

DROP TYPE IF EXISTS order_status;

-- +goose StatementEnd