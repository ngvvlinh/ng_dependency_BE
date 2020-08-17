CREATE TABLE user_ref_saff (
    user_id BIGINT PRIMARY KEY REFERENCES public."user"(id),
    ref_sale text NULL,
    ref_aff text NULL
);

CREATE INDEX user_ref_saff_ref_sale_idx
ON public.user_ref_saff (ref_sale);

CREATE INDEX user_ref_saff_ref_aff_idx
ON public.user_ref_saff (ref_aff);
