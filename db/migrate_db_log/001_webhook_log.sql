-- Create new database `etop_log`
-- Run sql below in this database 
CREATE DATABASE etop_log;

CREATE TYPE shipping_provider AS ENUM (
    'ghn',
    'ghtk',
		'vtpost'
);

CREATE TABLE shipping_provider_webhook (
	id BIGINT PRIMARY KEY,
	shipping_provider shipping_provider,
	data jsonb,
	created_at TIMESTAMP WITH TIME ZONE,
	updated_at TIMESTAMP WITH TIME ZONE,
	shipping_code TEXT,
	shipping_state TEXT,
	external_shipping_state TEXT
);
