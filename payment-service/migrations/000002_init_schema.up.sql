-- 2025 MAY 05 UPDATE: ADDED CARTS TABLE
CREATE TABLE IF NOT EXISTS cart (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  student_id TEXT NOT NULL,
  course_id TEXT NOT NULL,
  price NUMERIC(10, 2) NOT NULL,
  course_name TEXT NOT NULL,
  added_at TIMESTAMP DEFAULT NOW()
);

-- Add a unique constraint to ensure no duplicate course entries for the same student
ALTER TABLE cart ADD CONSTRAINT unique_cart_item UNIQUE (student_id, course_id);