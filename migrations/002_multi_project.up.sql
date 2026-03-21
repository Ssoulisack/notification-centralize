-- 002_multi_project.up.sql
-- Multi-project (multi-tenancy) support

-- Projects table
CREATE TABLE IF NOT EXISTS projects (
    id          VARCHAR(36) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    api_key     VARCHAR(255) NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_projects_api_key ON projects (api_key);

-- Default project for backwards compatibility
INSERT INTO projects (id, name, api_key)
VALUES ('00000000-0000-0000-0000-000000000000', 'default', 'default-api-key');

-- Add project_id to notifications
ALTER TABLE notifications ADD COLUMN project_id VARCHAR(36) REFERENCES projects(id);
UPDATE notifications SET project_id = '00000000-0000-0000-0000-000000000000';
ALTER TABLE notifications ALTER COLUMN project_id SET NOT NULL;
CREATE INDEX idx_notifications_project_id ON notifications (project_id);

-- Add project_id to device_tokens (update unique constraint)
ALTER TABLE device_tokens ADD COLUMN project_id VARCHAR(36) REFERENCES projects(id);
UPDATE device_tokens SET project_id = '00000000-0000-0000-0000-000000000000';
ALTER TABLE device_tokens ALTER COLUMN project_id SET NOT NULL;
ALTER TABLE device_tokens DROP CONSTRAINT device_tokens_token_key;
ALTER TABLE device_tokens ADD CONSTRAINT device_tokens_project_token_key UNIQUE (project_id, token);

-- Add project_id to user_preferences (change primary key)
ALTER TABLE user_preferences ADD COLUMN project_id VARCHAR(36) REFERENCES projects(id);
UPDATE user_preferences SET project_id = '00000000-0000-0000-0000-000000000000';
ALTER TABLE user_preferences ALTER COLUMN project_id SET NOT NULL;
ALTER TABLE user_preferences DROP CONSTRAINT user_preferences_pkey;
ALTER TABLE user_preferences ADD PRIMARY KEY (project_id, user_id);

-- Add project_id to notification_templates (update unique constraint)
ALTER TABLE notification_templates ADD COLUMN project_id VARCHAR(36) REFERENCES projects(id);
UPDATE notification_templates SET project_id = '00000000-0000-0000-0000-000000000000';
ALTER TABLE notification_templates ALTER COLUMN project_id SET NOT NULL;
ALTER TABLE notification_templates DROP CONSTRAINT notification_templates_name_channel_key;
ALTER TABLE notification_templates ADD CONSTRAINT notification_templates_project_name_channel_key UNIQUE (project_id, name, channel);
