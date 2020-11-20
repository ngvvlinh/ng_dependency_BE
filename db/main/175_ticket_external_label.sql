alter table ticket_label
    add constraint ticket_label_constraint primary key(id);

create table ticket_label_external (
   id int8 primary key,
   connection_id int8 references connection(id),
   external_id text,
   external_name text,
   created_at timestamptz,
   updated_at timestamptz,
   deleted_at timestamptz
);
select init_history('ticket_label_external', '{id}');

create table ticket_label_ticket_label_external (
    ticket_label_id int8 references ticket_label(id),
    ticket_label_external_id int8 references ticket_label_external(id),
    deleted_at timestamptz
);
select init_history('ticket_label_ticket_label_external', '{ticket_label_id, ticket_label_external_id}');

alter table ticket
    add column external_id text,
    add column wl_partner_id int8 references partner(id),
    add column deleted_at timestamptz,
    add column connection_id int8 references connection(id);
alter table history.ticket
    add column external_id text,
    add column wl_partner_id int8,
    add column deleted_at timestamptz,
    add column connection_id int8;

alter table ticket_comment
    add column external_id text,
    add column external_created_at timestamptz;
alter table history.ticket_comment
    add column external_id text,
    add column external_created_at timestamptz;

alter table ticket_label
    add column created_at timestamptz,
    add column updated_at timestamptz,
    add column deleted_at timestamptz,
    add column wl_partner_id int8 references partner(id);
