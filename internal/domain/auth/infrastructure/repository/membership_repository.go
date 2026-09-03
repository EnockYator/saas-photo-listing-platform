package repository

import (
	"context"

	"github.com/EnockYator/saas-photo-listing-platform/internal/domain/auth/domain"
	"github.com/google/uuid"
)

type MembershipRepository interface {
	Create(
		ctx context.Context,
		membership *domain.Membership,
	) error

	Get(
		ctx context.Context,
		userID uuid.UUID,
		tenantID uuid.UUID,
	) (*domain.Membership, error)

	ListByUserID(
		ctx context.Context,
		userID uuid.UUID,
	) ([]domain.Membership, error)
}
