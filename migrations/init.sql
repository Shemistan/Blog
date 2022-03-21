create database "blog";
create user "shem" with password '12345678';
grant all privileges on database "blog" to "shem";

\c blog shem

create table notes(
    id serial primary key,
    title varchar(20) not null,
    note_text text not null,
    tag varchar(20) not null,
    creating_data integer not null
);

alter table notes owner to shem;