create table conversations (
  id bigint primary key generated always as identity,
  user_id uuid references users (id) on delete cascade,
  conversation_id uuid default gen_random_uuid (),
  created_at timestamptz default now(),
  request text not null,
  response text not null,
  model_used text,
  role text
);

