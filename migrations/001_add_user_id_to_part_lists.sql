-- Migration: Add user_id column to part list tables
-- Created: 2026-01-18

-- Add user_id column to listening_part_list table
ALTER TABLE listening_part_list
ADD COLUMN user_id VARCHAR(255) NULL
COMMENT '创建者用户ID';

-- Add user_id column to reading_part_list table
ALTER TABLE reading_part_list
ADD COLUMN user_id VARCHAR(255) NULL
COMMENT '创建者用户ID';

-- Add user_id column to writing_part_list table
ALTER TABLE writing_part_list
ADD COLUMN user_id VARCHAR(255) NULL
COMMENT '创建者用户ID';

-- Add index for better query performance
CREATE INDEX idx_listening_part_user_id ON listening_part_list(user_id);
CREATE INDEX idx_reading_part_user_id ON reading_part_list(user_id);
CREATE INDEX idx_writing_part_user_id ON writing_part_list(user_id);
