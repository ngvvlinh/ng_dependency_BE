create table fb_shop_user_tag (
    id int8 primary key,
    name text not null,
    color text not null,
    shop_id int8 REFERENCES "shop"(id),
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
);

alter table fb_external_user add tag_ids int8[];
alter table history.fb_external_user add tag_ids int8[];

create index on fb_shop_user_tag(shop_id);
