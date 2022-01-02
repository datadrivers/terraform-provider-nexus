//go:build tools
// +build tools

package tools

import (
	_ "github.com/client9/misspell/cmd/misspell"
	_ "github.com/datadrivers/terraform-plugin-docs/cmd/tfplugindocs"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
)
