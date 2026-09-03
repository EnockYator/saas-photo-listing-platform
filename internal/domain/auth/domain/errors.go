package domain

import (
	"github.com/EnockYator/saas-photo-listing-platform/internal/shared/apperror"
)

var (
	ErrUserNotFound          apperror.ErrorCode = "USER_NOT_FOUND"
	ErrIdentityNotFound      apperror.ErrorCode = "IDENTITY_NOT_FOUND"
	ErrIdentityAlreadyLinked apperror.ErrorCode = "IDENTITY_ALREADY_LINKED"
	ErrMembershipNotFound    apperror.ErrorCode = "MEMBERSHIP_NOT_FOUND"
	ErrInvalidRole           apperror.ErrorCode = "INVALID_ROLE"
	ErrMissingIDToken        apperror.ErrorCode = "MISSING_TOKEN_ID"
)
