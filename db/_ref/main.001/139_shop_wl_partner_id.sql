alter table shop add column wl_partner_id int8 references partner(id);
alter table history.shop add column wl_partner_id int8;

-- migration
update shop set wl_partner_id = "user".wl_partner_id
    from "user" where shop.owner_id = "user".id;
