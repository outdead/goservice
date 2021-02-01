CREATE DATABASE migrago WITH TEMPLATE = template0 OWNER = postgres LC_COLLATE = 'en_US.UTF-8' LC_CTYPE = 'en_US.UTF-8';

\connect migrago

create table migration (
  project varchar not null,
  database varchar not null,
  version varchar not null,
  apply_time bigint default 0 not null,
  rollback boolean default true not null,
  constraint migration_pk primary key (project, database, version)
);
