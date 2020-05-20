UPDATE "order" SET customer_phone = customer_address->>'phone' WHERE customer_phone IS NULL;
