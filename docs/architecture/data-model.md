
## /docs/architecture/data-model.md

```markdown
# Data Model

This document describes the database schema, storage architecture, and data access patterns for Photo Listing SaaS.

## 🎯 Data Design Principles

### Core Principles
1. **Multi-Tenancy First**: All tenant-owned data includes tenant_id
2. **Row-Level Security**: Database-enforced tenant isolation
3. **Soft Deletes**: Preserve data for audit and recovery
4. **Audit Trail**: Track all data changes
5. **Performance**: Optimized indexes for common queries
6. **Scalability**: Design for horizontal scaling

### Data Isolation Strategy
- **Logical Isolation**: Row-Level Security with tenant context
- **Physical Isolation**: Separate storage paths per tenant
- **Cache Isolation**: Tenant-prefixed cache keys
- **Event Isolation**: Tenant context in all events

## 🏗️ Database Schema

### Core Tables

#### Tenants Table
```sql
CREATE TABLE tenants (
    -- Primary identifier
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Business information
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    
    -- Subscription information
    plan VARCHAR(50) NOT NULL DEFAULT 'free',
    subscription_id VARCHAR(255), -- External subscription ID
    subscription_status VARCHAR(50) NOT NULL DEFAULT 'active',
    trial_ends_at TIMESTAMP,
    
    -- Limits and usage
    limits JSONB NOT NULL DEFAULT '{
        "storage_gb": 5,
        "albums": 10,
        "users": 1,
        "bandwidth_gb": 10,
        "custom_domains": 0
    }',
    
    usage JSONB NOT NULL DEFAULT '{
        "storage_bytes": 0,
        "album_count": 0,
        "user_count": 1,
        "bandwidth_bytes": 0
    }',
    
    -- Settings and configuration
    settings JSONB NOT NULL DEFAULT '{
        "branding": {},
        "features": {},
        "notifications": {}
    }',
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    
    -- Indexes
    INDEX idx_tenants_slug ON tenants(slug),
    INDEX idx_tenants_plan ON tenants(plan),
    INDEX idx_tenants_created_at ON tenants(created_at DESC)
);

-- Row-Level Security for system tables (admin access only)
ALTER TABLE tenants ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenants_isolation ON tenants
    USING (false) WITH CHECK (false); -- No direct access

    Users Table
sql

CREATE TABLE users (
    -- Primary identifier
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    -- Authentication (linked to Supabase Auth)
    auth_user_id UUID UNIQUE, -- References auth.users.id
    email VARCHAR(255) NOT NULL,
    email_verified BOOLEAN NOT NULL DEFAULT false,
    
    -- Profile information
    name VARCHAR(255),
    avatar_url VARCHAR(500),
    bio TEXT,
    
    -- Role and permissions
    role VARCHAR(50) NOT NULL DEFAULT 'viewer',
    permissions JSONB NOT NULL DEFAULT '[]',
    
    -- Preferences
    preferences JSONB NOT NULL DEFAULT '{
        "notifications": true,
        "theme": "light",
        "language": "en"
    }',
    
    -- Activity tracking
    last_login_at TIMESTAMP,
    last_active_at TIMESTAMP,
    login_count INTEGER NOT NULL DEFAULT 0,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    
    -- Constraints
    UNIQUE(tenant_id, email),
    
    -- Indexes
    INDEX idx_users_tenant_id ON users(tenant_id),
    INDEX idx_users_email ON users(email),
    INDEX idx_users_role ON users(role),
    INDEX idx_users_last_active_at ON users(last_active_at DESC)
);

-- Row-Level Security for tenant isolation
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
CREATE POLICY users_tenant_isolation ON users
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid);

Albums Table
sql

CREATE TABLE albums (
    -- Primary identifier
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    -- Album information
    title VARCHAR(255) NOT NULL,
    description TEXT,
    slug VARCHAR(100) NOT NULL, -- URL-friendly identifier
    
    -- Status and visibility
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    visibility VARCHAR(50) NOT NULL DEFAULT 'private',
    password_hash VARCHAR(255), -- For password-protected albums
    
    -- Cover photo
    cover_photo_id UUID, -- References photos.id
    
    -- Ordering and organization
    sort_order INTEGER NOT NULL DEFAULT 0,
    category VARCHAR(100),
    tags VARCHAR(255)[] DEFAULT '{}',
    
    -- Settings
    settings JSONB NOT NULL DEFAULT '{
        "watermark": {
            "enabled": false,
            "text": "",
            "opacity": 0.3,
            "position": "bottom-right"
        },
        "download": {
            "enabled": true,
            "require_email": false,
            "max_downloads": null
        },
        "proofing": {
            "enabled": false,
            "max_selections": 25,
            "allow_comments": true
        },
        "sharing": {
            "allow_sharing": true,
            "embed_enabled": false
        }
    }',
    
    -- Statistics (denormalized for performance)
    stats JSONB NOT NULL DEFAULT '{
        "photo_count": 0,
        "total_size_bytes": 0,
        "view_count": 0,
        "download_count": 0,
        "share_count": 0,
        "avg_view_time_seconds": 0
    }',
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    published_at TIMESTAMP,
    archived_at TIMESTAMP,
    deleted_at TIMESTAMP,
    
    -- Optimistic concurrency control
    version INTEGER NOT NULL DEFAULT 1,
    
    -- Constraints
    UNIQUE(tenant_id, slug),
    
    -- Indexes
    INDEX idx_albums_tenant_id ON albums(tenant_id),
    INDEX idx_albums_status ON albums(status),
    INDEX idx_albums_visibility ON albums(visibility),
    INDEX idx_albums_created_at ON albums(created_at DESC),
    INDEX idx_albums_published_at ON albums(published_at DESC),
    
    -- Partial indexes for common queries
    INDEX idx_albums_tenant_status ON albums(tenant_id, status)
        WHERE status IN ('draft', 'published'),
    INDEX idx_albums_public ON albums(tenant_id, slug)
        WHERE status = 'published' AND visibility = 'public'
);

-- Row-Level Security for tenant isolation
ALTER TABLE albums ENABLE ROW LEVEL SECURITY;
CREATE POLICY albums_tenant_isolation ON albums
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid);

-- Public read policy for published albums
CREATE POLICY albums_public_read ON albums
    FOR SELECT
    USING (
        status = 'published' 
        AND visibility = 'public'
    );

Photos Table
sql

CREATE TABLE photos (
    -- Primary identifier
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    album_id UUID REFERENCES albums(id) ON DELETE SET NULL,
    
    -- File information
    original_filename VARCHAR(500) NOT NULL,
    storage_path VARCHAR(500) NOT NULL, -- Path in object storage
    file_hash VARCHAR(64) NOT NULL, -- SHA-256 for deduplication
    
    -- File metadata
    mime_type VARCHAR(100),
    file_size BIGINT NOT NULL,
    width INTEGER,
    height INTEGER,
    aspect_ratio DECIMAL(5, 3),
    color_profile VARCHAR(50),
    format VARCHAR(20),
    
    -- Processing information
    processing_status VARCHAR(50) NOT NULL DEFAULT 'pending',
    processing_errors JSONB DEFAULT '[]',
    processing_started_at TIMESTAMP,
    processing_completed_at TIMESTAMP,
    
    -- EXIF metadata (extracted from image)
    exif JSONB DEFAULT '{}',
    
    -- GPS data (extracted from EXIF)
    location GEOGRAPHY(POINT, 4326),
    location_name VARCHAR(255),
    
    -- Custom metadata
    metadata JSONB NOT NULL DEFAULT '{
        "title": null,
        "description": null,
        "tags": [],
        "rating": 0,
        "color_palette": [],
        "faces": []
    }',
    
    -- Version information (multiple sizes/formats)
    versions JSONB NOT NULL DEFAULT '{
        "original": {
            "path": "",
            "size_bytes": 0,
            "width": 0,
            "height": 0
        },
        "web": {
            "path": "",
            "size_bytes": 0,
            "width": 1920,
            "height": 0,
            "format": "webp"
        },
        "thumbnails": {
            "small": {"path": "", "width": 320, "height": 0},
            "medium": {"path": "", "width": 640, "height": 0},
            "large": {"path": "", "width": 1200, "height": 0}
        }
    }',
    
    -- Statistics
    stats JSONB NOT NULL DEFAULT '{
        "view_count": 0,
        "download_count": 0,
        "share_count": 0,
        "selection_count": 0
    }',
    
    -- Ordering
    sort_order INTEGER NOT NULL DEFAULT 0,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    uploaded_at TIMESTAMP NOT NULL DEFAULT NOW(),
    taken_at TIMESTAMP, -- From EXIF DateTimeOriginal
    deleted_at TIMESTAMP,
    
    -- Optimistic concurrency control
    version INTEGER NOT NULL DEFAULT 1,
    
    -- Indexes
    INDEX idx_photos_tenant_id ON photos(tenant_id),
    INDEX idx_photos_album_id ON photos(album_id),
    INDEX idx_photos_created_at ON photos(created_at DESC),
    INDEX idx_photos_taken_at ON photos(taken_at DESC),
    INDEX idx_photos_processing_status ON photos(processing_status),
    INDEX idx_photos_file_hash ON photos(file_hash),
    
    -- GIN indexes for JSON fields
    INDEX idx_photos_metadata_tags ON photos USING GIN ((metadata->'tags')),
    INDEX idx_photos_exif ON photos USING GIN (exif),
    
    -- Partial indexes
    INDEX idx_photos_active ON photos(tenant_id, album_id, sort_order)
        WHERE deleted_at IS NULL,
    
    -- Spatial index for location queries
    INDEX idx_photos_location ON photos USING GIST (location)
);

-- Row-Level Security
ALTER TABLE photos ENABLE ROW LEVEL SECURITY;
CREATE POLICY photos_tenant_isolation ON photos
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid);

Shares Table (Client Access)
sql

CREATE TABLE shares (
    -- Primary identifier
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    album_id UUID NOT NULL REFERENCES albums(id) ON DELETE CASCADE,
    
    -- Share token (public identifier)
    token VARCHAR(64) UNIQUE NOT NULL,
    
    -- Share configuration
    type VARCHAR(50) NOT NULL DEFAULT 'view', -- view, download, proofing
    name VARCHAR(255),
    password_hash VARCHAR(255),
    
    -- Access controls
    permissions JSONB NOT NULL DEFAULT '{
        "download": false,
        "share": false,
        "comment": false,
        "select": false,
        "watermark": true
    }',
    
    -- Watermark settings
    watermark JSONB DEFAULT '{
        "text": "Proof - Do Not Share",
        "opacity": 0.5,
        "position": "bottom-right"
    }',
    
    -- Expiration
    expires_at TIMESTAMP,
    max_views INTEGER,
    max_downloads INTEGER,
    
    -- Client information (for proofing)
    client_email VARCHAR(255),
    client_name VARCHAR(255),
    client_message TEXT,
    
    -- Statistics
    stats JSONB NOT NULL DEFAULT '{
        "view_count": 0,
        "download_count": 0,
        "selection_count": 0,
        "comment_count": 0
    }',
    
    -- Notification settings
    notifications JSONB NOT NULL DEFAULT '{
        "on_view": false,
        "on_download": false,
        "on_selection": false,
        "on_comment": false
    }',
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    accessed_at TIMESTAMP,
    deleted_at TIMESTAMP,
    
    -- Indexes
    INDEX idx_shares_token ON shares(token),
    INDEX idx_shares_tenant_id ON shares(tenant_id),
    INDEX idx_shares_album_id ON shares(album_id),
    INDEX idx_shares_expires_at ON shares(expires_at),
    INDEX idx_shares_created_at ON shares(created_at DESC)
);

-- Row-Level Security
ALTER TABLE shares ENABLE ROW LEVEL SECURITY;
CREATE POLICY shares_tenant_isolation ON shares
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid);

-- Special policy for token-based access (no authentication required)
CREATE POLICY shares_token_access ON shares
    FOR SELECT
    USING (
        token = current_setting('app.share_token', true)
        AND (expires_at IS NULL OR expires_at > NOW())
        AND (max_views IS NULL OR stats->>'view_count'::integer < max_views)
        AND deleted_at IS NULL
    );

Selections Table (Client Proofing)
sql

CREATE TABLE selections (
    -- Primary identifier
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    share_id UUID NOT NULL REFERENCES shares(id) ON DELETE CASCADE,
    photo_id UUID NOT NULL REFERENCES photos(id) ON DELETE CASCADE,
    
    -- Selection data
    rating INTEGER CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    tags VARCHAR(50)[] DEFAULT '{}',
    
    -- Metadata
    metadata JSONB NOT NULL DEFAULT '{
        "coordinates": null,
        "color_notes": null,
        "crop_suggestions": null
    }',
    
    -- Status
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    submitted_at TIMESTAMP,
    reviewed_at TIMESTAMP,
    
    -- Indexes
    INDEX idx_selections_tenant_id ON selections(tenant_id),
    INDEX idx_selections_share_id ON selections(share_id),
    INDEX idx_selections_photo_id ON selections(photo_id),
    INDEX idx_selections_status ON selections(status),
    INDEX idx_selections_created_at ON selections(created_at DESC),
    
    -- Unique constraint
    UNIQUE(share_id, photo_id)
);

-- Row-Level Security
ALTER TABLE selections ENABLE ROW LEVEL SECURITY;
CREATE POLICY selections_tenant_isolation ON selections
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid);

Events Table (Audit Trail)
sql

CREATE TABLE events (
    -- Primary identifier
    id BIGSERIAL PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    -- Event information
    event_id UUID NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    aggregate_type VARCHAR(100) NOT NULL,
    aggregate_id UUID NOT NULL,
    
    -- Event data
    data JSONB NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}',
    
    -- Version and ordering
    version INTEGER NOT NULL,
    sequence_number BIGINT NOT NULL,
    
    -- User context
    user_id UUID,
    user_ip INET,
    user_agent TEXT,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Indexes
    INDEX idx_events_tenant_id ON events(tenant_id),
    INDEX idx_events_event_type ON events(event_type),
    INDEX idx_events_aggregate ON events(aggregate_type, aggregate_id),
    INDEX idx_events_created_at ON events(created_at DESC),
    INDEX idx_events_sequence ON events(tenant_id, aggregate_type, sequence_number),
    
    -- Unique constraint for idempotency
    UNIQUE(tenant_id, aggregate_type, aggregate_id, version)
);

-- Row-Level Security
ALTER TABLE events ENABLE ROW LEVEL SECURITY;
CREATE POLICY events_tenant_isolation ON events
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid);

Event Outbox Table (Reliable Messaging)
sql

CREATE TABLE event_outbox (
    -- Primary identifier
    id BIGSERIAL PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    -- Event information
    event_id UUID NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    aggregate_type VARCHAR(100) NOT NULL,
    aggregate_id UUID NOT NULL,
    
    -- Event payload
    payload JSONB NOT NULL,
    
    -- Processing status
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    retry_count INTEGER NOT NULL DEFAULT 0,
    error_message TEXT,
    last_error_at TIMESTAMP,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMP,
    
    -- Indexes for efficient processing
    INDEX idx_outbox_pending ON event_outbox(tenant_id, status, created_at)
        WHERE status = 'pending',
    INDEX idx_outbox_retry ON event_outbox(tenant_id, status, retry_count, last_error_at)
        WHERE status = 'failed' AND retry_count < 3,
    
    -- Partitioning by tenant for large scale
    UNIQUE(tenant_id, event_id)
) PARTITION BY LIST (tenant_id);

-- Create initial partition (additional partitions created dynamically)
CREATE TABLE event_outbox_default PARTITION OF event_outbox
    DEFAULT;

🔄 Relationships & Constraints
Foreign Key Relationships
sql

-- Album-Photo relationship (one-to-many)
ALTER TABLE photos 
ADD CONSTRAINT fk_photos_album 
FOREIGN KEY (album_id) 
REFERENCES albums(id) 
ON DELETE SET NULL;

-- Album cover photo relationship
ALTER TABLE albums 
ADD CONSTRAINT fk_albums_cover_photo 
FOREIGN KEY (cover_photo_id) 
REFERENCES photos(id) 
ON DELETE SET NULL;

-- Share-Album relationship
ALTER TABLE shares 
ADD CONSTRAINT fk_shares_album 
FOREIGN KEY (album_id) 
REFERENCES albums(id) 
ON DELETE CASCADE;

-- Selection-Share relationship
ALTER TABLE selections 
ADD CONSTRAINT fk_selections_share 
FOREIGN KEY (share_id) 
REFERENCES shares(id) 
ON DELETE CASCADE;

-- Selection-Photo relationship
ALTER TABLE selections 
ADD CONSTRAINT fk_selections_photo 
FOREIGN KEY (photo_id) 
REFERENCES photos(id) 
ON DELETE CASCADE;

Check Constraints
sql

-- Album status validation
ALTER TABLE albums 
ADD CONSTRAINT check_album_status 
CHECK (status IN ('draft', 'review', 'published', 'archived'));

-- Album visibility validation
ALTER TABLE albums 
ADD CONSTRAINT check_album_visibility 
CHECK (visibility IN ('public', 'private', 'password', 'client'));

-- Photo processing status validation
ALTER TABLE photos 
ADD CONSTRAINT check_photo_processing_status 
CHECK (processing_status IN ('pending', 'processing', 'completed', 'failed'));

-- Share type validation
ALTER TABLE shares 
ADD CONSTRAINT check_share_type 
CHECK (type IN ('view', 'download', 'proofing'));

-- Selection rating validation
ALTER TABLE selections 
ADD CONSTRAINT check_selection_rating 
CHECK (rating IS NULL OR (rating >= 1 AND rating <= 5));

🏪 Storage Architecture
Cloudflare R2 Structure
text

r2://photo-listing/
├── tenants/
│   ├── {tenant-id}/
│   │   ├── originals/          # Original uploaded files
│   │   │   ├── {year}/{month}/{day}/
│   │   │   │   ├── {uuid}.{ext}
│   │   │   │   └── ...
│   │   │   └── ...
│   │   ├── web/                # Web-optimized versions
│   │   │   ├── {uuid}.webp
│   │   │   └── ...
│   │   ├── thumbnails/         # Multiple thumbnail sizes
│   │   │   ├── {uuid}-sm.webp  # 320px width
│   │   │   ├── {uuid}-md.webp  # 640px width
│   │   │   ├── {uuid}-lg.webp  # 1200px width
│   │   │   └── ...
│   │   ├── watermarks/         # Watermarked versions
│   │   │   ├── {share-token}/
│   │   │   │   ├── {uuid}.webp
│   │   │   │   └── ...
│   │   │   └── ...
│   │   └── exports/            # Downloaded/exported content
│   │       ├── {album-id}/
│   │       │   ├── full-resolution.zip
│   │       │   ├── web-resolution.zip
│   │       │   └── ...
│   │       └── ...
│   └── ...
├── system/
│   ├── avatars/               # User profile pictures
│   │   ├── {user-id}.webp
│   │   └── ...
│   ├── templates/             # Email/PDF templates
│   │   ├── emails/
│   │   │   ├── invitation.html
│   │   │   └── ...
│   │   └── watermarks/
│   │       ├── default.png
│   │       └── ...
│   └── backups/               # System backups
│       └── ...
└── cdn/
    ├── cache/                 # CDN cached content
    └── ...

Storage Path Generation
go

// Storage path generator
type StoragePathGenerator struct {
    bucketName string
}

func (g *StoragePathGenerator) PhotoPath(tenantID, photoID, version, extension string) string {
    switch version {
    case "original":
        // Organized by date for easier management
        date := time.Now().Format("2006/01/02")
        return fmt.Sprintf("tenants/%s/originals/%s/%s.%s",
            tenantID, date, photoID, extension)
    case "web":
        return fmt.Sprintf("tenants/%s/web/%s.webp", tenantID, photoID)
    case "thumbnail-sm":
        return fmt.Sprintf("tenants/%s/thumbnails/%s-sm.webp", tenantID, photoID)
    case "thumbnail-md":
        return fmt.Sprintf("tenants/%s/thumbnails/%s-md.webp", tenantID, photoID)
    case "thumbnail-lg":
        return fmt.Sprintf("tenants/%s/thumbnails/%s-lg.webp", tenantID, photoID)
    default:
        return fmt.Sprintf("tenants/%s/%s/%s", tenantID, version, photoID)
    }
}

func (g *StoragePathGenerator) WatermarkedPath(tenantID, shareToken, photoID string) string {
    return fmt.Sprintf("tenants/%s/watermarks/%s/%s.webp",
        tenantID, shareToken, photoID)
}

func (g *StoragePathGenerator) ExportPath(tenantID, albumID, filename string) string {
    return fmt.Sprintf("tenants/%s/exports/%s/%s",
        tenantID, albumID, filename)
}

📊 Indexing Strategy
Primary Indexes
sql

-- B-tree indexes for equality and range queries
CREATE INDEX idx_albums_tenant_created ON albums(tenant_id, created_at DESC);
CREATE INDEX idx_photos_album_order ON photos(album_id, sort_order ASC);
CREATE INDEX idx_shares_token_expires ON shares(token, expires_at);

-- Hash indexes for exact match queries
CREATE INDEX idx_users_email_hash ON users USING hash(email);

-- GIN indexes for JSONB and array queries
CREATE INDEX idx_photos_metadata_gin ON photos USING GIN (metadata);
CREATE INDEX idx_albums_tags_gin ON albums USING GIN (tags);
CREATE INDEX idx_selections_tags_gin ON selections USING GIN (tags);

-- GiST indexes for spatial queries
CREATE INDEX idx_photos_location_gist ON photos USING GIST (location);

-- BRIN indexes for large timestamp ranges
CREATE INDEX idx_events_created_brin ON events USING BRIN (created_at);

Partial Indexes
sql

-- Active albums only
CREATE INDEX idx_albums_active ON albums(tenant_id, status, visibility)
    WHERE status = 'published' AND deleted_at IS NULL;

-- Pending processing photos
CREATE INDEX idx_photos_pending ON photos(processing_status, created_at)
    WHERE processing_status = 'pending';

-- Active shares (not expired)
CREATE INDEX idx_shares_active ON shares(token, expires_at)
    WHERE expires_at IS NULL OR expires_at > NOW();

-- Recent events (last 30 days)
CREATE INDEX idx_events_recent ON events(tenant_id, created_at DESC)
    WHERE created_at > NOW() - INTERVAL '30 days';

Covering Indexes
sql

-- Covering index for album listings
CREATE INDEX idx_albums_listing ON albums(
    tenant_id, 
    status, 
    created_at DESC
) INCLUDE (
    title, 
    description, 
    cover_photo_id, 
    stats
) WHERE deleted_at IS NULL;

-- Covering index for photo listings
CREATE INDEX idx_photos_listing ON photos(
    album_id, 
    sort_order ASC
) INCLUDE (
    id,
    original_filename,
    width,
    height,
    versions,
    metadata
) WHERE deleted_at IS NULL AND processing_status = 'completed';

🔄 Data Access Patterns
Common Query Patterns
Album Listing with Pagination
sql

-- Get paginated album list for tenant
SELECT 
    a.id,
    a.title,
    a.description,
    a.status,
    a.visibility,
    a.cover_photo_id,
    a.stats->>'photo_count' as photo_count,
    a.stats->>'view_count' as view_count,
    a.created_at,
    a.published_at,
    -- Use COUNT(*) OVER() for total count without extra query
    COUNT(*) OVER() as total_count
FROM albums a
WHERE a.tenant_id = $1
AND a.deleted_at IS NULL
AND ($2::text IS NULL OR a.status = $2)
AND ($3::text IS NULL OR a.visibility = $3)
ORDER BY 
    CASE 
        WHEN $4 = 'created' THEN a.created_at
        WHEN $4 = 'published' THEN a.published_at
        ELSE a.created_at
    END DESC
LIMIT $5 OFFSET $6;

Photo Listing with Efficient Join
sql

-- Get photos with album info in single query
SELECT 
    p.id,
    p.original_filename,
    p.width,
    p.height,
    p.versions,
    p.metadata,
    p.created_at,
    p.sort_order,
    a.title as album_title,
    a.slug as album_slug
FROM photos p
INNER JOIN albums a ON p.album_id = a.id
WHERE p.tenant_id = $1
AND p.album_id = $2
AND p.deleted_at IS NULL
AND p.processing_status = 'completed'
ORDER BY p.sort_order ASC, p.created_at ASC
LIMIT $3 OFFSET $4;

Aggregated Statistics
sql

-- Get tenant statistics with single query
SELECT 
    -- Album counts by status
    COUNT(*) FILTER (WHERE status = 'draft') as draft_count,
    COUNT(*) FILTER (WHERE status = 'published') as published_count,
    COUNT(*) FILTER (WHERE status = 'archived') as archived_count,
    
    -- Photo statistics
    COALESCE(SUM(file_size), 0) as total_storage_bytes,
    COUNT(DISTINCT p.id) as total_photos,
    
    -- User statistics
    COUNT(DISTINCT u.id) as total_users,
    
    -- Share statistics
    COUNT(DISTINCT s.id) as active_shares,
    
    -- Recent activity
    MAX(a.created_at) as last_album_created,
    MAX(p.created_at) as last_photo_uploaded
FROM tenants t
LEFT JOIN albums a ON t.id = a.tenant_id AND a.deleted_at IS NULL
LEFT JOIN photos p ON a.id = p.album_id AND p.deleted_at IS NULL
LEFT JOIN users u ON t.id = u.tenant_id AND u.deleted_at IS NULL
LEFT JOIN shares s ON a.id = s.album_id AND s.deleted_at IS NULL
WHERE t.id = $1
GROUP BY t.id;

Materialized Views for Analytics
sql

-- Daily statistics materialized view
CREATE MATERIALIZED VIEW mv_daily_stats AS
SELECT 
    DATE(created_at) as date,
    tenant_id,
    COUNT(DISTINCT id) as album_count,
    COUNT(DISTINCT CASE WHEN status = 'published' THEN id END) as published_count,
    SUM((stats->>'view_count')::integer) as total_views,
    SUM((stats->>'download_count')::integer) as total_downloads
FROM albums
WHERE deleted_at IS NULL
GROUP BY DATE(created_at), tenant_id
WITH DATA;

-- Refresh on schedule
CREATE OR REPLACE FUNCTION refresh_daily_stats()
RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY mv_daily_stats;
END;
$$ LANGUAGE plpgsql;

-- Schedule refresh (using pg_cron or similar)
SELECT cron.schedule('refresh-daily-stats', '0 2 * * *', 
    'SELECT refresh_daily_stats()');

Full-Text Search
sql

-- Search albums by title and description
CREATE INDEX idx_albums_search ON albums 
USING gin(to_tsvector('english', title || ' ' || COALESCE(description, '')));

-- Search photos by filename and metadata
CREATE INDEX idx_photos_search ON photos 
USING gin(
    to_tsvector('english', 
        original_filename || ' ' || 
        COALESCE(metadata->>'title', '') || ' ' ||
        COALESCE(metadata->>'description', '') || ' ' ||
        array_to_string(metadata->'tags', ' ')
    )
);

-- Search query
SELECT 
    a.id,
    a.title,
    a.description,
    ts_rank_cd(
        to_tsvector('english', a.title || ' ' || COALESCE(a.description, '')),
        plainto_tsquery('english', $1)
    ) as rank
FROM albums a
WHERE a.tenant_id = $2
AND a.deleted_at IS NULL
AND to_tsvector('english', a.title || ' ' || COALESCE(a.description, '')) 
    @@ plainto_tsquery('english', $1)
ORDER BY rank DESC
LIMIT 20;

🗄️ Data Migration Strategy
Versioned Migrations
sql

-- migrations/001_initial_schema.up.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create tables as shown above...

-- migrations/002_add_fulltext_search.up.sql
CREATE INDEX idx_albums_search ON albums 
USING gin(to_tsvector('english', title || ' ' || COALESCE(description, '')));

-- migrations/003_add_spatial_support.up.sql
CREATE EXTENSION IF NOT EXISTS "postgis";
ALTER TABLE photos ADD COLUMN location GEOGRAPHY(POINT, 4326);
CREATE INDEX idx_photos_location ON photos USING GIST (location);

Data Transformation Migrations
sql

-- migrations/004_migrate_album_settings.up.sql
-- Convert old settings format to new format
UPDATE albums 
SET settings = jsonb_set(
    COALESCE(settings, '{}'::jsonb),
    '{watermark}',
    jsonb_build_object(
        'enabled', COALESCE(settings->>'watermark_enabled', 'false')::boolean,
        'text', COALESCE(settings->>'watermark_text', ''),
        'opacity', COALESCE((settings->>'watermark_opacity')::float, 0.3)
    )
)
WHERE settings IS NOT NULL;

Zero-Downtime Migrations
sql

-- Step 1: Add new column without constraint
ALTER TABLE photos ADD COLUMN file_hash VARCHAR(64);

-- Step 2: Backfill data (in batches)
UPDATE photos 
SET file_hash = encode(
    sha256(storage_path::bytea), 
    'hex'
)
WHERE file_hash IS NULL
AND id IN (
    SELECT id FROM photos 
    WHERE file_hash IS NULL 
    ORDER BY created_at 
    LIMIT 1000
);

-- Step 3: Add constraint after backfill complete
ALTER TABLE photos 
ADD CONSTRAINT photos_file_hash_not_null 
CHECK (file_hash IS NOT NULL) NOT VALID;

-- Step 4: Validate constraint (doesn't lock table)
ALTER TABLE photos VALIDATE CONSTRAINT photos_file_hash_not_null;

-- Step 5: Create index concurrently
CREATE INDEX CONCURRENTLY idx_photos_file_hash 
ON photos(file_hash);

🔒 Data Security
Row-Level Security Policies
sql

-- Enable RLS on all tenant tables
SELECT tablename 
FROM pg_tables 
WHERE schemaname = 'public'
AND tablename NOT LIKE 'pg_%'
AND tablename NOT IN ('tenants', 'schema_migrations')
ORDER BY tablename;

-- Function to set tenant context
CREATE OR REPLACE FUNCTION set_current_tenant(tenant_id UUID)
RETURNS void AS $$
BEGIN
    PERFORM set_config('app.current_tenant_id', tenant_id::text, true);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Function to get current tenant
CREATE OR REPLACE FUNCTION get_current_tenant_id()
RETURNS UUID AS $$
BEGIN
    RETURN current_setting('app.current_tenant_id')::UUID;
EXCEPTION
    WHEN OTHERS THEN
        RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- RLS policy function
CREATE OR REPLACE FUNCTION tenant_policy_check()
RETURNS boolean AS $$
BEGIN
    RETURN (
        current_setting('app.current_tenant_id') IS NOT NULL
        AND current_setting('app.current_tenant_id') != ''
        AND tenant_id = current_setting('app.current_tenant_id')::UUID
    );
EXCEPTION
    WHEN OTHERS THEN
        RETURN false;
END;
$$ LANGUAGE plpgsql STABLE;

-- Apply RLS policies
CREATE POLICY tenant_isolation ON albums
    FOR ALL
    USING (tenant_policy_check())
    WITH CHECK (tenant_policy_check());

Data Encryption
sql

-- Encrypt sensitive columns
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Function to encrypt data with tenant-specific key
CREATE OR REPLACE FUNCTION encrypt_data(
    data TEXT,
    tenant_id UUID
)
RETURNS BYTEA AS $$
DECLARE
    encryption_key TEXT;
BEGIN
    -- Get tenant-specific encryption key
    SELECT encryption_key INTO encryption_key
    FROM tenant_encryption_keys
    WHERE tenant_id = $2
    AND active = true;
    
    IF encryption_key IS NULL THEN
        RAISE EXCEPTION 'No encryption key found for tenant %', tenant_id;
    END IF;
    
    -- Encrypt data
    RETURN pgp_sym_encrypt(data, encryption_key);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Function to decrypt data
CREATE OR REPLACE FUNCTION decrypt_data(
    encrypted_data BYTEA,
    tenant_id UUID
)
RETURNS TEXT AS $$
DECLARE
    encryption_key TEXT;
BEGIN
    -- Get tenant-specific encryption key
    SELECT encryption_key INTO encryption_key
    FROM tenant_encryption_keys
    WHERE tenant_id = $2
    AND active = true;
    
    IF encryption_key IS NULL THEN
        RAISE EXCEPTION 'No encryption key found for tenant %', tenant_id;
    END IF;
    
    -- Decrypt data
    RETURN pgp_sym_decrypt(encrypted_data, encryption_key);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

📈 Performance Optimization
Query Optimization Examples
sql

-- Optimized album list with cover photo
SELECT 
    a.*,
    p.versions->'thumbnails'->'medium'->>'path' as cover_photo_url
FROM albums a
LEFT JOIN LATERAL (
    SELECT versions
    FROM photos p
    WHERE p.id = a.cover_photo_id
    LIMIT 1
) p ON true
WHERE a.tenant_id = $1
AND a.deleted_at IS NULL
ORDER BY a.created_at DESC
LIMIT $2 OFFSET $3;

-- Batch statistics update
WITH album_stats AS (
    SELECT 
        album_id,
        COUNT(*) as photo_count,
        SUM(file_size) as total_size,
        MAX(created_at) as last_photo_added
    FROM photos
    WHERE deleted_at IS NULL
    AND processing_status = 'completed'
    GROUP BY album_id
)
UPDATE albums a
SET 
    stats = jsonb_set(
        COALESCE(stats, '{}'::jsonb),
        '{photo_count,total_size_bytes,last_photo_added}',
        jsonb_build_array(
            COALESCE(s.photo_count, 0),
            COALESCE(s.total_size, 0),
            COALESCE(s.last_photo_added::text, '')
        )
    ),
    updated_at = NOW()
FROM album_stats s
WHERE a.id = s.album_id
AND a.tenant_id = $1;

Connection Pooling Configuration
sql

-- PgBouncer configuration (pgbouncer.ini)
[databases]
photo_listing = host=localhost port=5432 dbname=photo_listing

[pgbouncer]
pool_mode = transaction
max_client_conn = 1000
default_pool_size = 20
min_pool_size = 5
reserve_pool_size = 5
reserve_pool_timeout = 5
query_timeout = 120
idle_transaction_timeout = 60
server_idle_timeout = 60

🔄 Backup & Recovery
Backup Strategy
sql

-- Logical backup script
#!/bin/bash
# backup.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backups/$DATE"

mkdir -p $BACKUP_DIR

# Backup schema
pg_dump --schema-only --no-owner --no-privileges \
    $DATABASE_URL > $BACKUP_DIR/schema.sql

# Backup data (tenant by tenant)
TENANTS=$(psql $DATABASE_URL -t -c "SELECT id FROM tenants")

for TENANT_ID in $TENANTS; do
    # Backup tenant data
    pg_dump \
        --data-only \
        --table="albums" \
        --where="tenant_id='$TENANT_ID'" \
        $DATABASE_URL > $BACKUP_DIR/tenant_${TENANT_ID}_albums.sql
    
    # Continue for other tables...
done

# Backup configuration and users
pg_dump --table="tenants" --table="users" \
    $DATABASE_URL > $BACKUP_DIR/system.sql

Point-in-Time Recovery
sql

-- WAL archiving configuration
-- In postgresql.conf:
wal_level = replica
archive_mode = on
archive_command = 'rclone copy %p r2:photo-listing-wal/%f'

-- Recovery configuration
-- In recovery.conf:
restore_command = 'rclone copy r2:photo-listing-wal/%f %p'
recovery_target_time = '2024-01-15 14:30:00'
recovery_target_action = 'pause'

This data model provides a scalable, secure foundation for the Photo Listing SaaS platform. Regular reviews and optimizations should be performed as the application grows and usage patterns evolve.