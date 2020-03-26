alter table ws_website add column site_subdomain text;

CREATE UNIQUE INDEX  ON "ws_website" (site_subdomain);
