package handlers

import (
	"fmt"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"mmplat/internal/filesystem"
	"mmplat/internal/util"
	"strings"
)

// AssetDir TODO here is not cool
var AssetDir string

type Handler struct {
	fs          *fs.FSWorker
	router      *router.Router
	routes      map[string]HandlerFn
	optional    map[string]HandlerFn
	middlewares map[string][]MiddlewareFn
}

type HandlerFn func(ctx *fasthttp.RequestCtx)
type MiddlewareFn func(ctx *fasthttp.RequestCtx) error

// NewHandler ctror
func NewHandler(fs *fs.FSWorker, router *router.Router) *Handler {
	return &Handler{fs, router, make(map[string]HandlerFn), make(map[string]HandlerFn), make(map[string][]MiddlewareFn)}
}

func (h *Handler) Fs() *fs.FSWorker {
	return h.fs
}

// Handle reroutes request to concrete handler
// appends registered middlewares
func (h *Handler) Handle(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())
	pos := strings.LastIndex(path, "/")
	var route = h.routes[path]
	/*if route == nil {
		route = h.optional[path[:pos+1]]
	}*/
	if route == nil {
		// try optional
		// /anton/anton/.js idx 25
		// /14 idx 0
		// look decrementary
		for pos > -1 {
			route = h.optional[path[:pos+1]]
			path = path[:pos+1]
			if route == nil {
				path = path[:pos]
				pos = strings.LastIndex(path, "/")
				path = path + "/"
			} else {
				pos = -1
			}
		}
	}
	if route != nil {
		if h.middlewares[path] != nil && len(h.middlewares[path]) > 0 {
			for _, v := range h.middlewares[path] {
				if err := v(ctx); err != nil {
					// Walk the middlewares stack, call each on context of request
					fmt.Fprintf(ctx, "%v", err)
					return
				}
			}
		}
		// Call actual hadnle
		route(ctx)
	}
}

// Register registers routes to actual handler
func (h *Handler) Register(r *router.Router, method, path string, fn HandlerFn, mfn ...MiddlewareFn) {
	if h.routes[path] != nil {
		panic(fmt.Sprintf("error route alredy bound to %s", path))
	}
	// register as optional
	// then check main, if empty, try optional
	pos := strings.Index(path, "{")
	if pos > -1 {
		h.optional[path[:pos]] = fn
		r.Handle(method, path, h.Handle)
		path = path[:pos]
	} else {
		h.routes[path] = fn
		r.Handle(method, path, h.Handle)
	}
	if len(mfn) > 0 {
		for _, middleware := range mfn {
			h.middlewares[path] = append(h.middlewares[path], middleware)
		}
	}
}

// Middleware "*" corresponds to all routes
func (h *Handler) Middleware(mfn MiddlewareFn, except []string, paths ...string) {
	if len(paths) == 1 && slices.Contains(paths, "*") {
		paths = nil
		paths = append(paths, append(maps.Keys(h.optional), maps.Keys(h.routes)...)...)
	}
outer:
	for _, path := range paths {
		if util.NotEmpty(except) {
			for _, exc := range except {
				pos := strings.Index(exc, "{")
				if pos > -1 {
					if util.Equals(path, exc[:pos]) {
						continue outer
					}
				} else {
					if util.Equals(path, exc) {
						continue outer
					}
				}
			}
		}
		h.middlewares[path] = append(h.middlewares[path], mfn)
	}
}

func (h *Handler) InternalServerError(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Internal server error")
}

func (h *Handler) NotFound(ctx *fasthttp.RequestCtx) {
	// TODO can create cutsome errorpage.qpl w/ title etc
	fmt.Fprintf(ctx, "Page %s wasn't found", append(ctx.Host(), ctx.RequestURI()...))
	ctx.SetStatusCode(fasthttp.StatusNotFound)
}

