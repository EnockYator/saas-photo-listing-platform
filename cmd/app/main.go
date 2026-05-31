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
