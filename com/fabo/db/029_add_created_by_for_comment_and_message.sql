alter table fb_external_message add column created_by int8;
alter table history.fb_external_message add column created_by int8;

alter table fb_external_comment add column created_by int8;
alter table history.fb_external_comment add column created_by int8;
