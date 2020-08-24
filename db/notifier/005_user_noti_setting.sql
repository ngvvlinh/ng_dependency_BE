CREATE TABLE user_noti_setting (
    user_id int8 PRIMARY KEY,
    disable_topics text[]
);

ALTER TABLE notification ADD COLUMN user_id int8;
