CREATE TABLE invitation (
    id INT8 PRIMARY KEY,
    account_id INT8 NOT NULL REFERENCES shop(id),
    email TEXT NOT NULL,
    roles TEXT[],
    token TEXT,
    status INT2,
    invited_by INT8 NOT NULL REFERENCES "user"(id),
    accepted_at TIMESTAMP WITH TIME ZONE,
    rejected_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

select init_history('invitation', '{id,account_id}');
CREATE INDEX ON invitation(token);

-- add shop_owner for each record account_user
update account_user
set roles = array_append(roles, 'owner')
from
	account
where
	account_user.account_id = account.id
	and account.owner_id = account_user.user_id;