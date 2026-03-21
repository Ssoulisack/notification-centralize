-- 001_init.up.sql
-- Core notification log
CREATE TABLE IF NOT EXISTS notifications (
    id              VARCHAR(36) PRIMARY KEY,
    user_id         VARCHAR(255) NOT NULL,
    channel         VARCHAR(20)  NOT NULL,
    recipient       VARCHAR(500) NOT NULL,
    subject         VARCHAR(500) DEFAULT '',
    body            TEXT         DEFAULT '',
    template_id     VARCHAR(36)  DEFAULT '',
    priority        VARCHAR(20)  NOT NULL DEFAULT 'normal',
    status          VARCHAR(20)  NOT NULL DEFAULT 'pending',
    retry_count     INT          NOT NULL DEFAULT 0,
    metadata        JSONB        DEFAULT '{}',
    error_message   TEXT         DEFAULT '',
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    sent_at         TIMESTAMPTZ
);

CREATE INDEX idx_notifications_user_id    ON notifications (user_id);
CREATE INDEX idx_notifications_status     ON notifications (status);
CREATE INDEX idx_notifications_created_at ON notifications (created_at DESC);
CREATE INDEX idx_notifications_channel    ON notifications (channel);

-- Device tokens for push notifications
CREATE TABLE IF NOT EXISTS device_tokens (
    id          VARCHAR(36)  PRIMARY KEY,
    user_id     VARCHAR(255) NOT NULL,
    token       VARCHAR(500) NOT NULL UNIQUE,
    platform    VARCHAR(20)  NOT NULL,
    app_version VARCHAR(50)  DEFAULT '',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    last_used_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_device_tokens_user_id ON device_tokens (user_id);

-- User notification preferences
CREATE TABLE IF NOT EXISTS user_preferences (
    user_id          VARCHAR(255) PRIMARY KEY,
    enabled_channels TEXT[]       NOT NULL DEFAULT ARRAY['email']::TEXT[],
    quiet_start      INT          NOT NULL DEFAULT 22,
    quiet_end        INT          NOT NULL DEFAULT 8,
    opted_out_events TEXT[]       NOT NULL DEFAULT ARRAY[]::TEXT[],
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Notification templates
CREATE TABLE IF NOT EXISTS notification_templates (
    id               VARCHAR(36)  PRIMARY KEY,
    name             VARCHAR(100) NOT NULL,
    channel          VARCHAR(20)  NOT NULL,
    subject_template TEXT         DEFAULT '',
    body_template    TEXT         NOT NULL,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE(name, channel)
);

-- Seed a sample template
INSERT INTO notification_templates (id, name, channel, subject_template, body_template)
VALUES (
    'tmpl_order_shipped_email',
    'order_shipped',
    'email',
    'Your order {{.order_id}} has been shipped!',
    '<h1>Order Shipped</h1><p>Hi {{.name}},</p><p>Your order <strong>{{.order_id}}</strong> is on its way. Track it here: {{.tracking_url}}</p>'
) ON CONFLICT DO NOTHING;
