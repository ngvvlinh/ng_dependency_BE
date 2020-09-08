ALTER TABLE ticket ADD COLUMN ref_ticket_id int8 REFERENCES ticket(id);

CREATE INDEX ticket_ref_ticket_id_idx
ON public.ticket (ref_ticket_id);
