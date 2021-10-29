create table public.todos
(
    id          serial,
    task        text not null,
    completed   bool      default false,
    created_at  timestamp default now(),
    completed_at timestamp
);

comment
on table public.todos is 'todos v1 20211028';

create unique index todos_id_uindex on public.todos (id);

alter table public.todos
    add constraint todos_pk primary key (id);

