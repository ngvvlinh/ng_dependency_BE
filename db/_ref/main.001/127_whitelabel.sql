alter table partner
  add column white_label_key text;

alter table history.partner
  add column white_label_key text;

UPDATE partner SET white_label_key = 'itopx' WHERE id = '1057192413421863086';
