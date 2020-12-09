CREATE TABLE shop_setting (
  shop_id INT8 REFERENCES shop(id),
  return_address_id INT8 REFERENCES address(id),
  try_on try_on,
  shipping_note TEXT,
  payment_type_id INT8,
  weight INT,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ
);
