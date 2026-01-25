-- Migration: Add task_type field and migrate data
-- Created: 2026-01-25
-- Purpose:
--   1. Add task_type field to distinguish Task 1 / Task 2
--   2. Migrate existing 'type' data to 'task_type'
--   3. Migrate existing 'source' data to 'type' (to align with other tables)

-- Step 1: Add new task_type column
ALTER TABLE writing_part_list
ADD COLUMN task_type VARCHAR(255) NULL COMMENT '写作任务类型：1=Task1，2=Task2';

-- Step 2: Copy existing 'type' values to 'task_type'
-- Current 'type' field contains "1" or "2" representing Task 1/2
UPDATE writing_part_list
SET task_type = type;

-- Step 3: Copy existing 'source' values to 'type'
-- Current 'source' field contains "1", "2", or "3" representing data source
UPDATE writing_part_list
SET type = source;

-- Step 4: Set source to NULL or keep for future use
-- Option 1: Set to NULL (recommended if not needed)
UPDATE writing_part_list
SET source = NULL;

-- Option 2: Keep source field for actual source information (e.g., "Cambridge IELTS 15")
-- In this case, do nothing here

-- Add index for task_type
CREATE INDEX idx_writing_part_task_type ON writing_part_list(task_type);

-- Note: After this migration:
-- - 'type' field: data source (1=system, 2=official, 3=user) - aligned with other tables
-- - 'task_type' field: task category (1=Task1, 2=Task2) - new field
-- - 'source' field: NULL or repurposed for actual source information
