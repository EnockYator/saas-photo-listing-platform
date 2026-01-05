-- User roles (global)
CREATE TYPE user_role AS ENUM (
    'super_admin',
    'tenant_admin',
    'tenat_editor',
    'viewer'
);

-- Tenant scoped roles
CREATE TYPE tenant_user_role AS ENUM (
    'tenant_admin',
    'tenant_editor',
    'viewer'
);

-- -- Listing enums
-- CREATE TYPE listing_status AS ENUM ('draft', 'published');
-- CREATE TYPE listing_visibility AS ENUM ('private', 'public');

-- -- File watermark type
-- CREATE TYPE watermark_type_enum AS ENUM ('text', 'image');

-- -- Audit action types
-- CREATE TYPE audit_action AS ENUM (
--     'CREATE',
--     'UPDATE',
--     'DELETE',
--     'LOGIN',
--     'LOGOUT',
--     'PUBLISH',
--     'UNPUBLISH',
--     'SHARE',
--     'OTHER'
-- );

-- -- Subscription status
-- CREATE TYPE subscription_status AS ENUM ('active', 'inactive', 'canceled', 'past_due');

-- -- Share permissions
-- CREATE TYPE share_permission AS ENUM ('read', 'write', 'admin');
