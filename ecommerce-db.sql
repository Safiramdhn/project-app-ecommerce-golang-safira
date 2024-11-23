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
  total_stock INT DEFAULT 0,
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
INSERT INTO products (name, description, category_id, price, discount, rating, photo_url, has_variant, total_stock, created_at)
VALUES
('Casual T-Shirt', 'A comfortable cotton T-shirt, perfect for everyday wear.', 3, 15.99, 2.00, 4.5, 'https://example.com/tshirt.jpg', TRUE, 120, NOW() - INTERVAL '10 days'),
('Formal Shirt', 'A sleek and stylish formal shirt, ideal for office and events.', 3, 25.99, 5.00, 4.8, 'https://example.com/shirt.jpg', TRUE, 33, NOW() - INTERVAL '20 days'),
('Stuffed Bear', 'A cute and cuddly stuffed bear, perfect as a gift or decoration.', 5, 12.99, 0.00, 4.2, 'https://example.com/bear.jpg', FALSE, 100, NOW() - INTERVAL '5 days'),
('Skin Care Kit', 'A premium skin care kit to rejuvenate and nourish your skin.', 8, 45.99, 10.00, 4.7, 'https://example.com/skincare.jpg', TRUE, 45, NOW() - INTERVAL '15 days'),
('Wooden Chair', 'A sturdy wooden chair, perfect for dining or working.', 10, 89.99, 0.00, 4.3, 'https://example.com/chair.jpg', FALSE, 100, NOW() - INTERVAL '40 days');

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
(2, 'Red', 0.00, 60, NOW()),
(2, 'Blue', 0.00, 60, NOW());

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
CREATE TABLE wishlist (
  id SERIAL PRIMARY KEY,
  user_id VARCHAR REFERENCES users(id) ON DELETE CASCADE,
  product_id INT REFERENCES products(id) ON DELETE CASCADE,
  status status_enum DEFAULT 'active',
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);

SELECT * FROM wishlist

CREATE TABLE weekly_promos (
	id SERIAL PRIMARY KEY,
	product_id INT REFERENCES products(id) ON DELETE CASCADE,
	promo_discount  DECIMAL(10, 2) DEFAULT 0 NOT NULL,
	start_date DATE,
	end_date DATE,
	status status_enum DEFAULT 'active',
  	created_at TIMESTAMP DEFAULT NOW(),
  	updated_at TIMESTAMP,
  	deleted_at TIMESTAMP
)

INSERT INTO weekly_promos (product_id, promo_discount, start_date, end_date)
VALUES
-- Promo for Casual T-Shirt
(1, 3.00, CURRENT_DATE - INTERVAL '5 days', CURRENT_DATE + INTERVAL '2 days'),

-- Promo for Formal Shirt
(2, 7.00, CURRENT_DATE - INTERVAL '10 days', CURRENT_DATE - INTERVAL '3 days'),

-- Promo for Stuffed Bear
(3, 1.50, CURRENT_DATE - INTERVAL '2 days', CURRENT_DATE + INTERVAL '5 days'),

-- Promo for Skin Care Kit
(4, 12.00, CURRENT_DATE - INTERVAL '7 days', CURRENT_DATE + INTERVAL '7 days'),

-- Promo for Wooden Chair
(5, 10.00, CURRENT_DATE - INTERVAL '15 days', CURRENT_DATE - INTERVAL '8 days');

-- Tabel addresses
CREATE TABLE addresses (
  id SERIAL PRIMARY KEY,
  user_id VARCHAR NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  is_default BOOLEAN DEFAULT FALSE NOT NULL,
  name VARCHAR(100) NOT NULL, -- Specify a reasonable length for names
  street VARCHAR(255) NOT NULL, -- Allow for longer street addresses
  district VARCHAR(100), -- Add a length constraint for district
  city VARCHAR(100), -- Add city for better address structuring
  state VARCHAR(100), -- Add state for more detailed geographic data
  postal_code VARCHAR(20), -- Add postal code to improve address accuracy
  country VARCHAR(100) NOT NULL, -- Make country mandatory
  status status_enum DEFAULT 'active' NOT NULL,
  created_at TIMESTAMP DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP -- Nullable for soft deletes
);

SELECT * FROM addresses


-- -- Tabel cart
CREATE TABLE carts (
    id SERIAL PRIMARY KEY,                      -- Auto-incrementing ID for the cart
    user_id VARCHAR REFERENCES users(id) DELETE ON CASCADE,              -- User ID (assuming it's a string)
    product_id INT REFERENCES products(id) DELETE ON CASCADE,                    -- Foreign key referencing products
    amount INT NOT NULL,                        -- Amount of the product in the cart
    total_price DECIMAL(10, 2) NOT NULL,        -- Total price of the cart
    VariantIDs INT[],                       -- Array of variant IDs
    VariantOptionIDs INT[],                 -- Array of variant option IDs
	status status_enum DEFAULT 'active',
  	created_at TIMESTAMP DEFAULT NOW(),
  	updated_at TIMESTAMP,
  	deleted_at TIMESTAMP
);



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