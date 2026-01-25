-- Migration: Drop source field from writing_part_list
-- Created: 2026-01-25
-- Purpose: Remove the source field as it's no longer needed after migration 003

-- Drop the source column
ALTER TABLE writing_part_list
DROP COLUMN source;

-- Note: This operation is irreversible. Make sure you have a backup if needed.
