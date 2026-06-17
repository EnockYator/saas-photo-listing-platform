// @title           SaaS Photo Listing Platform API
// @version         1.0
// @description     REST API for the SaaS Photo Listing Platform.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

package main

import (
	_ "github.com/EnockYator/saas-photo-listing-platform/internal/cli/api"
	_ "github.com/EnockYator/saas-photo-listing-platform/internal/cli/migrate"
	_ "github.com/EnockYator/saas-photo-listing-platform/internal/cli/tools"
	_ "github.com/EnockYator/saas-photo-listing-platform/internal/cli/worker"

	"github.com/EnockYator/saas-photo-listing-platform/internal/cli"
)

func main() {
	cli.Execute()
}
