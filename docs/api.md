
## /docs/api.md

```markdown
# API Reference

Photo Listing SaaS provides a RESTful API for managing albums, photos, users, and subscriptions. All API endpoints require authentication unless explicitly marked as public.

## 📋 Base Information

### Base URL

https://api.photolisting.dev/v1
text


For development:

http://localhost:8080/v1
text


### Authentication
All endpoints (except public ones) require a JWT token in the Authorization header:

```http
Authorization: Bearer <your-jwt-token>

Response Format

All responses are in JSON format with consistent structure:

Success Response:
json

{
  "data": {...},
  "meta": {
    "request_id": "req_123456",
    "timestamp": "2024-01-15T10:30:00Z"
  }
}

Error Response:
json

{
  "error": {
    "code": "NOT_FOUND",
    "message": "Album not found",
    "details": {"album_id": "550e8400-e29b-41d4-a716-446655440000"},
    "request_id": "req_123456",
    "timestamp": "2024-01-15T10:30:00Z"
  }
}

Error Codes
Code	Description	HTTP Status
UNAUTHORIZED	Missing or invalid authentication	401
FORBIDDEN	Insufficient permissions	403
NOT_FOUND	Resource not found	404
VALIDATION_ERROR	Invalid request data	422
RATE_LIMIT_EXCEEDED	Too many requests	429
INTERNAL_ERROR	Server error	500
Rate Limiting

    Free Tier: 100 requests/hour

    Professional: 1,000 requests/hour

    Studio: 10,000 requests/hour

Rate limit headers are included in responses:
http

X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 985
X-RateLimit-Reset: 1705312800

🏢 Tenants
Get Tenant Details
http

GET /tenants/{tenant_id}

Response:
json

{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John's Photography",
    "slug": "johns-photography",
    "plan": "professional",
    "settings": {
      "branding": {
        "logo_url": "https://...",
        "primary_color": "#3b82f6"
      },
      "features": {
        "client_proofing": true,
        "custom_domains": true,
        "watermarking": true
      }
    },
    "limits": {
      "storage_gb": 50,
      "albums": null,
      "users": 5,
      "monthly_bandwidth_gb": 100
    },
    "usage": {
      "storage_gb": 12.5,
      "albums_created": 8,
      "active_users": 3,
      "bandwidth_this_month_gb": 24.7
    },
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}

Update Tenant Settings
http

PATCH /tenants/{tenant_id}
Content-Type: application/json

{
  "settings": {
    "branding": {
      "primary_color": "#10b981"
    }
  }
}

Response:
json

{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "settings": {
      "branding": {
        "primary_color": "#10b981"
      }
    },
    "updated_at": "2024-01-15T10:35:00Z"
  }
}

📸 Albums
List Albums
http

GET /albums

Query Parameters:
Parameter	Type	Description
status	string	Filter by status (draft, review, published, archived)
visibility	string	Filter by visibility (public, private, password, client)
sort_by	string	Sort field (created_at, updated_at, title)
sort_order	string	Sort order (asc, desc)
page	integer	Page number (default: 1)
limit	integer	Items per page (default: 20, max: 100)

Response:
json

{
  "data": [
    {
      "id": "album_123",
      "title": "Summer Wedding",
      "description": "Beautiful outdoor wedding ceremony",
      "status": "published",
      "visibility": "public",
      "cover_photo_url": "https://...",
      "photo_count": 45,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z",
      "_links": {
        "self": "/albums/album_123",
        "photos": "/albums/album_123/photos",
        "share": "/albums/album_123/share"
      }
    }
  ],
  "meta": {
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 45,
      "total_pages": 3
    }
  }
}

Create Album
http

POST /albums
Content-Type: application/json

{
  "title": "Family Portraits 2024",
  "description": "Annual family portrait session",
  "visibility": "private",
  "password": "optional-for-password-protected",
  "settings": {
    "watermark": {
      "enabled": true,
      "text": "© John's Photography",
      "opacity": 0.3,
      "position": "bottom-right"
    },
    "download": {
      "enabled": true,
      "require_email": true,
      "max_downloads": 10
    },
    "proofing": {
      "enabled": false,
      "max_selections": 25
    }
  }
}

Response:
json

{
  "data": {
    "id": "album_456",
    "title": "Family Portraits 2024",
    "description": "Annual family portrait session",
    "status": "draft",
    "visibility": "private",
    "settings": {
      "watermark": {
        "enabled": true,
        "text": "© John's Photography",
        "opacity": 0.3,
        "position": "bottom-right"
      }
    },
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "_links": {
      "self": "/albums/album_456",
      "photos": "/albums/album_456/photos",
      "upload": "/albums/album_456/upload"
    }
  }
}

Get Album Details
http

GET /albums/{album_id}

Response:
json

{
  "data": {
    "id": "album_123",
    "title": "Summer Wedding",
    "description": "Beautiful outdoor wedding ceremony",
    "status": "published",
    "visibility": "public",
    "cover_photo": {
      "id": "photo_789",
      "url": "https://...",
      "thumbnail_url": "https://..."
    },
    "statistics": {
      "views": 1245,
      "downloads": 89,
      "shares": 23,
      "avg_view_time": "2m 34s"
    },
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "published_at": "2024-01-16T09:00:00Z"
  }
}

Update Album
http

PATCH /albums/{album_id}
Content-Type: application/json

{
  "title": "Updated Album Title",
  "description": "Updated description",
  "status": "published",
  "settings": {
    "download": {
      "enabled": false
    }
  }
}

Delete Album
http

DELETE /albums/{album_id}

Response: 204 No Content
Publish Album
http

POST /albums/{album_id}/publish

Response:
json

{
  "data": {
    "id": "album_123",
    "status": "published",
    "published_at": "2024-01-15T10:35:00Z",
    "public_url": "https://photolisting.dev/johns-photography/album_123"
  }
}

🖼️ Photos
List Photos in Album
http

GET /albums/{album_id}/photos

Query Parameters:
Parameter	Type	Description
sort_by	string	Sort field (created_at, filename, size)
sort_order	string	Sort order (asc, desc)
page	integer	Page number
limit	integer	Items per page

Response:
json

{
  "data": [
    {
      "id": "photo_123",
      "filename": "DSC_1234.jpg",
      "url": "https://...",
      "thumbnail_url": "https://...",
      "web_url": "https://...",
      "width": 1920,
      "height": 1280,
      "size_bytes": 2456789,
      "mime_type": "image/jpeg",
      "exif": {
        "camera": "Canon EOS R5",
        "lens": "RF 24-70mm f/2.8",
        "iso": 100,
        "shutter_speed": "1/250",
        "aperture": "f/2.8",
        "focal_length": "50mm"
      },
      "metadata": {
        "title": "Bride and Groom",
        "description": "First look moment",
        "tags": ["wedding", "portrait", "emotional"],
        "rating": 5
      },
      "created_at": "2024-01-15T10:30:00Z",
      "_links": {
        "self": "/photos/photo_123",
        "download": "/photos/photo_123/download",
        "edit": "/photos/photo_123"
      }
    }
  ]
}

Upload Photo(s)
http

POST /albums/{album_id}/photos
Content-Type: multipart/form-data

album_id: album_123
photos: [file1.jpg, file2.jpg, file3.jpg]

Alternative: Direct Upload
http

POST /photos/upload
Content-Type: application/json

{
  "album_id": "album_123",
  "uploads": [
    {
      "filename": "photo1.jpg",
      "size": 2456789,
      "mime_type": "image/jpeg",
      "upload_url": "https://storage.com/presigned-url"
    }
  ]
}

Response:
json

{
  "data": {
    "batch_id": "batch_789",
    "status": "processing",
    "total_files": 3,
    "processed": 0,
    "upload_urls": [
      {
        "filename": "photo1.jpg",
        "upload_url": "https://storage.com/presigned-url",
        "fields": {
          "key": "tenants/tenant_123/originals/photo_456.jpg"
        }
      }
    ]
  }
}

Get Photo Details
http

GET /photos/{photo_id}

Response:
json

{
  "data": {
    "id": "photo_123",
    "album_id": "album_456",
    "filename": "DSC_1234.jpg",
    "original_url": "https://...",
    "web_url": "https://...",
    "thumbnail_urls": {
      "small": "https://...",
      "medium": "https://...",
      "large": "https://..."
    },
    "dimensions": {
      "width": 1920,
      "height": 1280,
      "aspect_ratio": "3:2"
    },
    "file_info": {
      "size_bytes": 2456789,
      "mime_type": "image/jpeg",
      "format": "JPEG",
      "color_profile": "sRGB"
    },
    "exif": {
      "camera": "Canon EOS R5",
      "lens": "RF 24-70mm f/2.8",
      "settings": {
        "iso": 100,
        "shutter_speed": "1/250",
        "aperture": "f/2.8",
        "focal_length": "50mm"
      },
      "location": {
        "latitude": 40.7128,
        "longitude": -74.0060
      },
      "date_taken": "2024-01-15T10:30:00Z"
    },
    "metadata": {
      "title": "Bride and Groom",
      "description": "First look moment",
      "tags": ["wedding", "portrait", "emotional"],
      "rating": 5,
      "color_palette": ["#f3f4f6", "#1f2937", "#dc2626"]
    },
    "statistics": {
      "views": 456,
      "downloads": 23,
      "shares": 12
    },
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}

Update Photo Metadata
http

PATCH /photos/{photo_id}
Content-Type: application/json

{
  "metadata": {
    "title": "Updated Title",
    "description": "Updated description",
    "tags": ["updated", "tags"],
    "rating": 4
  }
}

Delete Photo
http

DELETE /photos/{photo_id}

Response: 204 No Content
Get Photo Download URL
http

GET /photos/{photo_id}/download

Query Parameters:
Parameter	Type	Description
size	string	Download size (original, web, print)
watermark	boolean	Include watermark (default: true)
expires_in	integer	URL expiry in seconds (default: 3600)

Response:
json

{
  "data": {
    "url": "https://storage.com/signed-url?signature=...",
    "expires_at": "2024-01-15T11:30:00Z",
    "size": "original",
    "watermark": true
  }
}

👥 Users & Permissions
List Users
http

GET /users

Response:
json

{
  "data": [
    {
      "id": "user_123",
      "email": "john@example.com",
      "name": "John Photographer",
      "role": "owner",
      "avatar_url": "https://...",
      "last_login_at": "2024-01-15T09:30:00Z",
      "created_at": "2024-01-01T00:00:00Z",
      "permissions": [
        "albums:create",
        "albums:edit",
        "photos:upload",
        "users:manage"
      ]
    }
  ]
}

Invite User
http

POST /users/invite
Content-Type: application/json

{
  "email": "new@example.com",
  "role": "editor",
  "permissions": ["albums:create", "photos:upload"],
  "message": "Welcome to our team!"
}

Response:
json

{
  "data": {
    "invitation_id": "invite_123",
    "email": "new@example.com",
    "role": "editor",
    "status": "pending",
    "expires_at": "2024-01-22T10:30:00Z",
    "invitation_url": "https://photolisting.dev/invite/invite_123"
  }
}

Update User Role
http

PATCH /users/{user_id}/role
Content-Type: application/json

{
  "role": "admin",
  "permissions": ["albums:create", "users:manage"]
}

Remove User
http

DELETE /users/{user_id}

🔗 Sharing & Client Access
Create Share Link
http

POST /albums/{album_id}/share
Content-Type: application/json

{
  "type": "proofing", // or "download", "view"
  "expires_in": "7d", // or "24h", "30m", "never"
  "password": "optional-password",
  "permissions": {
    "download": true,
    "share": false,
    "comment": true,
    "select": true
  },
  "watermark": {
    "enabled": true,
    "text": "Proof - Do Not Share",
    "opacity": 0.5
  },
  "notifications": {
    "on_view": true,
    "on_download": true,
    "on_selection": true
  }
}

Response:
json

{
  "data": {
    "token": "share_abc123def456",
    "url": "https://photolisting.dev/s/share_abc123def456",
    "type": "proofing",
    "expires_at": "2024-01-22T10:30:00Z",
    "permissions": {
      "download": true,
      "share": false,
      "comment": true,
      "select": true
    },
    "statistics": {
      "views": 0,
      "downloads": 0,
      "selections": 0
    }
  }
}

Get Share Info
http

GET /shares/{share_token}

Response:
json

{
  "data": {
    "token": "share_abc123def456",
    "album": {
      "id": "album_123",
      "title": "Summer Wedding",
      "description": "Beautiful outdoor wedding ceremony",
      "photo_count": 45
    },
    "type": "proofing",
    "expires_at": "2024-01-22T10:30:00Z",
    "permissions": {
      "download": true,
      "share": false,
      "comment": true,
      "select": true
    },
    "watermark": {
      "enabled": true,
      "text": "Proof - Do Not Share"
    }
  }
}

Submit Proofing Selections
http

POST /shares/{share_token}/selections
Content-Type: application/json

{
  "selections": [
    {
      "photo_id": "photo_123",
      "rating": 5,
      "comment": "Love this one!",
      "tags": ["favorite", "print"]
    }
  ],
  "client_info": {
    "name": "Jane Client",
    "email": "jane@example.com",
    "message": "These are my selections for editing"
  }
}

💳 Subscriptions & Billing
Get Subscription Details
http

GET /subscription

Response:
json

{
  "data": {
    "plan": "professional",
    "status": "active",
    "current_period_start": "2024-01-01T00:00:00Z",
    "current_period_end": "2024-02-01T00:00:00Z",
    "cancel_at_period_end": false,
    "limits": {
      "storage_gb": 50,
      "albums": null,
      "users": 5,
      "bandwidth_gb": 100
    },
    "usage": {
      "storage_gb": 12.5,
      "albums_created": 8,
      "active_users": 3,
      "bandwidth_this_month_gb": 24.7
    },
    "payment_method": {
      "type": "card",
      "last4": "4242",
      "exp_month": 12,
      "exp_year": 2025
    },
    "next_payment": {
      "amount": 2900, // in cents
      "currency": "usd",
      "date": "2024-02-01T00:00:00Z"
    }
  }
}

Update Subscription
http

POST /subscription/upgrade
Content-Type: application/json

{
  "plan": "studio",
  "payment_method_id": "pm_123456"
}

Cancel Subscription
http

POST /subscription/cancel

Response:
json

{
  "data": {
    "status": "canceled",
    "cancel_at_period_end": true,
    "current_period_end": "2024-02-01T00:00:00Z",
    "access_until": "2024-02-01T00:00:00Z"
  }
}

Get Invoices
http

GET /subscription/invoices

Response:
json

{
  "data": [
    {
      "id": "in_123456",
      "number": "INV-2024-001",
      "amount": 2900,
      "currency": "usd",
      "status": "paid",
      "period_start": "2024-01-01T00:00:00Z",
      "period_end": "2024-02-01T00:00:00Z",
      "paid_at": "2024-01-01T00:05:00Z",
      "pdf_url": "https://...",
      "line_items": [
        {
          "description": "Professional Plan",
          "amount": 2900,
          "quantity": 1
        }
      ]
    }
  ]
}

📊 Analytics
Get Dashboard Stats
http

GET /analytics/dashboard

Query Parameters:
Parameter	Type	Description
period	string	Time period (7d, 30d, 90d, year, all)
start_date	string	Custom start date (ISO 8601)
end_date	string	Custom end date (ISO 8601)

Response:
json

{
  "data": {
    "overview": {
      "total_albums": 12,
      "total_photos": 456,
      "total_storage_gb": 23.4,
      "active_clients": 8,
      "revenue_this_month": 2900
    },
    "engagement": {
      "total_views": 12456,
      "total_downloads": 345,
      "total_shares": 89,
      "avg_session_duration": "2m 34s",
      "bounce_rate": 0.32
    },
    "trends": {
      "views_last_30d": [
        {"date": "2024-01-01", "value": 123},
        {"date": "2024-01-02", "value": 145}
      ],
      "storage_growth": 12.5 // percent
    },
    "top_content": {
      "albums": [
        {"id": "album_123", "title": "Summer Wedding", "views": 1245},
        {"id": "album_456", "title": "Family Portraits", "views": 890}
      ],
      "photos": [
        {"id": "photo_123", "filename": "DSC_1234.jpg", "views": 456},
        {"id": "photo_456", "filename": "DSC_5678.jpg", "views": 389}
      ]
    },
    "clients": {
      "most_active": [
        {"email": "client1@example.com", "views": 234, "downloads": 12},
        {"email": "client2@example.com", "views": 189, "downloads": 8}
      ],
      "geography": [
        {"country": "US", "views": 4567},
        {"country": "UK", "views": 2345},
        {"country": "CA", "views": 1234}
      ]
    }
  }
}

Export Analytics
http

GET /analytics/export
Content-Type: application/json

{
  "format": "csv", // or "json", "xlsx"
  "data_types": ["engagement", "clients", "revenue"],
  "period": "90d",
  "include_metadata": true
}

Response: File download with appropriate Content-Type header.
🛠️ Batch Operations
Batch Upload Status
http

GET /batch/{batch_id}

Response:
json

{
  "data": {
    "id": "batch_123",
    "status": "completed",
    "total_files": 45,
    "processed": 45,
    "successful": 43,
    "failed": 2,
    "started_at": "2024-01-15T10:30:00Z",
    "completed_at": "2024-01-15T10:35:00Z",
    "results": [
      {
        "filename": "photo1.jpg",
        "status": "success",
        "photo_id": "photo_123",
        "url": "https://..."
      },
      {
        "filename": "photo2.jpg",
        "status": "failed",
        "error": "Invalid file format",
        "error_code": "INVALID_FORMAT"
      }
    ]
  }
}

Batch Update Metadata
http

POST /photos/batch/update
Content-Type: application/json

{
  "photo_ids": ["photo_123", "photo_456", "photo_789"],
  "metadata": {
    "tags": ["updated", "batch"],
    "rating": 4
  }
}

Batch Delete
http

POST /photos/batch/delete
Content-Type: application/json

{
  "photo_ids": ["photo_123", "photo_456", "photo_789"],
  "confirm": true
}

🔧 System Endpoints
Health Check
http

GET /health

Response:
json

{
  "status": "healthy",
  "version": "1.0.0",
  "timestamp": "2024-01-15T10:30:00Z",
  "services": {
    "database": "connected",
    "redis": "connected",
    "storage": "connected",
    "email": "connected"
  }
}

API Status
http

GET /status

Response:
json

{
  "data": {
    "status": "operational",
    "incidents": [],
    "maintenance": null,
    "uptime": 99.95,
    "response_time_ms": 124,
    "last_updated": "2024-01-15T10:30:00Z"
  }
}

🔐 Webhooks
Webhook Configuration
http

POST /webhooks
Content-Type: application/json

{
  "url": "https://your-server.com/webhooks",
  "events": [
    "album.created",
    "album.published",
    "photo.uploaded",
    "client.selection_submitted",
    "payment.succeeded"
  ],
  "secret": "your-webhook-secret",
  "enabled": true
}

Webhook Events Payload
json

{
  "id": "evt_123456",
  "type": "album.published",
  "timestamp": "2024-01-15T10:30:00Z",
  "tenant_id": "tenant_123",
  "data": {
    "album_id": "album_456",
    "title": "Summer Wedding",
    "published_at": "2024-01-15T10:30:00Z",
    "public_url": "https://photolisting.dev/..."
  }
}

🎮 SDKs & Client Libraries
Go Client
go

import "github.com/photolisting/sdk-go"

client := photolisting.NewClient("your-api-key")
album, err := client.Albums.Create(context.Background(), &photolisting.CreateAlbumRequest{
    Title: "My Album",
    Visibility: "private",
})

JavaScript/TypeScript Client
typescript

import { PhotoListingClient } from '@photolisting/sdk-js';

const client = new PhotoListingClient({
  apiKey: 'your-api-key',
});

const album = await client.albums.create({
  title: 'My Album',
  visibility: 'private',
});

Python Client
python

from photolisting import Client

client = Client(api_key="your-api-key")
album = client.albums.create(
    title="My Album",
    visibility="private"
)

🚀 Quick Examples
Complete Workflow Example
bash

# 1. Create an album
curl -X POST https://api.photolisting.dev/v1/albums \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title": "Client Gallery", "visibility": "private"}'

# 2. Upload photos
curl -X POST https://api.photolisting.dev/v1/albums/album_123/photos \
  -H "Authorization: Bearer $TOKEN" \
  -F "photos=@photo1.jpg" \
  -F "photos=@photo2.jpg"

# 3. Create share link for client
curl -X POST https://api.photolisting.dev/v1/albums/album_123/share \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"type": "proofing", "expires_in": "7d"}'

# 4. Get client selections
curl -X GET https://api.photolisting.dev/v1/albums/album_123/selections \
  -H "Authorization: Bearer $TOKEN"

Error Handling Example
go

package main

import (
    "fmt"
    "github.com/photolisting/sdk-go"
)

func main() {
    client := photolisting.NewClient("your-api-key")
    
    album, err := client.Albums.Create(context.Background(), &photolisting.CreateAlbumRequest{
        Title: "", // Empty title will cause validation error
    })
    
    if err != nil {
        if apiErr, ok := err.(*photolisting.APIError); ok {
            switch apiErr.Code {
            case "VALIDATION_ERROR":
                fmt.Printf("Validation error: %v\n", apiErr.Details)
            case "RATE_LIMIT_EXCEEDED":
                fmt.Printf("Rate limit exceeded, retry after: %v\n", apiErr.RetryAfter)
            default:
                fmt.Printf("API error: %v\n", apiErr.Message)
            }
        } else {
            fmt.Printf("Network error: %v\n", err)
        }
        return
    }
    
    fmt.Printf("Created album: %v\n", album.ID)
}

📚 Additional Resources

    Authentication Guide

    Rate Limiting

    Error Handling

    Webhook Guide

    SDK Documentation

🔄 API Versioning

The API is versioned in the URL path (/v1/). Breaking changes will result in a new version (/v2/).

Deprecation Policy:

    Endpoints marked as deprecated will continue to work for 6 months

    Deprecated endpoints will return a Deprecation header

    Migration guides will be provided for breaking changes

*Last Updated: 2024-01-15*
API Version: v1
For support: api-support@photolisting.dev