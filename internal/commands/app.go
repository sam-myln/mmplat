package commands

import (
	"context"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/fasthttp/session/v2"
	"github.com/fasthttp/session/v2/providers/memory"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
	filesystem "mmplat/internal/filesystem"
	"mmplat/internal/handlers"
	"mmplat/internal/middleware"
	"mmplat/internal/util"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// NewAppCmd THE entry point of app
func NewAppCmd() *cobra.Command {
	version := version
	ctx := NewCmdCtx()

	cmd := &cobra.Command{
		Use:               "mmplat",
		Short:             fmt.Sprintf(fmtCmdShort, version),
		Long:              fmt.Sprintf(fmtCmdLong, version),
		Example:           fmtCmdEx,
		Version:           version,
		Args:              cobra.NoArgs,
		PreRunE:           ctx.PreRunE,
		RunE:              ctx.AppRunE,
		DisableAutoGenTag: true,
	}
	cmd.PersistentFlags().BoolP(
		cmdFlagNameRecursive, cmdFlagNameRecursiveShort,
		true, "Indicates deep traverse of provided paths",
	)

	cmd.PersistentFlags().String(
		cmdFlagNameLogLevel, cmdFlagNameLogLevelDefault, "Log verbosity level: info, debug, error",
	)

	cmd.PersistentFlags().StringSlice(
		cmdFlagNameFolders, nil, "List of directories to search",
	)
	cmd.MarkFlagRequired(cmdFlagNameFolders)
	cmd.PersistentFlags().StringSlice(
		cmdFlagNameFormats, []string{`*`}, //[]string{cmdFlagNamePredefinedAudio, cmdFlagNamePredefinedVideo},
		"List of formats to include with search.",
	)

	placeHold := strings.Join([]string{cmdFlagNameAddrAddress,
		strings.Join([]string{"--", cmdFlagNameBindAddress}, "")}, ",")
	cmd.PersistentFlags().String(
		placeHold, "", "Listen, bind address to, default localhost:8080",
	)
	cmd.PersistentFlags().String(
		cmdFlagNameAddrAddress, cmdFlagNameBindAddressDefault, "Listen, bind address to, default localhost:8080",
	)
	cmd.PersistentFlags().String(
		cmdFlagNameBindAddress, cmdFlagNameBindAddressDefault, "Listen, bind address to, default localhost:8080",
	)
	cmd.PersistentFlags().MarkHidden(cmdFlagNameAddrAddress)
	cmd.PersistentFlags().MarkHidden(cmdFlagNameBindAddress)

	cmd.PersistentFlags().Bool(
		cmdFlagNameAllowGuest, false, "allow guest",
	)

	cmd.PersistentFlags().StringSlice(
		cmdFlagNameLoginData, []string{} /*[]string{"login-data"}*/, "Login data providers. File, w/ 2 entries: login+pass. One per user",
	)

	cmd.MarkFlagsMutuallyExclusive(cmdFlagNameAllowGuest, cmdFlagNameLoginData)

	ctx.flags = cmd.PersistentFlags()

	// TODO yeah, here too
	//  const or global or config or etc
	handlers.AssetDir, _ = os.Getwd()
	handlers.AssetDir += "/assets"

	cmd.AddCommand(
		newOtpCmd(ctx),
	)

	return cmd
}

func (ctx *CmdCtx) PreRunE(cmd *cobra.Command, _ []string) error {
	ctx.log.SetLevel(
		util.LogStringToLevel(
			cmd.PersistentFlags().GetString(cmdFlagNameLogLevel),
		),
	)
	return nil
}

func (ctx *CmdCtx) AppRunE(_ *cobra.Command, _ []string) error {
	ctx.Run()
	return nil
}

func (ctx *CmdCtx) Run() {
	cctx, cancel := context.WithCancel(ctx)
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)
	var fs *filesystem.FSWorker
	{
		a, _ := ctx.flags.GetBool(cmdFlagNameRecursive)
		b, _ := ctx.flags.GetStringSlice(cmdFlagNameFolders)
		c, _ := ctx.flags.GetStringSlice(cmdFlagNameFormats)
		fs = filesystem.NewFSWorker(a, b, c)
	}

	err := fs.BuildTree()
	if err != nil {
		ctx.log.Debugf("Processo finished due to directory parsing error: %v", err)
	}
	r := router.New()
	handler := handlers.NewHandler(fs, r)
	provider, err := memory.New(memory.Config{})
	if err != nil {
		ctx.log.Errorf("memory provider error %v :", err)
	}

	var credMngr *util.CredMngr
	{
		a, _ := ctx.flags.GetBool(cmdFlagNameAllowGuest)
		t, _ := ctx.flags.GetStringSlice(cmdFlagNameLoginData)
		if util.Not((util.NotEmpty(t) || a) && !(util.NotEmpty(t) && a)) {
			panic(fmt.Sprintf("%s, %v=%v %v=%v", "error: mutually exclusive conditions",
				cmdFlagNameAllowGuest, a, cmdFlagNameLoginData, t))
		}
		credMngr = util.CreateCredMngr(t...)
		handlers.AllowGuest = a
	}
	credMngr.ParseCredentials()
	auth := handlers.CreateAuth(session.New(
		session.NewDefaultConfig(),
	), credMngr)
	_ = auth.Session().SetProvider(provider)

	r.NotFound = handler.NotFound
	handler.Register(r, "GET", "/", handler.Index, middleware.ContentTypeHtmlMiddleware)
	handler.Register(r, "GET", "/favicon.ico", handler.FaviconPieceOfShit) // whatever TF that is
	handler.Register(r, "GET", "/{item:[0-9]{1,3}}", handler.Item, middleware.ContentTypeHtmlMiddleware)
	handler.Register(r, "GET", "/assets/{asset:*}", handler.Asset)
	// api(or inner)-handler for stream request
	handler.Register(r, "GET", "/stream/{name}", handler.Stream)
	// auth
	if util.Not(handlers.AllowGuest) {
		handler.Register(r, "GET", "/login", auth.Login, middleware.ContentTypeHtmlMiddleware)
		handler.Register(r, "GET", "/totp", auth.Totp, middleware.ContentTypeHtmlMiddleware)
		handler.Register(r, "POST", "/auth", auth.Auth, middleware.ContentTypeHtmlMiddleware)
		handler.Register(r, "POST", "/verify", auth.VerifyTotp, middleware.ContentTypeHtmlMiddleware)

		handler.Middleware(auth.AuthMiddleware, []string{
			"/auth", "/login", "/assets/{asset:*}", "/totp", "/verify", "/favicon.ico",
		}, "*")
	}
	userSession := &middleware.UserSession{Session: auth.Session()}
	handler.Middleware(userSession.Csrf, nil,
		"/login", "/totp",
	)
	addr, _ := ctx.flags.GetString(cmdFlagNameAddrAddress)
	bind, _ := ctx.flags.GetString(cmdFlagNameBindAddress)
	var ln net.Listener
	if addr != cmdFlagNameBindAddressDefault {
		ln, err = new(net.ListenConfig).Listen(cctx, "tcp4", addr)
	} else if bind != cmdFlagNameBindAddressDefault {
		ln, err = new(net.ListenConfig).Listen(cctx, "tcp4", bind)
	} else {
		ln, err = new(net.ListenConfig).Listen(cctx, "tcp4", cmdFlagNameBindAddressDefault)
	}
	ctx.log.Infof("init done, listening on: %s", addr)
	if err != nil {
		ctx.log.Errorf("listen error: %v", err)
	}

	if err != nil {
		os.Exit(1)
	}

	// workerpool:185 race condition?
	// racecondition occurs earlier, server.go:1786
	go func() error {
		return fasthttp.Serve(ln, r.Handler)
	}()
	ctx.log.Debug("Startup complete")
	select {
	case s := <-quit:
		ctx.log.WithField("signal", s.String()).Debug("Shutdown initiated due to process signal")
		_ = ln.Close()
	case <-cctx.Done():
		ctx.log.Debug("Shutdown initiated due to context completion")
		_ = ln.Close()
	}
}
