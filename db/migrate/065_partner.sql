ALTER TABLE partner 
  ADD COLUMN available_from_etop boolean,
  ADD COLUMN available_from_etop_config JSONB;

ALTER TABLE history.partner 
  ADD COLUMN available_from_etop boolean,
  ADD COLUMN available_from_etop_config JSONB;
