alter table fb_external_comment
    add column is_hidden boolean,
    add column is_liked boolean,
    add column is_private_replied boolean;

alter table history.fb_external_comment
    add column is_hidden boolean,
    add column is_liked boolean,
    add column is_private_replied boolean;
