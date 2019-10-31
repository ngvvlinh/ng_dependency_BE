alter table shop_trader alter column shop_id drop not null;
alter table shop_carrier alter column shop_id drop not null;

-- Create carrier topship
insert into shop_trader (id, shop_id, "type")
    values (1000030662086749358, null, 'carrier');

insert into shop_carrier (id, shop_id, full_name, status, created_at, updated_at)
    values (1000030662086749358, null, 'TopShip', 0, now(), now());

alter table receipt alter column created_by drop not null;

-- Create independent customer
alter table shop_customer alter column shop_id drop not null;

insert into shop_trader (id, shop_id, "type")
    values (1000080135776788835, null, 'customer');

insert into shop_customer (id, shop_id, full_name, status, created_at, updated_at)
    values (1000080135776788835, null, 'Khách lẻ', 1, now(), now());