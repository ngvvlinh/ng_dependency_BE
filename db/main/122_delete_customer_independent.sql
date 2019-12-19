ALTER TYPE customer_type ADD VALUE 'anonymous';

update "shop_customer"
set deleted_at = now()
where "type"='independent';
