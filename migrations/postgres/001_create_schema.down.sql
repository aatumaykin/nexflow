-- Drop indexes
DROP INDEX IF EXISTS idx_logs_created_at;
DROP INDEX IF EXISTS idx_logs_source;
DROP INDEX IF EXISTS idx_logs_level;
DROP INDEX IF EXISTS idx_schedules_skill;
DROP INDEX IF EXISTS idx_tasks_session_id;
DROP INDEX IF EXISTS idx_messages_session_id;
DROP INDEX IF EXISTS idx_sessions_user_id;

-- Drop tables
DROP TABLE IF EXISTS logs CASCADE;
DROP TABLE IF EXISTS schedules CASCADE;
DROP TABLE IF EXISTS skills CASCADE;
DROP TABLE IF EXISTS tasks CASCADE;
DROP TABLE IF EXISTS messages CASCADE;
DROP TABLE IF EXISTS sessions CASCADE;
DROP TABLE IF EXISTS users CASCADE;
