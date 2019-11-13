alter table inventory_voucher add column trader jsonb;

alter table shop_trader
    add column deleted_at TIMESTAMPTZ;

alter table history.shop_trader
    add column deleted_at TIMESTAMPTZ;

update shop_trader as st set deleted_at = sc.deleted_at
from shop_customer as sc
where st.id = sc.id and sc.deleted_at is not null;

update shop_trader as st set deleted_at = sc.deleted_at
from shop_carrier as sc
where st.id = sc.id and sc.deleted_at is not null;

update shop_trader as st set deleted_at = sp.deleted_at
from shop_supplier as sp
where st.id = sp.id and sp.deleted_at is not null;
