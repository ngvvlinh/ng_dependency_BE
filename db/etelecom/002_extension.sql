ALTER TABLE extension
    ADD COLUMN hotline_id INT8 REFERENCES hotline(id);

-- hotline builtin
INSERT INTO hotline (id, hotline, connection_id, connection_method, created_at, updated_at) VALUES (1161165496170941353, '19001008', 100085369475949390, 'builtin', NOW(), NOW());
