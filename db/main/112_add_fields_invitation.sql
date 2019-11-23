alter table invitation
    add column "full_name" TEXT,
    add column "short_name" TEXT,
    add column "position" TEXT;

alter table "history".invitation
    add column "full_name" TEXT,
    add column "short_name" TEXT,
    add column "position" TEXT;