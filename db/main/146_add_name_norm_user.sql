 ALTER TABLE "user" ADD COLUMN full_name_norm tsvector;
 ALTER TABLE history."user" ADD COLUMN full_name_norm tsvector;

 CREATE INDEX ON "user" USING gin(full_name_norm);
