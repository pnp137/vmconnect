-- VMConnect Database Seed Data
-- This file contains sample test data for frontend development
-- 
-- To load this data:
-- 1. Start the backend: docker compose -f build/docker-compose.yml up -d
-- 2. Import this file: cat database_seed.sql | docker exec -i vmconnect-db mysql -u vmconnect_user -pvmconnect_pass vmconnect

-- ==========================================
-- TEST USERS
-- ==========================================

-- Merchant User
-- Email: merchant@test.com
-- Password: should be hashed in production
INSERT INTO users (email, password, role) 
VALUES ('merchant@test.com', 'merchant123', 'MERCHANT');

-- Vendor User  
-- Email: vendor@test.com
INSERT INTO users (email, password, role)
VALUES ('vendor@test.com', 'vendor123', 'VENDOR');

-- ==========================================
-- SAMPLE MERCHANTS
-- ==========================================

-- Insert a test merchant
-- Note: Replace user_id with actual ID from users table
-- SELECT LAST_INSERT_ID() to get the ID after inserting a user
INSERT INTO merchants (user_id, name, email, phone, status)
VALUES (1, 'Fresh Mart', 'merchant@test.com', '9876543210', 'ACTIVE');

-- ==========================================
-- SAMPLE VENDORS
-- ==========================================

INSERT INTO vendors (user_id, name, email, phone, status, vendor_code)
VALUES (2, 'Premium Supplies', 'vendor@test.com', '9123456789', 'ACTIVE', 'VENDOR_001');

-- ==========================================
-- SAMPLE PRODUCTS
-- ==========================================

INSERT INTO products (vendor_id, name, description, price, category)
VALUES 
(1, 'Organic Tomato', 'Fresh organic tomatoes', 50.00, 'VEGETABLES'),
(1, 'Fresh Milk', '1 liter milk packet', 60.00, 'DAIRY'),
(1, 'Brown Bread', 'Whole wheat bread', 40.00, 'BAKERY');

-- ==========================================
-- LOADING INSTRUCTIONS
-- ==========================================

-- Step 1: Ensure backend is running
-- docker compose -f build/docker-compose.yml up -d

-- Step 2: Import this seed data
-- cat database_seed.sql | docker exec -i vmconnect-db mysql -u vmconnect_user -pvmconnect_pass vmconnect

-- Step 3: Verify data loaded
-- docker exec vmconnect-db mysql -u vmconnect_user -pvmconnect_pass vmconnect -e "SELECT * FROM users;"

-- ==========================================
-- TEST CREDENTIALS
-- ==========================================
-- Merchant Login:
--   Email: merchant@test.com
--   Password: merchant123

-- Vendor Login:
--   Email: vendor@test.com
--   Password: vendor123
