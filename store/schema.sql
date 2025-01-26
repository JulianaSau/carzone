-- Drop existing foreign key constraint
-- ALTER TABLE car
-- DROP CONSTRAINT IF EXISTS fk_engine_id;

-- Truncate car table to clear existing data
-- TRUNCATE TABLE car;

-- Truncate engine table to clear existing data
-- TRUNCATE TABLE engine;

-- Truncate users table to clear existing data
-- TRUNCATE TABLE "user";

-- Create user table
CREATE TABLE IF NOT EXISTS "user" (
    id UUID PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    phone_number VARCHAR(20),
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'manager', 'driver')),
    active BOOLEAN DEFAULT TRUE,
    created_by VARCHAR(50) NULLABLE DEFAULT NULL,
    deleted_at TIMESTAMP NULLABLE DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create engine table
CREATE TABLE IF NOT EXISTS engine (
    id UUID PRIMARY KEY,
    displacement INT NOT NULL,
    no_of_cylinders INT NOT NULL,
    car_range INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create car table
CREATE TABLE IF NOT EXISTS car (
    id UUID PRIMARY KEY,
    registration_number VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    year VARCHAR(4) NOT NULL,
    brand VARCHAR(255) NOT NULL,
    fuel_type VARCHAR(50) NOT NULL,
    engine_id UUID NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('Available', 'In Use', 'Maintenance', 'Decommissioned')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add foreign key constraint on engine_id in car table
ALTER TABLE car
ADD CONSTRAINT fk_engine_id
FOREIGN KEY (engine_id)
REFERENCES engine(id)
ON DELETE CASCADE;

-- Insert dummy data into the engine table
INSERT INTO engine (id, displacement, no_of_cylinders, car_range)
VALUES
    ('e1f86b1a-0873-4c19-bae2-fc60329d0140', 2000, 4, 600),
    ('f4a9c66b-8e38-419b-93c4-215d5cefb318', 1600, 4, 550),
    ('cc2c2a7d-2e21-4f59-b7b8-bd9e5e4cf04c', 3000, 6, 700),
    ('9746be12-07b7-42a3-b8ab-7d1f209b63d7', 1800, 4, 500);

-- Insert dummy data into the user table
INSERT INTO "user" (uuid, username, password, first_name, last_name, email, phone_number, role, created_by)
VALUES
    ('d3b07384-d9a1-4c4b-8a0d-4b1b1b1b1b1b', 'admin', 'admin123', 'System', 'Admin', 'admin@carmanagement.com', '244707070707', 'admin', NULL),
    ('e4c2f3a5-e5b2-4d5c-9b2e-5c2c2c2c2c2c', 'manager', 'manager123', 'System', 'Manager', 'manager@carmanagement.com', '244707070706', 'manager', NULL),
    ('f5d3e4b6-f6c3-4e6d-ac3f-6d3d3d3d3d3d', 'driver', 'driver123', 'System', 'Driver', 'driver@carmanagement.com', '244707070708', 'driver', NULL);

-- Insert dummy data into the car table
INSERT INTO car (id, registration_number, name, year, brand, fuel_type, engine_id, status, price)
VALUES
    ('c7c1a6d5-1ec4-4c64-a59a-8a2f6f3d2bf3', 'KCX 786T', 'Honda Civic', '2023', 'Honda', 'Gasoline', 'e1f86b1a-0873-4c19-bae2-fc60329d0140', 'Available', 25000.00),
    ('9d6a56f8-79c3-4931-a5c0-6b290c84ba2f', 'KCZ 883J', 'Toyota Corolla', '2022', 'Toyota', 'Gasoline', 'f4a9c66b-8e38-419b-93c4-215d5cefb318', 'Available', 22000.00),
    ('9b9437c4-3ed1-45a5-b240-0fe3e24e0e4e', 'KBX 284P', 'Ford Mustang', '2024', 'Ford', 'Gasoline', 'cc2c2a7d-2e21-4f59-b7b8-bd9e5e4cf04c', 'Available', 40000.00),
    ('5e9df51a-8d7a-4d84-9c58-4ccfe5c7db06', 'KDC 376C', 'BMW 3 Series', '2023', 'BMW', 'Gasoline', '9746be12-07b7-42a3-b8ab-7d1f209b63d7', 'Available', 35000.00);
