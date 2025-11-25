package command

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
	ctx, c := context.WithCancel(context.Background())

	return &CmdCtx{
		Context: ctx,
		cancel:  c,
		log:     logging.NewLogger(),
	}
}
func (ctx *CmdCtx) Cancel() {
	ctx.log.Info("ctx cancel called.")
	if ctx.cancel != nil {
		ctx.cancel()
	}
}

// CmdCtx app context
type CmdCtx struct {
	context.Context
	cancel context.CancelFunc
	flags  *pflag.FlagSet
	log    *logrus.Logger
}

func (ctx *CmdCtx) SetFlags(flags *pflag.FlagSet) {
	ctx.flags = flags
}

func (ctx *CmdCtx) GetFlags() *pflag.FlagSet {
	return ctx.flags
}

func (ctx *CmdCtx) Logger() *logrus.Logger {
	return ctx.log
}
