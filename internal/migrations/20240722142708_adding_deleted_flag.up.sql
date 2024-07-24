ALTER TABLE public.shortcuts ADD deleted_date timestamptz NULL;
CREATE INDEX shortcuts_deleted_date_idx ON public.shortcuts (deleted_date);
