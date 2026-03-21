-- 002_multi_project.down.sql
-- Rollback multi-project support

-- Remove project_id from notification_templates (restore original unique constraint)
ALTER TABLE notification_templates DROP CONSTRAINT notification_templates_project_name_channel_key;
ALTER TABLE notification_templates ADD CONSTRAINT notification_templates_name_channel_key UNIQUE (name, channel);
ALTER TABLE notification_templates DROP COLUMN project_id;

-- Remove project_id from user_preferences (restore original primary key)
ALTER TABLE user_preferences DROP CONSTRAINT user_preferences_pkey;
ALTER TABLE user_preferences ADD PRIMARY KEY (user_id);
ALTER TABLE user_preferences DROP COLUMN project_id;

-- Remove project_id from device_tokens (restore original unique constraint)
ALTER TABLE device_tokens DROP CONSTRAINT device_tokens_project_token_key;
ALTER TABLE device_tokens ADD CONSTRAINT device_tokens_token_key UNIQUE (token);
ALTER TABLE device_tokens DROP COLUMN project_id;

-- Remove project_id from notifications
DROP INDEX IF EXISTS idx_notifications_project_id;
ALTER TABLE notifications DROP COLUMN project_id;

-- Drop projects table
DROP INDEX IF EXISTS idx_projects_api_key;
DROP TABLE IF EXISTS projects;
