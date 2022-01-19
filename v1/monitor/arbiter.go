package monitor

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

type SetupCallback func(command *cobra.Command) error

type arbiter struct{}

func (e *arbiter) initialize(options *ArbiterOptions) error {
	_ = options

	return nil
}

func (e *arbiter) run(options *RunOptions) error {
	return nil
}

func (e *arbiter) setupLogging() error {
	return nil
}

func (e *arbiter) setupSignals() error {
	return nil
}

// ArbiterOption type for updating the options struct
type ArbiterOption func(o *ArbiterOptions)

type ArbiterOptions struct {
	// Name of the tool
	Name string

	// Description is the long explanation of what the command is supposed to do.
	Description string

	// ShortDescription is the short blurb about functionality.
	ShortDescription string

	// Value to use for the primary arbiter
	// Default: 'run'
	Command string

	// InitializeCallback will be called after the command line system has been
	// initialized and act as a setup hook for the application.
	InitializeCallback func()

	// RunSetup callback to set up the default run command
	RunSetup SetupCallback

	// VersionSetup Handler for setting up teh version command.
	VersionSetup SetupCallback
}

func (o *ArbiterOptions) Validate() error {
	return nil
}

func Execute(opts ...ArbiterOption) {
	options := new(ArbiterOptions)
	for _, opt := range opts {
		opt(options)
	}

	if err := options.Validate(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr,
			"Failed to validat ethe arbiter options given: %s\n", err)
		os.Exit(1)
	}

	exe := new(arbiter)
	if err := exe.initialize(options); err != nil {
		_, _ = fmt.Fprintf(os.Stderr,
			"Failed to initialize the arbiter: %s\n", err)
		os.Exit(1)
	}

	parser := &cobra.Command{
		Use:   options.Name,
		Long:  options.Description,
		Short: options.ShortDescription,
	}

	setupVersionCommand(parser, options.VersionSetup)
	setupRunCommand(parser, options)

	cobra.OnInitialize(options.InitializeCallback)
}

type RunOptions struct {
	daemon   bool
	interval int
	logfile  string
	loglevel string
	stdout   bool
	pidfile  string
	uid      string
	gid      string
	config   string
}

func setupRunCommand(cmd *cobra.Command, options *ArbiterOptions) {
	runCmd := &cobra.Command{
		Use:   options.Name,
		Long:  options.Description,
		Short: options.ShortDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			var values RunOptions

			cmd.Flags().BoolVarP(
				&values.daemon,
				"daemon",
				"d",
				false,
				"Daemonize the application (only supported on Linux)")

			cmd.Flags().BoolVarP(
				&values.stdout,
				"stdout",
				"o",
				false,
				"Log messages will be sent to STDOUT")

			cmd.Flags().IntVarP(
				&values.interval,
				"interval",
				"i",
				10,
				"Interval at which the runtime callback should be invoked.")

			cmd.Flags().StringVarP(
				&values.logfile,
				"logfile",
				"l",
				"",
				"Log level for the logging system.")

			cmd.Flags().StringVarP(
				&values.loglevel,
				"loglevel",
				"l",
				"INFO",
				"Default log level for logging.")

			return nil
		},
	}

	if err := options.RunSetup(runCmd); err != nil {
		_, _ = fmt.Fprintf(os.Stderr,
			"Setup command failed: %s\n", err)
		os.Exit(1)
	}

	cmd.AddCommand(runCmd)
}

func setupVersionCommand(cmd *cobra.Command, setup SetupCallback) {
	versionCmd := &cobra.Command{
		Use:   "version",
		Long:  "Retrieve the version information for this tool",
		Short: "Retrieve the version information for this tool",
	}

	if err := setup(versionCmd); err != nil {
		_, _ = fmt.Fprintf(os.Stderr,
			"Failed to setup the version command: %s\n", err)
		os.Exit(1)
	}

	cmd.AddCommand(versionCmd)
}
