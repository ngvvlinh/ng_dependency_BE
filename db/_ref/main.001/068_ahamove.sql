ALTER TABLE external_account_ahamove
    ADD COLUMN fanpage_url TEXT,
    ADD COLUMN business_license_imgs TEXT[],
    ADD COLUMN website_url TEXT,
    ADD COLUMN company_imgs TEXT[],
    ADD COLUMN external_data_verified JSONB;
