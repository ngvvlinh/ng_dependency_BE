UPDATE shop SET try_on = 'none' WHERE try_on IS NULL AND ghn_note_code ='KHONGCHOXEMHANG';
UPDATE shop SET try_on = 'open' WHERE try_on IS NULL AND ghn_note_code ='CHOXEMHANGKHONGTHU';
UPDATE shop SET try_on = 'try'  WHERE try_on IS NULL AND ghn_note_code ='CHOTHUHANG';
UPDATE shop SET try_on = 'open' WHERE try_on IS NULL;
ALTER TABLE shop ALTER COLUMN try_on SET NOT NULL;
