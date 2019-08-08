ALTER TABLE "user"
  ADD COLUMN agreed_tos_at TIMESTAMPTZ,
  ADD COLUMN agreed_email_info_at TIMESTAMPTZ,
  ADD COLUMN email_verified_at TIMESTAMPTZ,
  ADD COLUMN phone_verified_at TIMESTAMPTZ,
  ADD COLUMN email_verification_sent_at TIMESTAMPTZ,
  ADD COLUMN phone_verification_sent_at TIMESTAMPTZ;

ALTER TABLE history."user"
  ADD COLUMN agreed_tos_at TIMESTAMPTZ,
  ADD COLUMN agreed_email_info_at TIMESTAMPTZ,
  ADD COLUMN email_verified_at TIMESTAMPTZ,
  ADD COLUMN phone_verified_at TIMESTAMPTZ,
  ADD COLUMN email_verification_sent_at TIMESTAMPTZ,
  ADD COLUMN phone_verification_sent_at TIMESTAMPTZ;
