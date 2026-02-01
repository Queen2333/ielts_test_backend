-- Migration: Add audio_files field to listening_part_list
-- Created: 2026-01-25
-- Purpose: Store audio file URLs/paths for listening parts

ALTER TABLE listening_part_list
ADD COLUMN audio_files TEXT NULL COMMENT '音频文件列表，JSON格式存储多个音频文件路径';

-- Note: This field will store JSON array of audio file URLs
-- Example: ["http://example.com/audio1.mp3", "http://example.com/audio2.mp3"]
