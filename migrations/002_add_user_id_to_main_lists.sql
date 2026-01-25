-- Migration: Add user_id column to main list tables
-- Created: 2026-01-25

-- Add user_id column to listening_list table
ALTER TABLE listening_list
ADD COLUMN user_id VARCHAR(255) NULL
COMMENT '创建者用户ID';

-- Add user_id column to reading_list table
ALTER TABLE reading_list
ADD COLUMN user_id VARCHAR(255) NULL
COMMENT '创建者用户ID';

-- Add user_id column to writing_list table
ALTER TABLE writing_list
ADD COLUMN user_id VARCHAR(255) NULL
COMMENT '创建者用户ID';

-- Add user_id column to testing_list table
ALTER TABLE testing_list
ADD COLUMN user_id VARCHAR(255) NULL
COMMENT '创建者用户ID';

-- Add index for better query performance
CREATE INDEX idx_listening_list_user_id ON listening_list(user_id);
CREATE INDEX idx_reading_list_user_id ON reading_list(user_id);
CREATE INDEX idx_writing_list_user_id ON writing_list(user_id);
CREATE INDEX idx_testing_list_user_id ON testing_list(user_id);
