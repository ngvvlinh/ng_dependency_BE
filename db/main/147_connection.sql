ALTER TABLE connection
    ADD COLUMN origin_connection_id INT8 REFERENCES connection(id);

ALTER TABLE history.connection
    ADD COLUMN origin_connection_id INT8;
