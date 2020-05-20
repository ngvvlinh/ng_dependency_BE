alter table receipt
    rename column user_id to created_by;

alter table history.receipt
    rename column user_id to created_by;
