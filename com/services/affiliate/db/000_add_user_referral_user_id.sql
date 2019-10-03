alter table affiliate_referral_code add column user_id bigint;

update affiliate_referral_code arc
set user_id = aff.owner_id
from
    affiliate aff
where arc.affiliate_id = aff.id;

alter table affiliate_referral_code alter column user_id set not null;