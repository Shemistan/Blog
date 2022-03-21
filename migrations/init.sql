create database "blog";
create user "shem" with password '12345678';
grant all privileges on database "blog" to "shem";

\c blog shem

create table notes(
    id serial primary key,
    title varchar(20) not null,
    note_text text not null,
    tag varchar(20) not null,
    created_at timestamp not null default now()
);

alter table notes owner to shem;