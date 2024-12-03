ALTER TABLE conversations
ADD COLUMN updated_at timestamptz default now();