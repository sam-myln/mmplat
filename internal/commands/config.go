package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

func NewAppCmd() *cobra.Command {
	version := version
	ctx := NewCmdCtx()

	cmd := &cobra.Command{
		Use:               fmtCmdUse,
		Short:             fmt.Sprintf(fmtCmdShort, version),
		Long:              fmt.Sprintf(fmtCmdLong, version),
		Example:           fmtCmdEx,
		Version:           version,
		Args:              cobra.MinimumNArgs(1),
		PreRunE:           ctx.PreRunE,
		RunE:              ctx.AppRunE,
		SilenceUsage: true,
		DisableAutoGenTag: true,
	}

	cmd.PersistentFlags().StringP(
		cmdFlagNameLogLevel, "l", "error", "Log verbosity level: info, debug, error. Default error.",
	)

	cmd.PersistentFlags().StringP(
		cmdFlagNameFormat, "f", "formats", "Only include specific formats",
	)

	cmd.PersistentFlags().BoolP(cmdFlagNameRecurse, "r", false, "process directories recursively")

	cmd.PersistentFlags().StringP(
		cmdFlagNameTest, "t", "test", "Run in localhost",
	)

	cmd.PersistentFlags().StringP(
		cmdFlagNamePort, "p", "8080", "Define listen port, default: 8080",
	)

	ctx.flags = cmd.PersistentFlags()

	return cmd
}
