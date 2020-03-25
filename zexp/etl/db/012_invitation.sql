CREATE TABLE IF NOT EXISTS invitation (
    id INT8 PRIMARY KEY,
    account_id INT8,
    email TEXT,
    roles TEXT[],
    token TEXT,
    status INT2,
    invited_by INT8,
    accepted_at TIMESTAMP WITH TIME ZONE,
    rejected_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    full_name TEXT,
    short_name TEXT,
    position TEXT,
    phone TEXT,
    rid INT8
);