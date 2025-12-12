package command

import (
	"context"
	"errors"
	"github.com/fasthttp/router"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
	filesystem "mmplat/internal/filesystem"
	handlers "mmplat/internal/handler"
	"mmplat/internal/middleware"
	"mmplat/internal/util"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func (ctx *CmdCtx) PreRunE(cmd *cobra.Command, _ []string) (err error) {
	var flag string
	defer func() {
		ctx.log.SetLevel(util.LogStringToLevel(flag))
	}()
	flag, err = cmd.PersistentFlags().GetString(cmdFlagNameLogLevel)
	if err != nil {
		return errors.New("missing flag, using default" + err.Error())
	}
	return
}

func (ctx *CmdCtx) AppRunE(_ *cobra.Command, args []string) error {

	defer ctx.Cancel()
	return ctx.run(args)
}

func (ctx *CmdCtx) run(args []string) error {
	wrappterCtx, cancel := context.WithCancel(ctx)
	cctx := context.WithValue(wrappterCtx, "logger", ctx.Logger())
	defer func() {
		ctx.log.Info("cctx cancel called.")
		cancel()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)
	var fs *filesystem.FSWorker
	{
		var err error
		var r bool
		r, err = ctx.flags.GetBool(cmdFlagNameRecurse)

		ctx.log.Info("provided Args(): " + strings.Join(args, ","))
		if len(args) == 0 {
			ctx.log.Error("Error, need to provide at least one directory")
			cancel()
			return errors.New("error, some required arguments are missing")
		}

		fmt, err := ctx.flags.GetStringSlice(cmdFlagNameFormat)
		if err != nil {
			ctx.log.Info("Flag " + cmdFlagNameFormat + " omitted, using default preset.")
		}
		fs, err = filesystem.NewFSWorker(r, args, fmt, ctx.log)
		if err != nil {
			ctx.log.Error("failed initializing FSWorker: " + err.Error())
			cancel()
			return err
		}
	}

	err := fs.BuildTree()
	if err != nil {
		ctx.log.Errorf("Process finished due to directory parsing error: %v", err)
		cancel()
		return err
	}
	r := router.New()
	handler := handlers.NewHandler(fs, r)
	r.NotFound = handler.NotFound
	handler.Register(r, "GET", "/", handler.Index, middleware.ContentTypeHtmlMiddleware)
	handler.Register(r, "POST", "/upload", handler.Upload, middleware.ContentTypeHtmlMiddlewareRedirect)
	handler.Register(r, "GET", "/favicon.ico", handler.FaviconPieceOfShit) // whatever TF that is
	handler.Register(r, "GET", "/{item:[0-9]{1,3}}", handler.Item, middleware.ContentTypeHtmlMiddleware)
	handler.Register(r, "GET", "/assets/{asset:*}", handler.Asset)
	// api(or inner)-handler for stream request
	handler.Register(r, "GET", "/stream/{name}", handler.Stream)

	port, _ := ctx.flags.GetString(cmdFlagNamePort)
	t, err := ctx.flags.GetBool(cmdFlagNameTest)
	var ln net.Listener
	if err == nil && t {
		ctx.log.Infof("binding to localhost")
		if strings.Contains(port, ":") {
			ctx.log.Errorf("Error, --test is incompatible with this port format")
			cancel()
			return errors.New("error, --test is incompatible with this port format")
		}
		ln, err = new(net.ListenConfig).Listen(cctx, "tcp", "localhost:"+port)
	} else {
		ctx.log.Infof("binding to iface: %v", err)

		if strings.Contains(port, ":") {
			ln, err = new(net.ListenConfig).Listen(cctx, "tcp", port)
		} else {
			ln, err = new(net.ListenConfig).Listen(cctx, "tcp", ":"+port)
		}
	}
	ctx.log.Infof("init done, listening on :%s", port)
	if err != nil {
		ctx.log.Errorf("listen error: %v", err)
	}

	if err != nil {
		cancel()
		return err
	}

	// workerpool:185 race condition?
	// racecondition occurs earlier, server.go:1786
	server := &fasthttp.Server{
		Handler:            r.Handler,
		MaxRequestBodySize: 100 << 20, // 100MB (100 * 1024 * 1024)
		ReadBufferSize:     1 << 20,   // 1MB (1024 * 1024)
		WriteBufferSize:    1 << 20,   // 1MB
	}
	go func() {
		err := func() error {
			return server.Serve(ln)
		}()
		if err != nil {
			ctx.log.Error("serve failed: " + err.Error())
		}
	}()
	ctx.log.Info("Startup complete")
	select {
	case s := <-quit:
		ctx.log.WithField("signal", s.String()).Info("Shutdown initiated due to process signal")
		_ = ln.Close()
		return cctx.Err()
	case <-cctx.Done():
		ctx.log.Info("Shutdown initiated due to context completion")
		_ = ln.Close()
		return nil
	}
}
