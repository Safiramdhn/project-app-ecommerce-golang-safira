-- Enum status_enum
CREATE TYPE status_enum AS ENUM ('active', 'deleted');

-- Tabel users
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
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
  category_id INT REFERENCES categories(id) ON DELETE SET NULL,
  price DECIMAL(10, 2) NOT NULL,
  discount DECIMAL(10, 2) DEFAULT 0,
  rating DECIMAL(1, 2) CHECK (rating BETWEEN 1 AND 5),
  photo_url TEXT,
  is_new_product BOOLEAN GENERATED ALWAYS AS (created_at > NOW() - INTERVAL '30 days') STORED,
  has_variations BOOLEAN DEFAULT FALSE,
  status status_enum DEFAULT 'active',
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);

-- Tabel wishlist
CREATE TABLE wishlist (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users(id) ON DELETE CASCADE,
  product_id INT REFERENCES products(id) ON DELETE CASCADE,
  status status_enum DEFAULT 'active',
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);

-- Tabel cart
CREATE TABLE cart (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users(id) ON DELETE CASCADE,
  product_id INT REFERENCES products(id) ON DELETE CASCADE,
  amount INT DEFAULT 1 CHECK (amount > 0),
  status status_enum DEFAULT 'active',
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);

-- Tabel addresses
CREATE TABLE addresses (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users(id) ON DELETE CASCADE,
  is_default BOOLEAN DEFAULT FALSE,
  name VARCHAR NOT NULL,
  street VARCHAR NOT NULL,
  district VARCHAR,
  country VARCHAR,
  status status_enum DEFAULT 'active',
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);

-- Tabel shipping_types
CREATE TABLE shipping_types (
  id SERIAL PRIMARY KEY,
  name VARCHAR NOT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  status status_enum DEFAULT 'active',
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);

-- Tabel payment_methods
CREATE TABLE payment_methods (
  id SERIAL PRIMARY KEY,
  name VARCHAR NOT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  status status_enum DEFAULT 'active',
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);

-- Tabel order_items
CREATE TABLE order_items (
  id SERIAL PRIMARY KEY,
  product_id INT REFERENCES products(id) ON DELETE CASCADE,
  amount INT CHECK (amount > 0),
  price DECIMAL(10, 2),
  order_id INT REFERENCES orders(id) ON DELETE CASCADE,
  status status_enum DEFAULT 'active',
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);

-- Tabel orders
CREATE TABLE orders (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users(id) ON DELETE CASCADE,
  total_price DECIMAL(10, 2),
  address_id INT REFERENCES addresses(id),
  shipping_type_id INT REFERENCES shipping_types(id),
  payment_method_id INT REFERENCES payment_methods(id),
  status status_enum DEFAULT 'active',
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);

-- Tabel recommendations
CREATE TABLE recommendations (
  id SERIAL PRIMARY KEY,
  product_id INT REFERENCES products(id) ON DELETE CASCADE,
  title VARCHAR,
  subtitle VARCHAR,
  photo_url TEXT,
  status status_enum DEFAULT 'active',
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);

CREATE TABLE variations (
    id SERIAL PRIMARY KEY,
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    attribute_name VARCHAR NOT NULL, -- E.g., Size, Color
    created_at TIMESTAMP DEFAULT NOW()
);

-- Table Variation_Options (only for products with variations)
CREATE TABLE variation_options (
    id SERIAL PRIMARY KEY,
    variation_id INT REFERENCES variations(id) ON DELETE CASCADE,
    option_value VARCHAR NOT NULL, -- E.g., Red, Large
    additional_price DECIMAL(10, 2) DEFAULT 0,
    stock INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);
