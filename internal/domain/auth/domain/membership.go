package domain

import (
	"github.com/google/uuid"
	"time"
)

type Membership struct {
	UserID    uuid.UUID `json:"user_id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
