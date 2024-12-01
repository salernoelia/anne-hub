alter table devices
drop constraint devices_user_id_fkey;

alter table tasks
drop constraint tasks_user_id_fkey;

alter table companion_apps
drop constraint companion_apps_user_id_fkey;

alter table users
drop constraint users_pkey;

alter table users
add column new_id uuid default gen_random_uuid ();

alter table devices
add column new_user_id uuid;

alter table tasks
add column new_user_id uuid;

alter table companion_apps
add column new_user_id uuid;

update devices
set
  new_user_id = users.new_id
from
  users
where
  devices.user_id = users.id;

update tasks
set
  new_user_id = users.new_id
from
  users
where
  tasks.user_id = users.id;

update companion_apps
set
  new_user_id = users.new_id
from
  users
where
  companion_apps.user_id = users.id;

alter table users
drop column id;

alter table users
rename column new_id to id;

alter table users
add primary key (id);

alter table devices
drop column user_id;

alter table devices
rename column new_user_id to user_id;

alter table devices
add constraint devices_user_id_fkey foreign key (user_id) references users (id) on delete cascade;

alter table tasks
drop column user_id;

alter table tasks
rename column new_user_id to user_id;

alter table tasks
add constraint tasks_user_id_fkey foreign key (user_id) references users (id) on delete cascade;

alter table companion_apps
drop column user_id;

alter table companion_apps
rename column new_user_id to user_id;

alter table companion_apps
add constraint companion_apps_user_id_fkey foreign key (user_id) references users (id) on delete cascade;