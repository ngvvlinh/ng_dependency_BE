alter table fb_external_comment add column internal_source int2;
alter table fb_external_message add column internal_source int2;

alter table history.fb_external_comment add column internal_source int2;
alter table history.fb_external_message add column internal_source int2;
