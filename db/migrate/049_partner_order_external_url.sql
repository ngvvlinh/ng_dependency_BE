ALTER TABLE partner
  ADD COLUMN recognized_hosts TEXT[],
  ADD COLUMN redirect_urls TEXT[];

ALTER TABLE history.partner
  ADD COLUMN recognized_hosts TEXT[],
  ADD COLUMN redirect_urls TEXT[];

ALTER TABLE shop ADD COLUMN recognized_hosts TEXT[];
ALTER TABLE history.shop ADD COLUMN recognized_hosts TEXT[];

ALTER TABLE "order" ADD COLUMN external_url TEXT;
ALTER TABLE history."order" ADD COLUMN external_url TEXT;
