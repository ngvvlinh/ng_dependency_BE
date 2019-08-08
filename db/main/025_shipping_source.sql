CREATE UNIQUE INDEX ON shipping_source (name, type);

UPDATE shipping_source SET name = 'D503809' WHERE type = 'ghn';
