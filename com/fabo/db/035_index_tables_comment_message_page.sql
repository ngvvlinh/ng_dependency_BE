create index on fb_external_message(external_created_time, created_by);
create index on fb_external_comment(external_created_time, created_by);
create index on fb_external_page(external_id, shop_id);
