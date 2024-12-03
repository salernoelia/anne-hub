ALTER TABLE conversations
DROP COLUMN if exists request,
DROP COLUMN if exists response,
DROP COLUMN if exists role,
DROP COLUMN if exists conversation_history;

ALTER TABLE conversations ADD COLUMN conversation_history jsonb;



