ALTER TABLE "order"
    ADD COLUMN referral_meta JSONB;
ALTER TABLE history."order"
    ADD COLUMN referral_meta JSONB;
