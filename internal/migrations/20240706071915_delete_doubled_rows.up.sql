delete
from
    public.shortcuts
where
        uuid in (
        select
            distinct s2.uuid
        from
            public.shortcuts s2
                inner join public.shortcuts s on
                        s.original_url = s2.original_url
                    and s2.created_date>s.created_date);