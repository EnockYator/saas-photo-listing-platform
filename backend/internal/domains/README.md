
## /internal/domain/README.md

```markdown
# Domain Layer

The domain layer contains the core business logic and rules for Photo Listing SaaS. It follows Domain-Driven Design (DDD) principles and Clean Architecture.

## 🎯 Domain Overview

### Core Domains
1. **Tenant Management**: Multi-tenancy, subscription, limits
2. **Portfolio Management**: Albums, photos, organization
3. **Client Delivery**: Sharing, proofing, downloads
4. **User Management**: Roles, permissions, authentication
5. **Billing & Subscriptions**: Plans, payments, usage tracking

### Strategic Design

┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│ Portfolio │ │ Client │ │ Billing │
│ Management │◄──┤ Delivery │◄──┤ & Payments │
│ Context │ │ Context │ │ Context │
└─────────────────┘ └─────────────────┘ └─────────────────┘
│ │ │
└───────────────────────┼───────────────────────┘
│
┌─────────────────┐
│ User │
│ Management │
│ Context │
└─────────────────┘
text


## 🏗️ Domain Structure

/internal/domain/
├── tenant/ # Tenant aggregate
│ ├── tenant.go # Tenant entity
│ ├── subscription.go # Subscription value object
│ ├── limits.go # Resource limits value object
│ ├── events.go # Tenant domain events
│ └── repository.go # Tenant repository interface
│
├── album/ # Album aggregate
│ ├── album.go # Album entity
│ ├── photo.go # Photo entity
│ ├── visibility.go # Visibility value object
│ ├── events.go # Album domain events
│ └── repository.go # Album repository interface
│
├── user/ # User aggregate
│ ├── user.go # User entity
│ ├── role.go # Role value object
│ ├── permissions.go # Permissions value object
│ ├── events.go # User domain events
│ └── repository.go # User repository interface
│
├── client/ # Client proofing
│ ├── share.go # Share entity
│ ├── selection.go # Selection entity
│ ├── watermark.go # Watermark value object
│ ├── events.go # Client domain events
│ └── repository.go # Client repository interface
│
├── billing/ # Billing domain
│ ├── plan.go # Plan value object
│ ├── invoice.go # Invoice entity
│ ├── payment.go # Payment entity
│ ├── events.go # Billing domain events
│ └── repository.go # Billing repository interface
│
├── shared/ # Shared domain concepts
│ ├── valueobject/ # Shared value objects
│ │ ├── id.go # Identifier value object
│ │ ├── metadata.go # Metadata value object
│ │ ├── file.go # File information value object
│ │ └── timestamp.go # Timestamp value object
│ │
│ ├── events/ # Shared domain events
│ │ ├── event.go # Base event interface
│ │ └── publisher.go # Event publisher
│ │
│ └── errors/ # Domain errors
│ ├── domainerror.go # Base domain error
│ └── specific/ # Specific domain errors
│ ├── notfound.go
│ ├── validation.go
│ └── limitsexceeded.go
│
└── service/ # Domain services
├── tenantevaluator.go # Tenant limits evaluator
├── albumorganizer.go # Album organization service
├── photoprocessor.go # Photo processing service
└── billingcalculator.go # Billing calculation service
text


## 🎪 Aggregates

### Tenant Aggregate
The Tenant aggregate represents a photography business using our platform. It's the root of all tenant-specific operations.

```go
// internal/domain/tenant/tenant.go
package tenant

import (
    "time"
    "github.com/google/uuid"
)

// Tenant is the root aggregate for a photography business
type Tenant struct {
    id        ID
    name      string
    slug      string
    plan      Plan
    subscription *Subscription
    limits    Limits
    usage     Usage
    settings  Settings
    createdAt time.Time
    updatedAt time.Time
}

// NewTenant creates a new tenant with default settings
func NewTenant(name, slug string, plan Plan) (*Tenant, error) {
    if name == "" {
        return nil, ErrTenantNameRequired
    }
    if slug == "" {
        return nil, ErrTenantSlugRequired
    }
    
    tenant := &Tenant{
        id:        ID(uuid.New().String()),
        name:      name,
        slug:      slug,
        plan:      plan,
        subscription: NewSubscription(plan),
        limits:    plan.DefaultLimits(),
        usage:     NewUsage(),
        settings:  DefaultSettings(),
        createdAt: time.Now(),
        updatedAt: time.Now(),
    }
    
    // Publish domain event
    tenant.AddEvent(TenantCreated{
        TenantID:  tenant.id,
        Name:      tenant.name,
        Plan:      tenant.plan,
        CreatedAt: tenant.createdAt,
    })
    
    return tenant, nil
}

// ChangePlan changes the tenant's subscription plan
func (t *Tenant) ChangePlan(newPlan Plan) error {
    if t.plan == newPlan {
        return nil // No change needed
    }
    
    // Validate the plan change is allowed
    if err := t.subscription.CanChangeTo(newPlan); err != nil {
        return err
    }
    
    oldPlan := t.plan
    t.plan = newPlan
    t.limits = newPlan.DefaultLimits()
    t.updatedAt = time.Now()
    
    // Update subscription
    t.subscription.ChangePlan(newPlan)
    
    // Publish domain event
    t.AddEvent(TenantPlanChanged{
        TenantID:  t.id,
        OldPlan:   oldPlan,
        NewPlan:   newPlan,
        ChangedAt: t.updatedAt,
    })
    
    return nil
}

// CheckStorageLimit checks if tenant can upload more data
func (t *Tenant) CheckStorageLimit(additionalBytes int64) error {
    newTotal := t.usage.StorageBytes + additionalBytes
    limitBytes := t.limits.StorageGB * 1024 * 1024 * 1024
    
    if newTotal > limitBytes {
        return ErrStorageLimitExceeded{
            Current:    t.usage.StorageBytes,
            Additional: additionalBytes,
            Limit:      limitBytes,
        }
    }
    
    return nil
}

// RecordUsage updates tenant usage statistics
func (t *Tenant) RecordUsage(usageDelta Usage) error {
    // Validate against limits
    if err := t.limits.Validate(t.usage.Add(usageDelta)); err != nil {
        return err
    }
    
    t.usage = t.usage.Add(usageDelta)
    t.updatedAt = time.Now()
    
    // Check if we need to warn about approaching limits
    if t.usage.IsApproachingLimit(t.limits) {
        t.AddEvent(TenantApproachingLimit{
            TenantID:   t.id,
            Usage:      t.usage,
            Limits:     t.limits,
            CheckedAt:  time.Now(),
        })
    }
    
    return nil
}

Album Aggregate

The Album aggregate manages a collection of photos with business rules around visibility, sharing, and lifecycle.
go

// internal/domain/album/album.go
package album

import (
    "time"
    "github.com/google/uuid"
)

// Album is the aggregate root for a photo collection
type Album struct {
    id          ID
    tenantID    tenant.ID
    title       string
    description string
    slug        string
    status      Status
    visibility  Visibility
    photos      []*Photo
    settings    Settings
    statistics  Statistics
    createdAt   time.Time
    updatedAt   time.Time
    publishedAt *time.Time
    archivedAt  *time.Time
    version     int
    events      []interface{} // Domain events
}

// NewAlbum creates a new album in draft state
func NewAlbum(tenantID tenant.ID, title, description string) (*Album, error) {
    if title == "" {
        return nil, ErrAlbumTitleRequired
    }
    
    album := &Album{
        id:          ID(uuid.New().String()),
        tenantID:    tenantID,
        title:       title,
        description: description,
        slug:        generateSlug(title),
        status:      StatusDraft,
        visibility:  VisibilityPrivate,
        photos:      []*Photo{},
        settings:    DefaultSettings(),
        statistics:  NewStatistics(),
        createdAt:   time.Now(),
        updatedAt:   time.Now(),
        version:     1,
        events:      []interface{}{},
    }
    
    // Add creation event
    album.AddEvent(AlbumCreated{
        AlbumID:     album.id,
        TenantID:    album.tenantID,
        Title:       album.title,
        CreatedAt:   album.createdAt,
    })
    
    return album, nil
}

// Publish makes the album publicly accessible
func (a *Album) Publish() error {
    if a.status == StatusArchived {
        return ErrCannotPublishArchivedAlbum
    }
    
    if len(a.photos) == 0 {
        return ErrCannotPublishEmptyAlbum
    }
    
    if a.visibility == VisibilityClient && len(a.photos) < 5 {
        return ErrNotEnoughPhotosForClientAlbum
    }
    
    a.status = StatusPublished
    now := time.Now()
    a.publishedAt = &now
    a.updatedAt = now
    a.version++
    
    // Add publish event
    a.AddEvent(AlbumPublished{
        AlbumID:     a.id,
        TenantID:    a.tenantID,
        PublishedAt: now,
        PhotoCount:  len(a.photos),
    })
    
    return nil
}

// AddPhoto adds a photo to the album
func (a *Album) AddPhoto(photo *Photo) error {
    if a.status == StatusArchived {
        return ErrCannotModifyArchivedAlbum
    }
    
    // Check for duplicate photos
    for _, p := range a.photos {
        if p.FileHash == photo.FileHash {
            return ErrDuplicatePhoto
        }
    }
    
    // Set photo order
    photo.Order = len(a.photos)
    photo.AlbumID = a.id
    photo.CreatedAt = time.Now()
    
    a.photos = append(a.photos, photo)
    a.updatedAt = time.Now()
    a.version++
    
    // Update statistics
    a.statistics.PhotoCount++
    a.statistics.TotalSizeBytes += photo.FileSize
    
    // Add photo added event
    a.AddEvent(PhotoAdded{
        AlbumID:     a.id,
        TenantID:    a.tenantID,
        PhotoID:     photo.ID,
        Filename:    photo.OriginalFilename,
        FileSize:    photo.FileSize,
        AddedAt:     time.Now(),
    })
    
    return nil
}

// ChangeVisibility changes who can see the album
func (a *Album) ChangeVisibility(visibility Visibility, password string) error {
    if a.status == StatusArchived {
        return ErrCannotModifyArchivedAlbum
    }
    
    if visibility == VisibilityPassword && password == "" {
        return ErrPasswordRequired
    }
    
    oldVisibility := a.visibility
    a.visibility = visibility
    
    if visibility == VisibilityPassword {
        a.settings.PasswordHash = hashPassword(password)
    } else {
        a.settings.PasswordHash = ""
    }
    
    a.updatedAt = time.Now()
    a.version++
    
    // Add visibility changed event
    a.AddEvent(AlbumVisibilityChanged{
        AlbumID:         a.id,
        TenantID:        a.tenantID,
        OldVisibility:   oldVisibility,
        NewVisibility:   visibility,
        ChangedAt:       time.Now(),
    })
    
    return nil
}

// ValidatePassword checks if the provided password is correct
func (a *Album) ValidatePassword(password string) bool {
    if a.visibility != VisibilityPassword {
        return true // No password required
    }
    
    return checkPasswordHash(password, a.settings.PasswordHash)
}

Photo Entity

The Photo entity represents an individual photo within an album.
go

// internal/domain/album/photo.go
package album

import (
    "time"
    "mime"
    "path/filepath"
    "github.com/google/uuid"
)

// Photo represents a single photo in an album
type Photo struct {
    id                ID
    tenantID          tenant.ID
    albumID           ID
    originalFilename  string
    fileHash          string // SHA-256 for deduplication
    fileSize          int64
    mimeType          string
    dimensions        Dimensions
    exif              EXIFData
    metadata          Metadata
    versions          Versions
    processingStatus  ProcessingStatus
    processingErrors  []string
    order             int
    createdAt         time.Time
    updatedAt         time.Time
    uploadedAt        time.Time
}

// NewPhoto creates a new photo entity
func NewPhoto(
    tenantID tenant.ID,
    filename string,
    fileSize int64,
    fileHash string,
) (*Photo, error) {
    
    if filename == "" {
        return nil, ErrPhotoFilenameRequired
    }
    
    if fileSize <= 0 {
        return nil, ErrPhotoFileSizeInvalid
    }
    
    if fileHash == "" {
        return nil, ErrPhotoFileHashRequired
    }
    
    // Validate file type
    mimeType := mime.TypeByExtension(filepath.Ext(filename))
    if !isSupportedImageType(mimeType) {
        return nil, ErrUnsupportedImageType{MimeType: mimeType}
    }
    
    // Validate file size (max 100MB)
    if fileSize > 100*1024*1024 {
        return nil, ErrFileTooLarge{MaxSize: 100 * 1024 * 1024}
    }
    
    photo := &Photo{
        id:               ID(uuid.New().String()),
        tenantID:         tenantID,
        originalFilename: filename,
        fileHash:         fileHash,
        fileSize:         fileSize,
        mimeType:         mimeType,
        processingStatus: ProcessingStatusPending,
        processingErrors: []string{},
        metadata:         NewMetadata(),
        versions:         NewVersions(),
        createdAt:        time.Now(),
        updatedAt:        time.Now(),
        uploadedAt:       time.Now(),
    }
    
    return photo, nil
}

// UpdateMetadata updates photo metadata
func (p *Photo) UpdateMetadata(metadata Metadata) error {
    if p.processingStatus != ProcessingStatusCompleted {
        return ErrCannotUpdateUnprocessedPhoto
    }
    
    // Validate metadata
    if err := metadata.Validate(); err != nil {
        return err
    }
    
    p.metadata = metadata
    p.updatedAt = time.Now()
    
    return nil
}

// MarkProcessingStarted marks the photo as being processed
func (p *Photo) MarkProcessingStarted() error {
    if p.processingStatus != ProcessingStatusPending {
        return ErrPhotoAlreadyProcessing
    }
    
    p.processingStatus = ProcessingStatusProcessing
    p.updatedAt = time.Now()
    
    return nil
}

// MarkProcessingCompleted marks the photo processing as complete
func (p *Photo) MarkProcessingCompleted(versions Versions, exif EXIFData) error {
    if p.processingStatus != ProcessingStatusProcessing {
        return ErrPhotoNotProcessing
    }
    
    p.processingStatus = ProcessingStatusCompleted
    p.versions = versions
    p.exif = exif
    p.dimensions = versions.Web.Dimensions
    p.updatedAt = time.Now()
    
    return nil
}

// MarkProcessingFailed marks the photo processing as failed
func (p *Photo) MarkProcessingFailed(err error) {
    p.processingStatus = ProcessingStatusFailed
    p.processingErrors = append(p.processingErrors, err.Error())
    p.updatedAt = time.Now()
}

// GetURL returns the URL for a specific version of the photo
func (p *Photo) GetURL(version string, includeWatermark bool) string {
    if !p.IsProcessed() {
        return "" // Not available until processed
    }
    
    url := p.versions.GetURL(version)
    if includeWatermark && p.metadata.Watermark.Enabled {
        url = p.metadata.Watermark.ApplyToURL(url)
    }
    
    return url
}

// IsProcessed returns true if the photo has been successfully processed
func (p *Photo) IsProcessed() bool {
    return p.processingStatus == ProcessingStatusCompleted
}

💎 Value Objects
ID Value Object
go

// internal/domain/shared/valueobject/id.go
package valueobject

import (
    "fmt"
    "github.com/google/uuid"
)

// ID is a value object representing an identifier
type ID string

// NewID creates a new unique ID
func NewID() ID {
    return ID(uuid.New().String())
}

// String returns the string representation
func (id ID) String() string {
    return string(id)
}

// IsEmpty returns true if the ID is empty
func (id ID) IsEmpty() bool {
    return string(id) == ""
}

// Validate validates the ID format
func (id ID) Validate() error {
    if id.IsEmpty() {
        return fmt.Errorf("id cannot be empty")
    }
    
    // Check if it's a valid UUID
    if _, err := uuid.Parse(string(id)); err != nil {
        return fmt.Errorf("invalid id format: %w", err)
    }
    
    return nil
}

// Equals checks if two IDs are equal
func (id ID) Equals(other ID) bool {
    return id == other
}

Metadata Value Object
go

// internal/domain/shared/valueobject/metadata.go
package valueobject

import (
    "encoding/json"
    "fmt"
)

// Metadata is a value object for flexible metadata storage
type Metadata map[string]interface{}

// NewMetadata creates new empty metadata
func NewMetadata() Metadata {
    return make(Metadata)
}

// GetString returns a string value from metadata
func (m Metadata) GetString(key string) (string, bool) {
    if val, exists := m[key]; exists {
        if str, ok := val.(string); ok {
            return str, true
        }
    }
    return "", false
}

// SetString sets a string value in metadata
func (m Metadata) SetString(key, value string) {
    m[key] = value
}

// GetInt returns an integer value from metadata
func (m Metadata) GetInt(key string) (int, bool) {
    if val, exists := m[key]; exists {
        switch v := val.(type) {
        case int:
            return v, true
        case float64:
            return int(v), true
        }
    }
    return 0, false
}

// Merge merges another metadata into this one
func (m Metadata) Merge(other Metadata) Metadata {
    result := make(Metadata)
    
    // Copy current metadata
    for k, v := range m {
        result[k] = v
    }
    
    // Merge other metadata (overwrites existing keys)
    for k, v := range other {
        result[k] = v
    }
    
    return result
}

// Validate validates metadata against schema
func (m Metadata) Validate(schema map[string]interface{}) error {
    // Implementation would validate against JSON schema
    // or custom validation rules
    return nil
}

// ToJSON converts metadata to JSON
func (m Metadata) ToJSON() ([]byte, error) {
    return json.Marshal(m)
}

// FromJSON parses metadata from JSON
func FromJSON(data []byte) (Metadata, error) {
    var m Metadata
    if err := json.Unmarshal(data, &m); err != nil {
        return nil, fmt.Errorf("failed to parse metadata: %w", err)
    }
    return m, nil
}

File Information Value Object
go

// internal/domain/shared/valueobject/file.go
package valueobject

import (
    "fmt"
    "mime"
    "path/filepath"
    "strings"
)

// FileInfo represents file information
type FileInfo struct {
    Filename    string
    Extension   string
    MimeType    string
    SizeBytes   int64
}

// NewFileInfo creates a new FileInfo from filename and size
func NewFileInfo(filename string, sizeBytes int64) (FileInfo, error) {
    if filename == "" {
        return FileInfo{}, fmt.Errorf("filename is required")
    }
    
    if sizeBytes < 0 {
        return FileInfo{}, fmt.Errorf("file size cannot be negative")
    }
    
    ext := strings.ToLower(filepath.Ext(filename))
    mimeType := mime.TypeByExtension(ext)
    
    return FileInfo{
        Filename:    filename,
        Extension:   ext,
        MimeType:    mimeType,
        SizeBytes:   sizeBytes,
    }, nil
}

// IsImage returns true if the file is an image
func (f FileInfo) IsImage() bool {
    return strings.HasPrefix(f.MimeType, "image/")
}

// IsSupportedImage returns true if the image type is supported
func (f FileInfo) IsSupportedImage() bool {
    if !f.IsImage() {
        return false
    }
    
    supportedTypes := map[string]bool{
        "image/jpeg": true,
        "image/png":  true,
        "image/webp": true,
        "image/gif":  true,
        "image/tiff": true,
        "image/bmp":  true,
    }
    
    return supportedTypes[f.MimeType]
}

// Validate validates the file information
func (f FileInfo) Validate() error {
    if f.Filename == "" {
        return fmt.Errorf("filename is required")
    }
    
    if f.SizeBytes <= 0 {
        return fmt.Errorf("file size must be positive")
    }
    
    if !f.IsSupportedImage() {
        return fmt.Errorf("unsupported image type: %s", f.MimeType)
    }
    
    // Check file size limits
    maxSize := int64(100 * 1024 * 1024) // 100MB
    if f.SizeBytes > maxSize {
        return fmt.Errorf("file size exceeds limit of %dMB", maxSize/1024/1024)
    }
    
    return nil
}

// BaseName returns the filename without extension
func (f FileInfo) BaseName() string {
    return strings.TrimSuffix(f.Filename, f.Extension)
}

🌟 Domain Events
Base Event Interface
go

// internal/domain/shared/events/event.go
package events

import (
    "time"
)

// Event represents a domain event
type Event interface {
    EventID() string
    AggregateID() string
    AggregateType() string
    EventType() string
    TenantID() string
    Timestamp() time.Time
    Data() interface{}
    Metadata() map[string]interface{}
    Version() int
}

// BaseEvent provides common event functionality
type BaseEvent struct {
    eventID      string
    aggregateID  string
    aggregateType string
    eventType    string
    tenantID     string
    timestamp    time.Time
    metadata     map[string]interface{}
    version      int
}

func (e BaseEvent) EventID() string {
    return e.eventID
}

func (e BaseEvent) AggregateID() string {
    return e.aggregateID
}

func (e BaseEvent) AggregateType() string {
    return e.aggregateType
}

func (e BaseEvent) EventType() string {
    return e.eventType
}

func (e BaseEvent) TenantID() string {
    return e.tenantID
}

func (e BaseEvent) Timestamp() time.Time {
    return e.timestamp
}

func (e BaseEvent) Metadata() map[string]interface{} {
    return e.metadata
}

func (e BaseEvent) Version() int {
    return e.version
}

// NewBaseEvent creates a new base event
func NewBaseEvent(
    aggregateID string,
    aggregateType string,
    eventType string,
    tenantID string,
    version int,
) BaseEvent {
    
    return BaseEvent{
        eventID:       generateEventID(),
        aggregateID:   aggregateID,
        aggregateType: aggregateType,
        eventType:     eventType,
        tenantID:      tenantID,
        timestamp:     time.Now(),
        metadata:      make(map[string]interface{}),
        version:       version,
    }
}

Album Domain Events
go

// internal/domain/album/events.go
package album

import (
    "time"
    "github.com/google/uuid"
    "github.com/photolisting/internal/domain/shared/events"
    "github.com/photolisting/internal/domain/tenant"
)

// AlbumCreated event
type AlbumCreated struct {
    events.BaseEvent
    Title       string
    Description string
    Visibility  Visibility
}

func NewAlbumCreated(albumID ID, tenantID tenant.ID, title, description string) AlbumCreated {
    return AlbumCreated{
        BaseEvent: events.NewBaseEvent(
            albumID.String(),
            "album",
            "album.created",
            tenantID.String(),
            1,
        ),
        Title:       title,
        Description: description,
        Visibility:  VisibilityPrivate,
    }
}

func (e AlbumCreated) Data() interface{} {
    return map[string]interface{}{
        "title":       e.Title,
        "description": e.Description,
        "visibility":  e.Visibility,
    }
}

// AlbumPublished event
type AlbumPublished struct {
    events.BaseEvent
    PublishedAt time.Time
    PhotoCount  int
    PublicURL   string
}

func NewAlbumPublished(albumID ID, tenantID tenant.ID, photoCount int, publicURL string) AlbumPublished {
    return AlbumPublished{
        BaseEvent: events.NewBaseEvent(
            albumID.String(),
            "album",
            "album.published",
            tenantID.String(),
            1,
        ),
        PublishedAt: time.Now(),
        PhotoCount:  photoCount,
        PublicURL:   publicURL,
    }
}

func (e AlbumPublished) Data() interface{} {
    return map[string]interface{}{
        "published_at": e.PublishedAt,
        "photo_count":  e.PhotoCount,
        "public_url":   e.PublicURL,
    }
}

// PhotoAdded event
type PhotoAdded struct {
    events.BaseEvent
    PhotoID     ID
    Filename    string
    FileSize    int64
    AddedAt     time.Time
}

func NewPhotoAdded(albumID, photoID ID, tenantID tenant.ID, filename string, fileSize int64) PhotoAdded {
    return PhotoAdded{
        BaseEvent: events.NewBaseEvent(
            albumID.String(),
            "album",
            "photo.added",
            tenantID.String(),
            1,
        ),
        PhotoID:     photoID,
        Filename:    filename,
        FileSize:    fileSize,
        AddedAt:     time.Now(),
    }
}

func (e PhotoAdded) Data() interface{} {
    return map[string]interface{}{
        "photo_id":  e.PhotoID.String(),
        "filename":  e.Filename,
        "file_size": e.FileSize,
        "added_at":  e.AddedAt,
    }
}

// AlbumVisibilityChanged event
type AlbumVisibilityChanged struct {
    events.BaseEvent
    OldVisibility Visibility
    NewVisibility Visibility
    ChangedAt     time.Time
}

func NewAlbumVisibilityChanged(albumID ID, tenantID tenant.ID, oldVis, newVis Visibility) AlbumVisibilityChanged {
    return AlbumVisibilityChanged{
        BaseEvent: events.NewBaseEvent(
            albumID.String(),
            "album",
            "album.visibility_changed",
            tenantID.String(),
            1,
        ),
        OldVisibility: oldVis,
        NewVisibility: newVis,
        ChangedAt:     time.Now(),
    }
}

func (e AlbumVisibilityChanged) Data() interface{} {
    return map[string]interface{}{
        "old_visibility": e.OldVisibility,
        "new_visibility": e.NewVisibility,
        "changed_at":     e.ChangedAt,
    }
}

🛡️ Domain Errors
Base Domain Error
go

// internal/domain/shared/errors/domainerror.go
package errors

import "fmt"

// DomainError is the base error type for domain errors
type DomainError struct {
    Code    string
    Message string
    Details map[string]interface{}
}

func (e DomainError) Error() string {
    return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewDomainError creates a new domain error
func NewDomainError(code, message string, details map[string]interface{}) DomainError {
    return DomainError{
        Code:    code,
        Message: message,
        Details: details,
    }
}

// IsDomainError checks if an error is a domain error
func IsDomainError(err error) bool {
    _, ok := err.(DomainError)
    return ok
}

Specific Domain Errors
go

// internal/domain/shared/errors/specific/notfound.go
package specific

import "github.com/photolisting/internal/domain/shared/errors"

var (
    // Tenant errors
    ErrTenantNotFound = errors.NewDomainError(
        "TENANT_NOT_FOUND",
        "Tenant not found",
        nil,
    )
    
    // Album errors
    ErrAlbumNotFound = errors.NewDomainError(
        "ALBUM_NOT_FOUND",
        "Album not found",
        nil,
    )
    
    // Photo errors
    ErrPhotoNotFound = errors.NewDomainError(
        "PHOTO_NOT_FOUND",
        "Photo not found",
        nil,
    )
    
    // User errors
    ErrUserNotFound = errors.NewDomainError(
        "USER_NOT_FOUND",
        "User not found",
        nil,
    )
)

// internal/domain/shared/errors/specific/validation.go
package specific

import "github.com/photolisting/internal/domain/shared/errors"

var (
    // Validation errors
    ErrInvalidInput = errors.NewDomainError(
        "INVALID_INPUT",
        "Invalid input provided",
        nil,
    )
    
    ErrValidationFailed = errors.NewDomainError(
        "VALIDATION_FAILED",
        "Validation failed",
        nil,
    )
    
    // Album validation errors
    ErrAlbumTitleRequired = errors.NewDomainError(
        "ALBUM_TITLE_REQUIRED",
        "Album title is required",
        nil,
    )
    
    ErrCannotPublishEmptyAlbum = errors.NewDomainError(
        "CANNOT_PUBLISH_EMPTY_ALBUM",
        "Cannot publish an empty album",
        nil,
    )
    
    // Photo validation errors
    ErrPhotoFilenameRequired = errors.NewDomainError(
        "PHOTO_FILENAME_REQUIRED",
        "Photo filename is required",
        nil,
    )
    
    ErrUnsupportedImageType = errors.NewDomainError(
        "UNSUPPORTED_IMAGE_TYPE",
        "Unsupported image type",
        nil,
    )
    
    ErrFileTooLarge = errors.NewDomainError(
        "FILE_TOO_LARGE",
        "File size exceeds maximum allowed size",
        nil,
    )
)

// internal/domain/shared/errors/specific/limitsexceeded.go
package specific

import "github.com/photolisting/internal/domain/shared/errors"

// StorageLimitExceeded error
type StorageLimitExceeded struct {
    errors.DomainError
    Current    int64
    Additional int64
    Limit      int64
}

func NewStorageLimitExceeded(current, additional, limit int64) StorageLimitExceeded {
    return StorageLimitExceeded{
        DomainError: errors.NewDomainError(
            "STORAGE_LIMIT_EXCEEDED",
            "Storage limit exceeded",
            map[string]interface{}{
                "current":    current,
                "additional": additional,
                "limit":      limit,
                "unit":       "bytes",
            },
        ),
        Current:    current,
        Additional: additional,
        Limit:      limit,
    }
}

// AlbumLimitExceeded error
type AlbumLimitExceeded struct {
    errors.DomainError
    Current    int
    Limit      int
}

func NewAlbumLimitExceeded(current, limit int) AlbumLimitExceeded {
    return AlbumLimitExceeded{
        DomainError: errors.NewDomainError(
            "ALBUM_LIMIT_EXCEEDED",
            "Album limit exceeded",
            map[string]interface{}{
                "current": current,
                "limit":   limit,
            },
        ),
        Current: current,
        Limit:   limit,
    }
}

🔧 Domain Services
Tenant Evaluator Service
go

// internal/domain/service/tenantevaluator.go
package service

import (
    "context"
    "github.com/photolisting/internal/domain/tenant"
    "github.com/photolisting/internal/domain/album"
)

// TenantEvaluator evaluates tenant limits and permissions
type TenantEvaluator struct {
    tenantRepo tenant.Repository
}

// NewTenantEvaluator creates a new tenant evaluator
func NewTenantEvaluator(tenantRepo tenant.Repository) *TenantEvaluator {
    return &TenantEvaluator{
        tenantRepo: tenantRepo,
    }
}

// CanCreateAlbum checks if a tenant can create a new album
func (e *TenantEvaluator) CanCreateAlbum(ctx context.Context, tenantID tenant.ID) error {
    t, err := e.tenantRepo.FindByID(ctx, tenantID)
    if err != nil {
        return err
    }
    
    // Check album limit
    if t.Limits.MaxAlbums > 0 && t.Usage.AlbumCount >= t.Limits.MaxAlbums {
        return tenant.ErrAlbumLimitExceeded{
            Current: t.Usage.AlbumCount,
            Limit:   t.Limits.MaxAlbums,
        }
    }
    
    return nil
}

// CanUploadPhoto checks if a tenant can upload a photo
func (e *TenantEvaluator) CanUploadPhoto(
    ctx context.Context,
    tenantID tenant.ID,
    fileSize int64,
) error {
    
    t, err := e.tenantRepo.FindByID(ctx, tenantID)
    if err != nil {
        return err
    }
    
    // Check storage limit
    if err := t.CheckStorageLimit(fileSize); err != nil {
        return err
    }
    
    // Check photo limit per album
    // This would depend on the specific album being uploaded to
    
    return nil
}

// EvaluateTenantStatus evaluates overall tenant status
func (e *TenantEvaluator) EvaluateTenantStatus(ctx context.Context, tenantID tenant.ID) (tenant.Status, error) {
    t, err := e.tenantRepo.FindByID(ctx, tenantID)
    if err != nil {
        return tenant.StatusUnknown, err
    }
    
    // Check subscription status
    if !t.Subscription.IsActive() {
        return tenant.StatusInactive, nil
    }
    
    // Check if over limits
    if t.Usage.IsOverLimit(t.Limits) {
        return tenant.StatusOverLimit, nil
    }
    
    // Check if approaching limits
    if t.Usage.IsApproachingLimit(t.Limits) {
        return tenant.StatusApproachingLimit, nil
    }
    
    return tenant.StatusActive, nil
}

Album Organizer Service
go

// internal/domain/service/albumorganizer.go
package service

import (
    "context"
    "sort"
    "github.com/photolisting/internal/domain/album"
)

// AlbumOrganizer provides album organization operations
type AlbumOrganizer struct {
    albumRepo album.Repository
}

// NewAlbumOrganizer creates a new album organizer
func NewAlbumOrganizer(albumRepo album.Repository) *AlbumOrganizer {
    return &AlbumOrganizer{
        albumRepo: albumRepo,
    }
}

// ReorderPhotos reorders photos in an album
func (o *AlbumOrganizer) ReorderPhotos(
    ctx context.Context,
    albumID album.ID,
    photoOrder map[album.ID]int,
) error {
    
    a, err := o.albumRepo.FindByID(ctx, albumID)
    if err != nil {
        return err
    }
    
    // Validate all photo IDs exist in the album
    for photoID := range photoOrder {
        found := false
        for _, photo := range a.Photos {
            if photo.ID == photoID {
                found = true
                break
            }
        }
        if !found {
            return album.ErrPhotoNotFound
        }
    }
    
    // Reorder photos
    for i, photo := range a.Photos {
        if newOrder, exists := photoOrder[photo.ID]; exists {
            a.Photos[i].Order = newOrder
        }
    }
    
    // Sort photos by new order
    sort.Slice(a.Photos, func(i, j int) bool {
        return a.Photos[i].Order < a.Photos[j].Order
    })
    
    // Save changes
    return o.albumRepo.Save(ctx, a)
}

// MovePhotos moves photos between albums
func (o *AlbumOrganizer) MovePhotos(
    ctx context.Context,
    sourceAlbumID album.ID,
    targetAlbumID album.ID,
    photoIDs []album.ID,
) error {
    
    sourceAlbum, err := o.albumRepo.FindByID(ctx, sourceAlbumID)
    if err != nil {
        return err
    }
    
    targetAlbum, err := o.albumRepo.FindByID(ctx, targetAlbumID)
    if err != nil {
        return err
    }
    
    // Check if albums belong to same tenant
    if sourceAlbum.TenantID != targetAlbum.TenantID {
        return album.ErrCannotMoveBetweenTenants
    }
    
    // Move photos
    var photosToMove []*album.Photo
    var remainingPhotos []*album.Photo
    
    for _, photo := range sourceAlbum.Photos {
        move := false
        for _, photoID := range photoIDs {
            if photo.ID == photoID {
                move = true
                break
            }
        }
        
        if move {
            photosToMove = append(photosToMove, photo)
        } else {
            remainingPhotos = append(remainingPhotos, photo)
        }
    }
    
    // Update source album
    sourceAlbum.Photos = remainingPhotos
    sourceAlbum.UpdateStatistics()
    
    // Update target album
    for _, photo := range photosToMove {
        photo.AlbumID = targetAlbumID
        photo.Order = len(targetAlbum.Photos)
        targetAlbum.Photos = append(targetAlbum.Photos, photo)
    }
    targetAlbum.UpdateStatistics()
    
    // Save both albums
    if err := o.albumRepo.Save(ctx, sourceAlbum); err != nil {
        return err
    }
    
    return o.albumRepo.Save(ctx, targetAlbum)
}

// CreateAlbumFromSelection creates a new album from selected photos
func (o *AlbumOrganizer) CreateAlbumFromSelection(
    ctx context.Context,
    tenantID tenant.ID,
    sourceAlbumID album.ID,
    photoIDs []album.ID,
    title string,
    description string,
) (*album.Album, error) {
    
    sourceAlbum, err := o.albumRepo.FindByID(ctx, sourceAlbumID)
    if err != nil {
        return nil, err
    }
    
    // Create new album
    newAlbum, err := album.NewAlbum(tenantID, title, description)
    if err != nil {
        return nil, err
    }
    
    // Copy selected photos
    for _, photo := range sourceAlbum.Photos {
        for _, photoID := range photoIDs {
            if photo.ID == photoID {
                // Create copy of photo
                copiedPhoto := *photo
                copiedPhoto.ID = album.NewID()
                copiedPhoto.AlbumID = newAlbum.ID
                copiedPhoto.Order = len(newAlbum.Photos)
                
                newAlbum.Photos = append(newAlbum.Photos, &copiedPhoto)
                break
            }
        }
    }
    
    // Save new album
    if err := o.albumRepo.Save(ctx, newAlbum); err != nil {
        return nil, err
    }
    
    return newAlbum, nil
}

📚 Repository Interfaces
Album Repository Interface
go

// internal/domain/album/repository.go
package album

import (
    "context"
    "github.com/photolisting/internal/domain/tenant"
)

// Repository defines the interface for album persistence
type Repository interface {
    // Basic CRUD operations
    FindByID(ctx context.Context, id ID) (*Album, error)
    Save(ctx context.Context, album *Album) error
    Delete(ctx context.Context, id ID) error
    
    // Query operations
    FindBySlug(ctx context.Context, tenantID tenant.ID, slug string) (*Album, error)
    ListByTenant(ctx context.Context, tenantID tenant.ID, filter ListFilter) ([]*Album, *Pagination, error)
    ListByStatus(ctx context.Context, tenantID tenant.ID, status Status) ([]*Album, error)
    ListPublicAlbums(ctx context.Context, tenantID tenant.ID) ([]*Album, error)
    
    // Photo operations
    FindPhotoByID(ctx context.Context, id ID) (*Photo, error)
    FindPhotosByAlbum(ctx context.Context, albumID ID, filter PhotoFilter) ([]*Photo, *Pagination, error)
    SavePhoto(ctx context.Context, photo *Photo) error
    DeletePhoto(ctx context.Context, id ID) error
    
    // Statistics
    CountByTenant(ctx context.Context, tenantID tenant.ID) (int, error)
    GetStorageUsage(ctx context.Context, tenantID tenant.ID) (int64, error)
    
    // Transaction support
    InTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// ListFilter represents album listing filters
type ListFilter struct {
    Status     *Status
    Visibility *Visibility
    Category   *string
    Tags       []string
    Search     *string
    SortBy     string // "created_at", "updated_at", "title", "published_at"
    SortOrder  string // "asc", "desc"
    Page       int
    Limit      int
}

// PhotoFilter represents photo listing filters
type PhotoFilter struct {
    Status            *ProcessingStatus
    HasExif           *bool
    HasLocation       *bool
    MinWidth          *int
    MinHeight         *int
    MetadataTags      []string
    SortBy            string // "created_at", "uploaded_at", "taken_at", "filename", "size"
    SortOrder         string // "asc", "desc"
    Page              int
    Limit             int
}

// Pagination represents pagination metadata
type Pagination struct {
    Page       int
    Limit      int
    Total      int64
    TotalPages int64
}

Tenant Repository Interface
go

// internal/domain/tenant/repository.go
package tenant

import (
    "context"
)

// Repository defines the interface for tenant persistence
type Repository interface {
    // Basic CRUD operations
    FindByID(ctx context.Context, id ID) (*Tenant, error)
    FindBySlug(ctx context.Context, slug string) (*Tenant, error)
    Save(ctx context.Context, tenant *Tenant) error
    Delete(ctx context.Context, id ID) error
    
    // Query operations
    List(ctx context.Context, filter ListFilter) ([]*Tenant, *Pagination, error)
    FindBySubscriptionID(ctx context.Context, subscriptionID string) (*Tenant, error)
    
    // Usage operations
    UpdateUsage(ctx context.Context, tenantID ID, usageDelta Usage) error
    ResetUsage(ctx context.Context, tenantID ID, period string) error
    
    // Statistics
    Count(ctx context.Context) (int64, error)
    GetActiveCount(ctx context.Context) (int64, error)
    
    // Transaction support
    InTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// ListFilter represents tenant listing filters
type ListFilter struct {
    Plan       *Plan
    Status     *string
    Search     *string
    SortBy     string // "created_at", "name", "plan"
    SortOrder  string // "asc", "desc"
    Page       int
    Limit      int
}

🧪 Testing Domain Models
Album Tests
go

// internal/domain/album/album_test.go
package album_test

import (
    "testing"
    "github.com/google/uuid"
    "github.com/photolisting/internal/domain/album"
    "github.com/photolisting/internal/domain/tenant"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestAlbum_Create(t *testing.T) {
    tenantID := tenant.ID(uuid.New().String())
    
    t.Run("creates album successfully", func(t *testing.T) {
        album, err := album.NewAlbum(tenantID, "Test Album", "Test Description")
        
        require.NoError(t, err)
        assert.NotEmpty(t, album.ID)
        assert.Equal(t, "Test Album", album.Title)
        assert.Equal(t, "Test Description", album.Description)
        assert.Equal(t, album.StatusDraft, album.Status)
        assert.Equal(t, album.VisibilityPrivate, album.Visibility)
        assert.NotEmpty(t, album.Slug)
        assert.Len(t, album.Events(), 1)
    })
    
    t.Run("fails with empty title", func(t *testing.T) {
        album, err := album.NewAlbum(tenantID, "", "Description")
        
        assert.Error(t, err)
        assert.Nil(t, album)
        assert.Equal(t, album.ErrAlbumTitleRequired, err)
    })
    
    t.Run("generates URL-friendly slug", func(t *testing.T) {
        album, err := album.NewAlbum(tenantID, "My Awesome Album!", "Test")
        
        require.NoError(t, err)
        assert.Equal(t, "my-awesome-album", album.Slug)
    })
}

func TestAlbum_Publish(t *testing.T) {
    tenantID := tenant.ID(uuid.New().String())
    
    t.Run("publishes album successfully", func(t *testing.T) {
        album, _ := album.NewAlbum(tenantID, "Test Album", "")
        
        // Add some photos
        for i := 0; i < 5; i++ {
            photo, _ := album.NewPhoto(
                "photo.jpg",
                1024*1024,
                "hash123",
            )
            album.AddPhoto(photo)
        }
        
        err := album.Publish()
        
        require.NoError(t, err)
        assert.Equal(t, album.StatusPublished, album.Status)
        assert.NotNil(t, album.PublishedAt)
        assert.Len(t, album.Events(), 6) // Created + 5 photos + published
    })
    
    t.Run("fails to publish empty album", func(t *testing.T) {
        album, _ := album.NewAlbum(tenantID, "Test Album", "")
        
        err := album.Publish()
        
        assert.Error(t, err)
        assert.Equal(t, album.ErrCannotPublishEmptyAlbum, err)
        assert.Equal(t, album.StatusDraft, album.Status)
    })
    
    t.Run("fails to publish archived album", func(t *testing.T) {
        album, _ := album.NewAlbum(tenantID, "Test Album", "")
        album.Status = album.StatusArchived
        
        err := album.Publish()
        
        assert.Error(t, err)
        assert.Equal(t, album.ErrCannotPublishArchivedAlbum, err)
    })
}

func TestAlbum_AddPhoto(t *testing.T) {
    tenantID := tenant.ID(uuid.New().String())
    album, _ := album.NewAlbum(tenantID, "Test Album", "")
    
    t.Run("adds photo successfully", func(t *testing.T) {
        photo, err := album.NewPhoto("photo.jpg", 1024*1024, "hash123")
        require.NoError(t, err)
        
        err = album.AddPhoto(photo)
        
        require.NoError(t, err)
        assert.Len(t, album.Photos, 1)
        assert.Equal(t, photo.ID, album.Photos[0].ID)
        assert.Equal(t, 0, album.Photos[0].Order)
        assert.Equal(t, int64(1024*1024), album.Statistics.TotalSizeBytes)
    })
    
    t.Run("prevents duplicate photos", func(t *testing.T) {
        photo, _ := album.NewPhoto("photo2.jpg", 1024*1024, "hash123")
        
        err := album.AddPhoto(photo)
        
        assert.Error(t, err)
        assert.Equal(t, album.ErrDuplicatePhoto, err)
        assert.Len(t, album.Photos, 1) // Still only one photo
    })
    
    t.Run("maintains photo order", func(t *testing.T) {
        photo2, _ := album.NewPhoto("photo2.jpg", 2048*1024, "hash456")
        photo3, _ := album.NewPhoto("photo3.jpg", 3072*1024, "hash789")
        
        album.AddPhoto(photo2)
        album.AddPhoto(photo3)
        
        assert.Len(t, album.Photos, 3)
        assert.Equal(t, 1, album.Photos[1].Order)
        assert.Equal(t, 2, album.Photos[2].Order)
        assert.Equal(t, int64(6144*1024), album.Statistics.TotalSizeBytes)
    })
}

Photo Tests
go

// internal/domain/album/photo_test.go
package album_test

import (
    "testing"
    "github.com/google/uuid"
    "github.com/photolisting/internal/domain/album"
    "github.com/photolisting/internal/domain/tenant"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestPhoto_Create(t *testing.T) {
    tenantID := tenant.ID(uuid.New().String())
    
    t.Run("creates photo successfully", func(t *testing.T) {
        photo, err := album.NewPhoto(
            tenantID,
            "photo.jpg",
            1024*1024,
            "abc123def456",
        )
        
        require.NoError(t, err)
        assert.NotEmpty(t, photo.ID)
        assert.Equal(t, "photo.jpg", photo.OriginalFilename)
        assert.Equal(t, int64(1024*1024), photo.FileSize)
        assert.Equal(t, "abc123def456", photo.FileHash)
        assert.Equal(t, album.ProcessingStatusPending, photo.ProcessingStatus)
        assert.NotEmpty(t, photo.CreatedAt)
    })
    
    t.Run("validates file size", func(t *testing.T) {
        photo, err := album.NewPhoto(
            tenantID,
            "photo.jpg",
            0,
            "hash123",
        )
        
        assert.Error(t, err)
        assert.Nil(t, photo)
        assert.Equal(t, album.ErrPhotoFileSizeInvalid, err)
    })
    
    t.Run("validates file type", func(t *testing.T) {
        photo, err := album.NewPhoto(
            tenantID,
            "document.pdf",
            1024*1024,
            "hash123",
        )
        
        assert.Error(t, err)
        assert.Nil(t, photo)
        assert.IsType(t, album.ErrUnsupportedImageType{}, err)
    })
    
    t.Run("validates file hash", func(t *testing.T) {
        photo, err := album.NewPhoto(
            tenantID,
            "photo.jpg",
            1024*1024,
            "",
        )
        
        assert.Error(t, err)
        assert.Nil(t, photo)
        assert.Equal(t, album.ErrPhotoFileHashRequired, err)
    })
}

func TestPhoto_Processing(t *testing.T) {
    tenantID := tenant.ID(uuid.New().String())
    photo, _ := album.NewPhoto(tenantID, "photo.jpg", 1024*1024, "hash123")
    
    t.Run("starts processing", func(t *testing.T) {
        err := photo.MarkProcessingStarted()
        
        require.NoError(t, err)
        assert.Equal(t, album.ProcessingStatusProcessing, photo.ProcessingStatus)
    })
    
    t.Run("completes processing", func(t *testing.T) {
        versions := album.NewVersions()
        exif := album.EXIFData{"Camera": "Canon", "ISO": "100"}
        
        err := photo.MarkProcessingCompleted(versions, exif)
        
        require.NoError(t, err)
        assert.Equal(t, album.ProcessingStatusCompleted, photo.ProcessingStatus)
        assert.Equal(t, versions, photo.Versions)
        assert.Equal(t, exif, photo.EXIF)
    })
    
    t.Run("fails processing", func(t *testing.T) {
        photo2, _ := album.NewPhoto(tenantID, "photo2.jpg", 1024*1024, "hash456")
        photo2.MarkProcessingStarted()
        
        photo2.MarkProcessingFailed(assert.AnError)
        
        assert.Equal(t, album.ProcessingStatusFailed, photo2.ProcessingStatus)
        assert.Len(t, photo2.ProcessingErrors, 1)
    })
    
    t.Run("prevents invalid state transitions", func(t *testing.T) {
        photo3, _ := album.NewPhoto(tenantID, "photo3.jpg", 1024*1024, "hash789")
        
        // Try to complete without starting
        err := photo3.MarkProcessingCompleted(album.NewVersions(), album.EXIFData{})
        
        assert.Error(t, err)
        assert.Equal(t, album.ErrPhotoNotProcessing, err)
    })
}

🎯 Domain Rules Summary
Business Rules
Album Rules

    Album Creation: Must have a title; automatically gets a URL-friendly slug

    Album Publishing: Cannot publish empty albums or archived albums

    Photo Limits: Client albums require at least 5 photos for publishing

    Visibility Changes: Password-protected albums require a password

    Archiving: Archived albums cannot be modified

Photo Rules

    File Validation: Must be a supported image type under 100MB

    Deduplication: Photos with identical file hashes cannot be added twice

    Processing Flow: Must follow pending → processing → completed/failed states

    Metadata Updates: Can only update metadata after processing completes

Tenant Rules

    Resource Limits: Cannot exceed storage, album, or user limits

    Plan Changes: Must follow allowed upgrade/downgrade paths

    Subscription Status: Inactive subscriptions limit functionality

    Usage Tracking: Monitored for billing and limit enforcement

User Rules

    Role Hierarchy: Owner > Admin > Editor > Viewer

    Permission Inheritance: Higher roles inherit lower role permissions

    Tenant Isolation: Users can only access their own tenant's data

    Invitation Limits: Limited by tenant's user limit

Invariants

    Tenant Isolation: All operations are scoped to a tenant

    Data Consistency: Aggregate roots enforce consistency boundaries

    Immutable Events: Domain events cannot be changed after creation

    Version Control: Optimistic concurrency control for aggregates

🔮 Future Domain Extensions
Planned Extensions

    Advanced Tagging: Hierarchical tags, AI-powered auto-tagging

    Smart Albums: Dynamic albums based on rules (date, location, tags)

    Collaboration Features: Real-time commenting, annotation tools

    Advanced Analytics: Predictive insights, trend analysis

    Integration Framework: Third-party service integrations

Extension Points
go

// Plugin system for photo processing
type PhotoProcessorPlugin interface {
    CanProcess(photo *Photo) bool
    Process(ctx context.Context, photo *Photo) (*Photo, error)
    Priority() int
}

// Event handlers for extensibility
type EventHandler interface {
    Handle(ctx context.Context, event events.Event) error
    EventTypes() []string
}

// Rule engine for smart albums
type AlbumRule interface {
    Evaluate(photo *Photo) bool
    Description() string
}

The domain layer encapsulates the core business logic and rules of Photo Listing SaaS. All application logic should defer to the domain layer for business decisions, keeping the domain pure and focused on business rules rather than technical concerns.