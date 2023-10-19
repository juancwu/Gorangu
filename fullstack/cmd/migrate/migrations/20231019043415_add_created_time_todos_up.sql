-- Write your UP migration SQL here.
ALTER TABLE todos
ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
