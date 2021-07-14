ALTER TABLE public."contact" ADD COLUMN full_name_norm tsvector;
ALTER TABLE history."contact" ADD COLUMN full_name_norm tsvector;

CREATE INDEX ON  public."contact" USING GIN(full_name_norm);
