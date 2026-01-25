-- Migration: Change title and sub_title to TEXT type
-- Created: 2026-01-25
-- Purpose: Allow longer content for title and sub_title fields (support HTML content)

-- Change title column from VARCHAR(255) to TEXT
ALTER TABLE writing_part_list
MODIFY COLUMN title TEXT NOT NULL COMMENT '标题';

-- Change sub_title column from VARCHAR(255) to TEXT
ALTER TABLE writing_part_list
MODIFY COLUMN sub_title TEXT NULL COMMENT '副标题';

-- Note: TEXT type can store up to 65,535 characters, which is sufficient for HTML content
