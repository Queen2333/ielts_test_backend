-- Rollback Migration: Remove user_id column from part list tables
-- Created: 2026-01-18

-- Remove index first
DROP INDEX IF EXISTS idx_listening_part_user_id ON listening_part_list;
DROP INDEX IF EXISTS idx_reading_part_user_id ON reading_part_list;
DROP INDEX IF EXISTS idx_writing_part_user_id ON writing_part_list;

-- Remove user_id column from listening_part_list table
ALTER TABLE listening_part_list
DROP COLUMN user_id;

-- Remove user_id column from reading_part_list table
ALTER TABLE reading_part_list
DROP COLUMN user_id;

-- Remove user_id column from writing_part_list table
ALTER TABLE writing_part_list
DROP COLUMN user_id;
