ALTER TABLE public.shortcuts ADD user_id uuid NULL;
COMMENT ON COLUMN public.shortcuts.user_id IS 'ID пользователя, создавшего сокращение';
CREATE INDEX shortcuts_user_id_idx ON public.shortcuts (user_id);