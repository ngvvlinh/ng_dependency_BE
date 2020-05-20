CREATE OR REPLACE FUNCTION fulfillment_expected_pick_at(created_at timestamptz)
RETURNS timestamptz
LANGUAGE plpgsql IMMUTABLE
AS $$
DECLARE
  created_hour INT;
begin
  created_hour := date_part('hour', created_at at time zone 'ict');
  IF (created_hour < 10) THEN
    return date_trunc('hour', created_at + (13 - created_hour) * interval '1 hour'); -- pick at 13h
  ELSIF (created_hour < 16) THEN
    return date_trunc('hour', created_at + (21 - created_hour) * interval '1 hour'); -- pick at 21h
  ELSE
    return date_trunc('hour', created_at + (37 - created_hour) * interval '1 hour'); -- pick at 13h tomorrow
  END IF;
END
$$;

CREATE INDEX ON fulfillment (fulfillment_expected_pick_at(created_at));
