package apperror

// =========================
// Generic / HTTP-aligned
// =========================

const (
	CodeInternal     = "INTERNAL_SERVER_ERROR"
	CodeBadRequest   = "BAD_REQUEST"
	CodeUnauthorized = "UNAUTHORIZED"
	CodeForbidden    = "FORBIDDEN"
	CodeNotFound     = "NOT_FOUND"
	CodeValidation   = "VALIDATION_ERROR"
	CodeUnavailable  = "SERVICE_UNAVAILABLE"
	CodeConflict     = "CONFLICT"
	CodeTooManyRequests = "TOO_MANY_REQUESTS"
)

// =========================
// Authentication / Authorization
// =========================

const (
	CodeAuthInvalidCredentials = "AUTH_INVALID_CREDENTIALS"
	CodeAuthTokenExpired       = "AUTH_TOKEN_EXPIRED"
	CodeAuthTokenInvalid       = "AUTH_TOKEN_INVALID"
	CodeAuthTokenMissing       = "AUTH_TOKEN_MISSING"
	CodeAuthSessionExpired     = "AUTH_SESSION_EXPIRED"
	CodeAuthAccountLocked      = "AUTH_ACCOUNT_LOCKED"
	CodeAuthMFARequired        = "AUTH_MFA_REQUIRED"
)

// =========================
// User / Account Domain
// =========================

const (
	CodeUserNotFound        = "USER_NOT_FOUND"
	CodeUserAlreadyExists   = "USER_ALREADY_EXISTS"
	CodeUserEmailInvalid    = "USER_EMAIL_INVALID"
	CodeUserProfileInvalid  = "USER_PROFILE_INVALID"
)

// =========================
// Validation (field-level / request-level)
// =========================

const (
	CodeValidationRequiredField = "VALIDATION_REQUIRED_FIELD"
	CodeValidationInvalidFormat = "VALIDATION_INVALID_FORMAT"
	CodeValidationOutOfRange    = "VALIDATION_OUT_OF_RANGE"
	CodeValidationFailed       = "VALIDATION_FAILED"
)

// =========================
// Database / Persistence
// =========================

const (
	CodeDBConnectionFailed = "DB_CONNECTION_FAILED"
	CodeDBQueryFailed      = "DB_QUERY_FAILED"
	CodeDBRecordNotFound   = "DB_RECORD_NOT_FOUND"
	CodeDBTransactionFail  = "DB_TRANSACTION_FAILED"
	CodeDBConstraintFail   = "DB_CONSTRAINT_VIOLATION"
)

// =========================
// External Services / Integrations
// =========================

const (
	CodeExternalServiceFailure = "EXTERNAL_SERVICE_FAILURE"
	CodeExternalTimeout        = "EXTERNAL_SERVICE_TIMEOUT"
	CodeExternalBadResponse    = "EXTERNAL_BAD_RESPONSE"
	CodeExternalUnavailable    = "EXTERNAL_SERVICE_UNAVAILABLE"
)

// =========================
// Payment / Billing (common SaaS need)
// =========================

const (
	CodePaymentFailed      = "PAYMENT_FAILED"
	CodePaymentDeclined    = "PAYMENT_DECLINED"
	CodePaymentRequired    = "PAYMENT_REQUIRED"
	CodeBillingNotFound    = "BILLING_NOT_FOUND"
	CodeSubscriptionExpired = "SUBSCRIPTION_EXPIRED"
	CodeSubscriptionInvalid = "SUBSCRIPTION_INVALID"
)

// =========================
// File
// =========================

const (
	CodeFileNotFound        = "FILE_NOT_FOUND"
	CodeFileUploadFailed    = "FILE_UPLOAD_FAILED"
	CodeFileTooLarge        = "FILE_TOO_LARGE"
	CodeFileInvalidType     = "FILE_INVALID_TYPE"
)

// =========================
// Storage
// =========================

const (
	CodeCloudStorageLimitReached  = "CLOUD_STORAGE_LIMIT_REACHED"
	CodeSloudStorageUnavailable  = "CLOUD_STORAGE_UNAVAILABLE"
)

// =========================
// Rate limiting / Abuse / Security
// =========================

const (
	CodeRateLimited        = "RATE_LIMIT_EXCEEDED"
	CodeIPBlocked          = "IP_BLOCKED"
	CodeSuspiciousActivity = "SUSPICIOUS_ACTIVITY"
	CodePermissionDenied   = "PERMISSION_DENIED"
)

// =========================
// System / Infrastructure
// =========================

const (
	CodeCircuitBreakerOpen = "CIRCUIT_BREAKER_OPEN"
	CodeDependencyFailure  = "DEPENDENCY_FAILURE"
	CodeTimeout            = "REQUEST_TIMEOUT"
	CodeShutdownInProgress = "SHUTDOWN_IN_PROGRESS"
)