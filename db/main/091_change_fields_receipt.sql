ALTER TABLE receipt
    RENAME COLUMN order_ids to ref_ids;
ALTER TABLE "history".receipt
    RENAME COLUMN order_ids to ref_ids;

ALTER TABLE receipt
    ADD COLUMN cancelled_reason TEXT;
ALTER TABLE "history".receipt
    ADD COLUMN cancelled_reason TEXT;

ALTER TABLE receipt
    ADD COLUMN code_norm int4;
ALTER TABLE "history".receipt
    ADD COLUMN code_norm int4;

ALTER TABLE receipt
    ADD COLUMN paid_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE "history".receipt
    ADD COLUMN paid_at TIMESTAMP WITH TIME ZONE;

ALTER TABLE receipt
    ADD COLUMN confirmed_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE "history".receipt
    ADD COLUMN confirmed_at TIMESTAMP WITH TIME ZONE;

ALTER TABLE receipt
    ADD COLUMN cancelled_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE "history".receipt
    ADD COLUMN cancelled_at TIMESTAMP WITH TIME ZONE;
