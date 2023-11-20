//go:build tools
// +build tools

package tools

import (
	_ "github.com/onsi/ginkgo/v2/ginkgo"
	_ "golang.org/x/tools/cmd/goimports"
)
