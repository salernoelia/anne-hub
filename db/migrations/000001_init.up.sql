-- Active: 1723672368124@@127.0.0.1@5432@anne_hub

create table users (
  id bigint primary key generated always as identity,
  username text not null unique,
  email text not null unique,
  password_hash text not null,
  created_at timestamptz default now()
);

create table devices (
  id bigint primary key generated always as identity,
  user_id bigint references users (id) on delete cascade,
  device_name text not null,
  device_type text,
  created_at timestamptz default now()
);

create table tasks (
  id bigint primary key generated always as identity,
  user_id bigint references users (id) on delete cascade,
  title text not null,
  description text,
  due_date timestamptz,
  completed boolean default false,
  created_at timestamptz default now()
);

alter table devices
drop device_type,
add column last_synced timestamptz,
add column companion_app_id bigint;

create table companion_apps (
  id bigint primary key generated always as identity,
  app_name text not null,
  settings jsonb,
  created_at timestamptz default now()
);

alter table devices
add constraint fk_companion_app foreign key (companion_app_id) references companion_apps (id) on delete set null;

alter table companion_apps
drop app_name;

alter table companion_apps
add column user_id bigint references users (id) on delete cascade;

alter table users
add column age int,
add column interests text[0];

alter table tasks
add column interest_links text[0];