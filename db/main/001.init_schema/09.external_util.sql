CREATE VIEW public.external_account_ahamove_view AS
SELECT external_account_ahamove.id,
       external_account_ahamove.owner_id,
       external_account_ahamove.phone,
       external_account_ahamove.name,
       external_account_ahamove.external_verified,
       external_account_ahamove.created_at,
       external_account_ahamove.updated_at,
       external_account_ahamove.external_created_at,
       external_account_ahamove.last_send_verified_at,
       external_account_ahamove.external_ticket_id,
       external_account_ahamove.external_id,
       (external_account_ahamove.id_card_front_img IS NOT NULL) AS id_card_front_img_uploaded,
       (external_account_ahamove.id_card_back_img IS NOT NULL) AS id_card_back_img_uploaded,
       (external_account_ahamove.portrait_img IS NOT NULL) AS portrait_img_uploaded,
       external_account_ahamove.uploaded_at
FROM public.external_account_ahamove;

CREATE VIEW public.partner_relation_view AS
SELECT partner_relation.partner_id,
       partner_relation.subject_id,
       partner_relation.subject_type,
       partner_relation.status,
       partner_relation.created_at,
       partner_relation.updated_at,
       partner_relation.deleted_at
FROM public.partner_relation;

-- utilities for metabase

CREATE FUNCTION public.convertinterval2hoursminutes(t interval) RETURNS text
  LANGUAGE plpgsql
AS $$
BEGIN
  return concat(EXTRACT(epoch FROM t)::int/3600 , ':', (EXTRACT(epoch FROM t)::int - EXTRACT(epoch FROM t)::int/3600*3600)/60);
END;
$$;

CREATE FUNCTION public.convertinterval2hoursminutes(t1 timestamp without time zone, t2 timestamp without time zone) RETURNS text
  LANGUAGE plpgsql
AS $$
DECLARE t INTERVAL;
BEGIN
  t = t1::INTERVAL - t2::INTERVAL;
  return concat(EXTRACT(epoch FROM t::INTERVAL)::int/3600 , ':', (EXTRACT(epoch FROM t::INTERVAL)::int - EXTRACT(epoch FROM t::INTERVAL)::int/3600*3600)/60);
END;
$$;

CREATE FUNCTION public.convertinterval2hoursminutes(t1 timestamp with time zone, t2 timestamp with time zone) RETURNS text
  LANGUAGE plpgsql
AS $$
DECLARE t INTERVAL;
BEGIN
  t = t1 - t2;
  return concat(EXTRACT(epoch FROM t::INTERVAL)::int/3600 , ':', (EXTRACT(epoch FROM t::INTERVAL)::int - EXTRACT(epoch FROM t::INTERVAL)::int/3600*3600)/60);
END;
$$;

