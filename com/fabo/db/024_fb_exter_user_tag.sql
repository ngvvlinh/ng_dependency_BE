create table fb_shop_tag (
    id int8 primary key,
    name text not null,
    color text not null,
    shop_id bigint REFERENCES "shop"(id),
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
);

alter table fb_external_user add tags integer[];
