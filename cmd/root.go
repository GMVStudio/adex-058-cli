package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/gmvstudio/adex-cli/errs"
	"github.com/gmvstudio/adex-cli/internal/client"
	"github.com/gmvstudio/adex-cli/internal/config"
	"github.com/spf13/cobra"
)

// Factory holds shared dependencies injected into commands.
type Factory struct {
	Config *config.Config
	Client *client.Client
	Out    io.Writer
	ErrOut io.Writer
}

var globalFactory *Factory

const rootLong = `adex — ADEX CLI tool.

USAGE:
    adex <command> [subcommand] [options]

EXAMPLES:
    # Query campaign daily report
    adex raw campaign daily --tenant 6 --campaign C-618-001-619 --range 1d

    # With pretty output
    adex raw campaign daily --tenant 6 --range 1d --format pretty

    # With table output
    adex raw campaign daily --tenant 6 --range 1d --format table

ENVIRONMENT:
    ADEX_API_BASE_URL  API base URL (default: http://localhost:8000)`

// NewRootCmd creates the root cobra command.
func NewRootCmd(f *Factory) *cobra.Command {
	root := &cobra.Command{
		Use:           "adex",
		Short:         "ADEX CLI tool",
		Long:          rootLong,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.PersistentFlags().String("base-url", "", "API base URL (overrides ADEX_API_BASE_URL env)")
	root.PersistentFlags().String("format", "json", "output format: json (default) | pretty | table")
	root.PersistentFlags().Bool("dry-run", false, "print request without executing")

	root.AddCommand(newInitCmd(f))
	root.AddCommand(newRawCmd(f))
	root.AddCommand(newKsCmd(f))
	root.AddCommand(newOeCmd(f))
	root.AddCommand(newTenantCmd(f))
	root.AddCommand(newUserCmd(f))
	root.AddCommand(newSkillCmd(f))

	return root
}

// Execute runs the root command and returns the process exit code.
func Execute() int {
	cfg := config.Load()
	f := &Factory{
		Config: cfg,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
	globalFactory = f

	root := NewRootCmd(f)

	err := root.Execute()
	if err != nil {
		return handleError(f, err)
	}
	return 0
}

// handleError dispatches a command error to the typed envelope writer
// and returns the process exit code.
func handleError(f *Factory, err error) int {
	if errs.IsTyped(err) {
		fmt.Fprintln(f.ErrOut, string(errs.RenderEnvelope(err)))
		prob, _ := errs.ProblemOf(err)
		return errs.ExitCodeForCategory(prob.Category)
	}

	// Wrap cobra usage errors (missing required flags, unknown commands, etc.)
	// as typed validation errors so they produce a parseable JSON envelope.
	wrapped := errs.NewValidationError(errs.SubtypeInvalidArgument, "%v", err).WithCause(err)
	fmt.Fprintln(f.ErrOut, string(errs.RenderEnvelope(wrapped)))
	return errs.ExitCodeForCategory(errs.CategoryValidation)
}

// resolveClient lazily creates an API client, applying --base-url override.
func (f *Factory) resolveClient(cmd *cobra.Command) *client.Client {
	if f.Client != nil {
		return f.Client
	}
	baseURL := f.Config.BaseURL
	if v, _ := cmd.Flags().GetString("base-url"); v != "" {
		baseURL = v
	}
	f.Client = client.New(baseURL,
		client.WithErrOut(f.ErrOut),
		client.WithAPIKey(f.Config.Authorization),
	)
	return f.Client
}

// resolveFormat returns the output format from --format flag.
func (f *Factory) resolveFormat(cmd *cobra.Command) string {
	format, _ := cmd.Flags().GetString("format")
	return format
}
