alter table invitation
    add column phone text;

alter table "history".invitation
    add column phone text;

alter table invitation alter column email drop not null;
