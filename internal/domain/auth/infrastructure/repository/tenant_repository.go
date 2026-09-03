package repository

import (
	"context"

	"github.com/EnockYator/saas-photo-listing-platform/internal/domain/auth/domain"
	"github.com/google/uuid"
)

type TenantRepository interface {
	Create(
		ctx context.Context,
		tenant *domain.Tenant,
	) error

	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (*domain.Tenant, error)
}
