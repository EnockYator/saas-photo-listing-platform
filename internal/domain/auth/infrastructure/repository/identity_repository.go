package repository

import (
	"context"

	"github.com/EnockYator/saas-photo-listing-platform/internal/domain/auth/domain"
	"github.com/google/uuid"
)

type IdentityRepository interface {
	Create(
		ctx context.Context,
		identity *domain.Identity,
	) error

	GetByIssuerAndSubject(
		ctx context.Context,
		issuer string,
		subject string,
	) (*domain.Identity, error)

	GetByUserID(
		ctx context.Context,
		userID uuid.UUID,
	) ([]domain.Identity, error)
}
