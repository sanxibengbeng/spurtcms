-- 初始化本地存储设置
-- 这个脚本将在数据库初始化后运行，设置 SpurtCMS 使用本地存储

-- 检查表是否存在，如果不存在则创建
CREATE TABLE IF NOT EXISTS tbl_storage_type (
    id SERIAL PRIMARY KEY,
    local VARCHAR(255),
    aws JSONB,
    azure JSONB,
    drive JSONB,
    selected_type VARCHAR(50)
);

-- 检查是否已有记录，如果没有则插入
INSERT INTO tbl_storage_type (id, local, selected_type)
SELECT 1, '/app/storage', 'local'
WHERE NOT EXISTS (SELECT 1 FROM tbl_storage_type WHERE id = 1);

-- 如果记录已存在，则更新为使用本地存储
UPDATE tbl_storage_type
SET selected_type = 'local', local = '/app/storage'
WHERE id = 1 AND selected_type != 'local';
