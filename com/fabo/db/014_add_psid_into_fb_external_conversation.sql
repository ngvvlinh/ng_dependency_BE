alter table fb_external_conversation
    add column psid text not null;
alter table history.fb_external_conversation
    add column psid text;