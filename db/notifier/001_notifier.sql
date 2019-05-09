CREATE DATABASE etop_notifier;

CREATE TABLE external_service (
  id INT2 PRIMARY KEY
, code TEXT NOT NULL UNIQUE
, name TEXT NOT NULL
);

INSERT INTO external_service (id, code, name) VALUES
(101, 'onesignal', 'One Signal');

CREATE TABLE notification (
  id INT8 PRIMARY KEY
, rid SERIAL
, title TEXT
, message TEXT
, is_read BOOLEAN
, data JSONB
, entity TEXT
, entity_id INT8
, account_id INT8
, sync_status INT2
, external_service_id INT2
, external_noti_id TEXT
, seen_at TIMESTAMP WITH TIME ZONE
, created_at TIMESTAMP WITH TIME ZONE
, updated_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE device (
  id INT8 PRIMARY KEY
, device_id TEXT UNIQUE
, device_name TEXT
, external_service_id INT2
, external_device_id TEXT UNIQUE
, account_id INT8
, created_at TIMESTAMP WITH TIME ZONE
, updated_at TIMESTAMP WITH TIME ZONE
);

CREATE FUNCTION notify_pgrid_insert() RETURNS trigger
  AS $$
BEGIN
  PERFORM pg_notify(
    'pgrid',
    TG_TABLE_NAME
    || ':' || NEW.rid
    || ':' || 'INSERT'
    || ':' || NEW.id
  );
  RETURN NULL;
END$$ LANGUAGE 'plpgsql';

CREATE TRIGGER notify_pgrid_insert AFTER INSERT ON notification
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_insert();
