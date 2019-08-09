INSERT INTO "user" (
    "id", "status", "created_at", "updated_at", "full_name", "short_name", "email", "phone", "is_test", "identifying"
) VALUES (
    '1000101010101010101', 1, NOW(), NOW(), 'Etop System', 'System', 'hi@etop.vn', '0101010101', 0, 'full'
);
INSERT INTO "account" (
    "id", "name", "type", "owner_id"
) VALUES (
    '1000421281650350414', 'Haravan', 'partner', '1000101010101010101'
);

INSERT INTO "partner" (
    "id", "name", "public_name", "owner_id", "status", "is_test", "phone", "email", "created_at", "updated_at"
) VALUES (
    '1000421281650350414', 'Haravan', 'Haravan', '1000101010101010101', '1', '0', '0903119101', 'hi@haravan.com', NOW(), NOW()
);

ALTER TABLE "order" ADD COLUMN external_meta JSONB;
ALTER TABLE history."order" ADD COLUMN  external_meta JSONB;
