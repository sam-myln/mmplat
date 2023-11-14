package commands

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"mmplat/internal/logging"
)

// CobraRunECmd describes a function that can be used as a *cobra.Command RunE, PreRunE, or PostRunE.
type CobraRunECmd func(cmd *cobra.Command, args []string) (err error)

// NewCmdCtx returns a new CmdCtx.
func NewCmdCtx() *CmdCtx {
	ctx := context.Background()

	return &CmdCtx{
		Context: ctx,
		log:     logging.Logger(),
	}
}

// CmdCtx app context
type CmdCtx struct {
	context.Context

	flags *pflag.FlagSet
	log *logrus.Logger
}

func (ctx* CmdCtx) SetFlags(flags *pflag.FlagSet) {
	ctx.flags = flags
}

func (ctx* CmdCtx) GetFlags() *pflag.FlagSet {
	return ctx.flags
}

