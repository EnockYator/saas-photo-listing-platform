// +build tools

package main

// This import is only for tracking the CLI tool version.
// The blank identifier (_) prevents it from being included in the final binary.
import (
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc" // SQL code generation
	_ "github.com/golang-migrate/migrate/v4/cmd/migrate" // Database migrations
	_ "github.com/air-verse/air" // Live reload during development
	_ "github.com/go-swagger/go-swagger/cmd/swagger" // Swagger API documentation generation
	_ "github.com/swaggo/swag/cmd/swag" // Swagger CLI tool
)


