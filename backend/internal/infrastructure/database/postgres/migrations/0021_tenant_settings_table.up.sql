CREATE TABLE IF NOT EXISTS tenant_settings (
    tenant_id UUID PRIMARY KEY REFERENCES tenants(id) ON DELETE CASCADE,

    theme TEXT NOT NULL DEFAULT 'light'
        CONSTRAINT tenant_settings_theme_check CHECK (theme IN ('light', 'dark', 'system')),

    watermark_enabled BOOLEAN NOT NULL DEFAULT True,
    watermark_text TEXT DEFAULT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Trigger to keep updated_at fresh
CREATE TRIGGER trg_tenant_settings_updated_at
BEFORE UPDATE ON tenant_settings
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();  

-- Index for querying by theme
CREATE INDEX idx_tenant_settings_theme
    ON tenant_settings(theme);

-- Index for querying by watermark_enabled
CREATE INDEX idx_tenant_settings_watermark_enabled
    ON tenant_settings(watermark_enabled);

