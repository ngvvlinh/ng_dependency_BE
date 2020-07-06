UPDATE account_user
SET roles = '{admin}'
WHERE account_id = 101 AND roles IS NULL;
