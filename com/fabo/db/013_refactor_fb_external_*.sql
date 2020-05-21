-- Remove id (table fb_external_user_internal and fb_external_user)
--   (1): Add external_id replace for id
--   (2): Map external_id between fb_external_user -> fb_external_user_internal
--   (3): Remove duplicate external_id (fb_external_user_internal)
--   (4): Remove id (fb_external_user_internal) and Set external_id to primary key
--   (5): Drop reference fb_user_id to table fb_external_user
--        Remove user_id
--   (6): Remove duplicate external_id in table fb_external_user
--   (7): Drop id,user_id and Set external_id to primary key
--   (8): Add foreign key external_id (fb_external_user_internal) to external_id (fb_external_user)


-- (1): Add external_id replace for id
alter table fb_external_user_internal
    add column external_id text;
alter table history.fb_external_user_internal
    add column external_id text;

-- (2): Map external_id between fb_external_user -> fb_external_user_internal
update fb_external_user_internal
set external_id  = fb_external_user.external_id
from fb_external_user
where fb_external_user_internal.id = fb_external_user.id;

-- (3): Remove duplicate external_id (fb_external_user_internal)
delete from fb_external_user_internal
where id not in (
	select min(id)
	from fb_external_user_internal
	group by external_id
);

-- (4): Remove id (fb_external_user_internal) and Set external_id to primary key
alter table fb_external_user_internal
	drop column id,
	add primary key ("external_id");
alter table history.fb_external_user_internal
	drop column id;

-- (5): Drop reference fb_user_id to table fb_external_user
--      Remove user_id
alter table fb_external_page
    drop column fb_user_id,
    drop column user_id;
alter table history.fb_external_page
    drop column fb_user_id,
    drop column user_id;

-- (6): Remove duplicate external_id in table fb_external_user
delete from fb_external_user
where id not in (
	select min(id)
	from fb_external_user
	group by external_id
);

-- (7): Drop id,user_id and Set external_id to primary key
alter table fb_external_user
	drop column id,
	drop column user_id,
	add primary key ("external_id");
alter table history.fb_external_user
	drop column id,
	drop column user_id;

-- (8): Add foreign key external_id (fb_external_user_internal) to external_id (fb_external_user)
alter table fb_external_user_internal
    add foreign key ("external_id") references fb_external_user(external_id);



-- Add external_id into fb_external_page_internal and map its
alter table fb_external_page_internal
    add column external_id text;
alter table history.fb_external_page_internal
    add column external_id text;

update fb_external_page_internal
set external_id  = fb_external_page.external_id
from fb_external_page
where fb_external_page_internal.id = fb_external_page.id;


-- Drop fb_page_id (unnecessary)
alter table fb_external_post
    drop column fb_page_id;
alter table history.fb_external_post
    drop column fb_page_id;

-- Drop fb_post_id, fb_page_id (unnecessary)
alter table fb_external_comment
    drop column fb_post_id,
    drop column fb_page_id;
alter table history.fb_external_comment
    drop column fb_post_id,
    drop column fb_page_id;

-- Drop fb_page_id (unnecessary)
alter table fb_external_conversation
    drop column fb_page_id;
alter table history.fb_external_conversation
    drop column fb_page_id;

-- Drop fb_conversation_id, fb_page_id (unnecessary)
alter table fb_external_message
    drop column fb_conversation_id,
    drop column fb_page_id;
alter table history.fb_external_message
    drop column fb_conversation_id,
    drop column fb_page_id;

-- Drop fb_page_id (unnecessary)
alter table fb_customer_conversation
    drop column fb_page_id;
alter table history.fb_customer_conversation
    drop column fb_page_id;

-- Remove duplicate external_id in table fb_external_page_internal and fb_external_page (where status = -1)
delete from fb_external_page_internal
where id in (
	select id
	from fb_external_page
	where status = -1
);

delete from fb_external_page
where status = -1;

CREATE UNIQUE INDEX ON "fb_external_page"(external_id);