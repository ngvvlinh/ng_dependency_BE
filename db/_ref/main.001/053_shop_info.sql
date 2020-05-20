ALTER TABLE shop
  ADD COLUMN company_info JSONB,
  ADD COLUMN money_transaction_rrule TEXT,
  ADD COLUMN survey_info JSONB,
  ADD COLUMN shipping_service_select_strategy JSONB;

ALTER TABLE history.shop
  ADD COLUMN company_info JSONB,
  ADD COLUMN money_transaction_rrule TEXT,
  ADD COLUMN survey_info JSONB,
  ADD COLUMN shipping_service_select_strategy JSONB;
