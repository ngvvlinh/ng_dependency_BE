CREATE UNIQUE INDEX ON "fb_external_post"(external_id);
CREATE UNIQUE INDEX ON "fb_external_comment"(external_id);
CREATE UNIQUE INDEX ON "fb_external_conversation"(external_id);
CREATE UNIQUE INDEX ON "fb_external_message"(external_id);
CREATE UNIQUE INDEX ON "fb_customer_conversation"(external_id, external_user_id);