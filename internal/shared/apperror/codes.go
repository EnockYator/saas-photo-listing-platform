package apperror

// ErrorCode identifies a stable application-level error condition.
type ErrorCode string

// =========================
// Generic / HTTP-aligned
// =========================

const (
	CodeInternalServerError ErrorCode = "INTERNAL_SERVER_ERROR" // 500
	CodeBadRequest          ErrorCode = "BAD_REQUEST"           // 400
	CodeMethodNotAllowed    ErrorCode = "METHOD_NOT_ALLOWED"    // 405
	CodeUnauthorized        ErrorCode = "UNAUTHORIZED"          // 401
	CodeForbidden           ErrorCode = "FORBIDDEN"             // 403
	CodeNotFound            ErrorCode = "NOT_FOUND"             // 404
	CodeValidation          ErrorCode = "VALIDATION_ERROR"      // 422
	CodeServiceUnavailable  ErrorCode = "SERVICE_UNAVAILABLE"   // 503
	CodeConflict            ErrorCode = "CONFLICT"              // 409
	CodeTooManyRequests     ErrorCode = "TOO_MANY_REQUESTS"     // 429
)

// =========================
// Authentication / Authorization
// =========================

const (
	CodeAuthInvalidCredentials ErrorCode = "AUTH_INVALID_CREDENTIALS" // 401
	CodeAuthTokenExpired       ErrorCode = "AUTH_TOKEN_EXPIRED"       // 401
	CodeAuthTokenInvalid       ErrorCode = "AUTH_TOKEN_INVALID"       // 401
	CodeAuthTokenMissing       ErrorCode = "AUTH_TOKEN_MISSING"       // 401
	CodeAuthSessionExpired     ErrorCode = "AUTH_SESSION_EXPIRED"     // 401
	CodeAuthAccountLocked      ErrorCode = "AUTH_ACCOUNT_LOCKED"      // 423
	CodeAuthMFARequired        ErrorCode = "AUTH_MFA_REQUIRED"        // 401
)

// =========================
// User / Account Domain
// =========================

const (
	CodeUserNotFound       ErrorCode = "USER_NOT_FOUND"       // 404
	CodeUserAlreadyExists  ErrorCode = "USER_ALREADY_EXISTS"  // 409
	CodeUserEmailInvalid   ErrorCode = "USER_EMAIL_INVALID"   // 422
	CodeUserProfileInvalid ErrorCode = "USER_PROFILE_INVALID" // 422
)

// =========================
// Input Validation
// =========================

const (
	CodeValidationRequiredField ErrorCode = "VALIDATION_REQUIRED_FIELD" // 422
	CodeValidationInvalidEmailFormat ErrorCode = "VALIDATION_INVALID_EMAIL_FORMAT" // 422
	CodeValidationOutOfRange    ErrorCode = "VALIDATION_OUT_OF_RANGE"   // 422
	CodeValidationFailed        ErrorCode = "VALIDATION_FAILED"         // 422
)

// =========================
// Database / Persistence
// =========================

const (
	CodeDBConnectionFailed ErrorCode = "DB_CONNECTION_FAILED"  // 503
	CodeDBQueryFailed      ErrorCode = "DB_QUERY_FAILED"       // 500
	CodeDBRecordNotFound   ErrorCode = "DB_RECORD_NOT_FOUND"   // 404
	CodeDBTransactionFail  ErrorCode = "DB_TRANSACTION_FAILED" // 500
	CodeDBConstraintFail   ErrorCode = "DB_CONSTRAINT_VIOLATION" // 409
)

// =========================
// External Services / Integrations
// =========================

const (
	CodeExternalServiceFailure ErrorCode = "EXTERNAL_SERVICE_FAILURE"    // 502
	CodeExternalTimeout        ErrorCode = "EXTERNAL_SERVICE_TIMEOUT"     // 504
	CodeExternalBadResponse    ErrorCode = "EXTERNAL_BAD_RESPONSE"        // 502
	CodeExternalUnavailable    ErrorCode = "EXTERNAL_SERVICE_UNAVAILABLE" // 503
)

// =========================
// Payment / Billing
// =========================

const (
	CodePaymentFailed       ErrorCode = "PAYMENT_FAILED"        // 402
	CodePaymentDeclined     ErrorCode = "PAYMENT_DECLINED"      // 402
	CodePaymentRequired     ErrorCode = "PAYMENT_REQUIRED"      // 402
	CodeBillingNotFound     ErrorCode = "BILLING_NOT_FOUND"     // 404
	CodeSubscriptionExpired ErrorCode = "SUBSCRIPTION_EXPIRED"  // 402
	CodeSubscriptionInvalid ErrorCode = "SUBSCRIPTION_INVALID"  // 422
)

// =========================
// File
// =========================

const (
	CodeFileNotFound     ErrorCode = "FILE_NOT_FOUND"     // 404
	CodeFileUploadFailed ErrorCode = "FILE_UPLOAD_FAILED" // 500
	CodeFileTooLarge     ErrorCode = "FILE_TOO_LARGE"     // 413
	CodeFileInvalidType  ErrorCode = "FILE_INVALID_TYPE"  // 415
)

// =========================
// Storage
// =========================

const (
	CodeCloudStorageLimitReached ErrorCode = "CLOUD_STORAGE_LIMIT_REACHED" // 507
	CodeCloudStorageUnavailable  ErrorCode = "CLOUD_STORAGE_UNAVAILABLE"   // 503
)

// =========================
// Rate limiting / Abuse / Security
// =========================

const (
	CodeRateLimited        ErrorCode = "RATE_LIMIT_EXCEEDED"  // 429
	CodeIPBlocked          ErrorCode = "IP_BLOCKED"           // 403
	CodeSuspiciousActivity ErrorCode = "SUSPICIOUS_ACTIVITY"  // 403
	CodePermissionDenied   ErrorCode = "PERMISSION_DENIED"    // 403
)

// =========================
// System / Infrastructure
// =========================

const (
	CodeCircuitBreakerOpen ErrorCode = "CIRCUIT_BREAKER_OPEN"  // 503
	CodeDependencyFailure  ErrorCode = "DEPENDENCY_FAILURE"    // 503
	CodeTimeout             ErrorCode = "REQUEST_TIMEOUT"        // 504
	CodeShutdownInProgress  ErrorCode = "SHUTDOWN_IN_PROGRESS"  // 503
)