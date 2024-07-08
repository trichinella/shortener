CREATE TABLE public.shortcuts (
                                  uuid uuid NOT NULL,
                                  original_url varchar NOT NULL,
                                  short_url varchar NOT NULL,
                                  created_date timestamp with time zone DEFAULT (now() at time zone 'Europe/Moscow') NOT NULL,
                                  CONSTRAINT shortcuts_pk PRIMARY KEY (uuid)
);
COMMENT ON TABLE public.shortcuts IS 'Сокращения ссылок';

-- Column comments

COMMENT ON COLUMN public.shortcuts.uuid IS 'Уникальный идентификатор (PK)';
COMMENT ON COLUMN public.shortcuts.original_url IS 'Оригинальная (пользовательская) ссылка';
COMMENT ON COLUMN public.shortcuts.short_url IS 'Сокращенная ссылка (хэш)';
COMMENT ON COLUMN public.shortcuts.created_date IS 'Дата создания сокращения';
