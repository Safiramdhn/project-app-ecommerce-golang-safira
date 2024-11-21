-- Enum status_enum
CREATE TYPE status_enum AS ENUM ('active', 'deleted');

-- Tabel users
CREATE TABLE users (
  id VARCHAR PRIMARY KEY UNIQUE NOT NULL,
  name VARCHAR NOT NULL,
  email VARCHAR UNIQUE,
  phone_number VARCHAR UNIQUE,
  password VARCHAR NOT NULL,
  status status_enum DEFAULT 'active',
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);

CREATE SEQUENCE category_id_seq START 1;

-- Tabel categories
CREATE TABLE categories (
  id SERIAL PRIMARY KEY,
  name VARCHAR NOT NULL UNIQUE,
  status status_enum DEFAULT 'active',
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);

-- Tabel products
CREATE TABLE products (
  id SERIAL PRIMARY KEY,
  name VARCHAR NOT NULL,
  description VARCHAR NOT NULL,
  category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
  price DECIMAL(10, 2) NOT NULL,
  discount DECIMAL(10, 2) DEFAULT 0,
  rating DECIMAL(2, 1) CHECK (rating BETWEEN 1 AND 5) DEFAULT 0,
  photo_url TEXT,
  has_variant BOOLEAN DEFAULT FALSE,
  status status_enum DEFAULT 'active',
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);

CREATE TABLE variations (
    id SERIAL PRIMARY KEY,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    attribute_name VARCHAR NOT NULL, -- E.g., Size, Color
    status status_enum DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW()
);

-- Table Variation_Options (only for products with variations)
CREATE TABLE variation_options (
    id SERIAL PRIMARY KEY,
    variation_id INT REFERENCES variations(id) ON DELETE CASCADE,
    option_value VARCHAR NOT NULL, -- E.g., Red, Large
    additional_price DECIMAL(10, 2) DEFAULT 0,
    stock INT DEFAULT 0,
    status status_enum DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW()
);

SELECT * FROM users

INSERT INTO categories (name) VALUES 
('Electronics'), ('Books'), ('Clothing'),
('Home Appliances'), ('Toys'), ('Groceries'),
('Sports'), ('Health & Beauty'), ('Automotive'),
('Furniture'), ('Jewelry'), ('Music'),
('Stationery'), ('Gardening'), ('Pet Supplies');

SELECT * FROM categories

-- Insert Products (Sample: Clothing, Toys, Health & Beauty, Furniture)
INSERT INTO products (name, description, category_id, price, discount, rating, photo_url, has_variant, created_at)
VALUES
('Casual T-Shirt', 'A comfortable cotton T-shirt, perfect for everyday wear.', 3, 15.99, 2.00, 4.5, 'https://example.com/tshirt.jpg', TRUE, NOW() - INTERVAL '10 days'),
('Formal Shirt', 'A sleek and stylish formal shirt, ideal for office and events.', 3, 25.99, 5.00, 4.8, 'https://example.com/shirt.jpg', TRUE, NOW() - INTERVAL '20 days'),
('Stuffed Bear', 'A cute and cuddly stuffed bear, perfect as a gift or decoration.', 5, 12.99, 0.00, 4.2, 'https://example.com/bear.jpg', FALSE, NOW() - INTERVAL '5 days'),
('Skin Care Kit', 'A premium skin care kit to rejuvenate and nourish your skin.', 8, 45.99, 10.00, 4.7, 'https://example.com/skincare.jpg', TRUE, NOW() - INTERVAL '15 days'),
('Wooden Chair', 'A sturdy wooden chair, perfect for dining or working.', 10, 89.99, 0.00, 4.3, 'https://example.com/chair.jpg', FALSE, NOW() - INTERVAL '40 days');

SELECT * FROM products
SELECT id, name, description, price, discount, rating, photo_url, has_variant FROM products WHERE 1=1 AND category_id = 5

-- Variations for products with variations
INSERT INTO variations (product_id, attribute_name, created_at)
VALUES
(1, 'Size', NOW()),  -- For Casual T-Shirt
(1, 'Color', NOW()), -- For Casual T-Shirt
(2, 'Size', NOW()),  -- For Formal Shirt
(4, 'Fragrance', NOW()); -- For Skin Care Kit

SELECT * FROM variations

-- Variation Options for "Casual T-Shirt" (Product ID 1)
INSERT INTO variation_options (variation_id, option_value, additional_price, stock, created_at)
VALUES
(1, 'Small', 0.00, 50, NOW()),
(1, 'Medium', 1.00, 40, NOW()),
(1, 'Large', 2.00, 30, NOW()),
(2, 'Red', 0.00, 20, NOW()),
(2, 'Blue', 0.00, 25, NOW());

-- Variation Options for "Formal Shirt" (Product ID 2)
INSERT INTO variation_options (variation_id, option_value, additional_price, stock, created_at)
VALUES
(3, 'Small', 0.00, 15, NOW()),
(3, 'Medium', 1.50, 10, NOW()),
(3, 'Large', 2.50, 8, NOW());

-- Variation Options for "Skin Care Kit" (Product ID 4)
INSERT INTO variation_options (variation_id, option_value, additional_price, stock, created_at)
VALUES
(4, 'Lavender', 2.00, 25, NOW()),
(4, 'Rose', 2.50, 20, NOW());


SELECT * FROM products

-- -- Tabel recommendations
CREATE TABLE recommendations (
  id SERIAL PRIMARY KEY,
  product_id INT REFERENCES products(id) ON DELETE CASCADE,
  is_recommended BOOLEAN DEFAULT FALSE, -- New field for IsRecommended
  set_in_banner BOOLEAN DEFAULT FALSE,  -- New field for SetInBanner
  title VARCHAR NOT NULL,
  subtitle VARCHAR NOT NULL,
  photo_url TEXT,
  status status_enum DEFAULT 'active',
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);

INSERT INTO recommendations (
  product_id, is_recommended, set_in_banner, title, subtitle, photo_url, status, created_at
) VALUES
-- Recommendation for Casual T-Shirt
(1, TRUE, FALSE, 'Casual Comfort', 'Perfect for daily wear and casual outings.', 
 'https://example.com/tshirt.jpg', 'active', NOW()),

-- Recommendation for Formal Shirt
(2, TRUE, TRUE, 'Formal Elegance', 'Upgrade your wardrobe with sleek style.', 
 'https://example.com/shirt.jpg', 'active', NOW()),

-- Recommendation for Stuffed Bear
(3, TRUE, FALSE, 'Gift of Cuteness', 'An adorable gift for your loved ones.', 
 'https://example.com/bear.jpg', 'active', NOW()),

-- Recommendation for Skin Care Kit
(4, TRUE, TRUE, 'Skin Care Deluxe', 'Rejuvenate your skin with premium care.', 
 'https://example.com/skincare.jpg', 'active', NOW()),

-- Recommendation for Wooden Chair
(5, FALSE, FALSE, 'Elegant Wood Design', 'Perfect for your dining or working space.', 
 'https://example.com/chair.jpg', 'active', NOW());

-- Tabel wishlist
-- CREATE TABLE wishlist (
--   id SERIAL PRIMARY KEY,
--   user_id INT REFERENCES users(id) ON DELETE CASCADE,
--   product_id INT REFERENCES products(id) ON DELETE CASCADE,
--   status status_enum DEFAULT 'active',
--   created_at TIMESTAMP DEFAULT NOW(),
--   updated_at TIMESTAMP,
--   deleted_at TIMESTAMP
-- );

-- -- Tabel cart
-- CREATE TABLE cart (
--   id SERIAL PRIMARY KEY,
--   user_id INT REFERENCES users(id) ON DELETE CASCADE,
--   product_id INT REFERENCES products(id) ON DELETE CASCADE,
--   amount INT DEFAULT 1 CHECK (amount > 0),
--   status status_enum DEFAULT 'active',
--   created_at TIMESTAMP DEFAULT NOW(),
--   updated_at TIMESTAMP,
--   deleted_at TIMESTAMP
-- );

-- -- Tabel addresses
-- CREATE TABLE addresses (
--   id SERIAL PRIMARY KEY,
--   user_id INT REFERENCES users(id) ON DELETE CASCADE,
--   is_default BOOLEAN DEFAULT FALSE,
--   name VARCHAR NOT NULL,
--   street VARCHAR NOT NULL,
--   district VARCHAR,
--   country VARCHAR,
--   status status_enum DEFAULT 'active',
--   created_at TIMESTAMP DEFAULT NOW(),
--   updated_at TIMESTAMP,
--   deleted_at TIMESTAMP
-- );

-- -- Tabel shipping_types
-- CREATE TABLE shipping_types (
--   id SERIAL PRIMARY KEY,
--   name VARCHAR NOT NULL,
--   is_active BOOLEAN DEFAULT TRUE,
--   status status_enum DEFAULT 'active',
--   created_at TIMESTAMP DEFAULT NOW(),
--   updated_at TIMESTAMP,
--   deleted_at TIMESTAMP
-- );

-- -- Tabel payment_methods
-- CREATE TABLE payment_methods (
--   id SERIAL PRIMARY KEY,
--   name VARCHAR NOT NULL,
--   is_active BOOLEAN DEFAULT TRUE,
--   status status_enum DEFAULT 'active',
--   created_at TIMESTAMP DEFAULT NOW(),
--   updated_at TIMESTAMP,
--   deleted_at TIMESTAMP
-- );

-- -- Tabel order_items
-- CREATE TABLE order_items (
--   id SERIAL PRIMARY KEY,
--   product_id INT REFERENCES products(id) ON DELETE CASCADE,
--   amount INT CHECK (amount > 0),
--   price DECIMAL(10, 2),
--   order_id INT REFERENCES orders(id) ON DELETE CASCADE,
--   status status_enum DEFAULT 'active',
--   created_at TIMESTAMP DEFAULT NOW(),
--   updated_at TIMESTAMP,
--   deleted_at TIMESTAMP
-- );

-- -- Tabel orders
-- CREATE TABLE orders (
--   id SERIAL PRIMARY KEY,
--   user_id INT REFERENCES users(id) ON DELETE CASCADE,
--   total_price DECIMAL(10, 2),
--   address_id INT REFERENCES addresses(id),
--   shipping_type_id INT REFERENCES shipping_types(id),
--   payment_method_id INT REFERENCES payment_methods(id),
--   status status_enum DEFAULT 'active',
--   created_at TIMESTAMP DEFAULT NOW(),
--   updated_at TIMESTAMP,
--   deleted_at TIMESTAMP
-- );