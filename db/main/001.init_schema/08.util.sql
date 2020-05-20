-- for generating next shipping_code for vtpost
-- TODO: rename to shipping_code_vtpost
CREATE SEQUENCE public.shipping_code
  START WITH 100001
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;
