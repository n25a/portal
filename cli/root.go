// Package cli is the commands collector for the CLI application.
package cli

// CLI is the main struct for the CLI application.
var CLI struct {
	Export exportCmd `cmd:"" help:"Export result."`
}
