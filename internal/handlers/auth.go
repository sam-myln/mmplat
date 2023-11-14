package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fasthttp/session/v2"
	"github.com/valyala/fasthttp"
	"mmplat/internal/templates"
	"mmplat/internal/util"
	"reflect"
)

var AllowGuest bool

type Auth struct {
	session  *session.Session
	credMngr *util.CredMngr
}

func CreateAuth(session *session.Session,
	credMngr *util.CredMngr) *Auth {
	return &Auth{session, credMngr}
}

// AuthMiddleware final stage verify
func (h *Auth) AuthMiddleware(ctx *fasthttp.RequestCtx) error {
	storage, _ := h.Session().Get(ctx)
	if util.Empty(storage.Get("UserAuth")) {
		strLocation := []byte("Location")
		ctx.Response.Header.SetCanonical(strLocation, []byte("/login"))
		ctx.Response.SetStatusCode(fasthttp.StatusFound)
		return errors.New("you must be authenticatied in order to access this resource")
	}
	return nil
}

// userExists verifies that user exists and credentials are correct
func (h *Auth) userExists(login, pass []byte) bool {
	strLogin := string(login)
	strPass := string(pass)

	if h.credMngr.Exists(strLogin, strPass) {
		return true
	}
	return false
}

// VerifyTotp 2nd factor
func (h *Auth) VerifyTotp(ctx *fasthttp.RequestCtx) {
	if AllowGuest {
		return
	}
	storage, _ := h.Session().Get(ctx)
	defer h.Session().Save(ctx, storage)
	passcode := ctx.Request.PostArgs().Peek("code")
	csrf := ctx.Request.PostArgs().Peek("csrf")
	loginPass := storage.Get("LoginPass")
	if util.Empty(passcode ) || util.Empty(csrf) || util.Empty(loginPass) {
		strLocation := []byte("Location")
		// TODO yeah, i know. don't mention
		ctx.Response.Header.Set("AuthFailed", "wrong credentials")
		ctx.Response.Header.SetCanonical(strLocation, []byte("/login"))
		ctx.Response.SetStatusCode(fasthttp.StatusFound)
	}
	valid, err := util.ValidateKeyPass(loginPass.(string), string(passcode))
	if err != nil {
		// something wrong, log, check
	}
	if valid && h.CheckCsrf(ctx, csrf) {
		strLocation := []byte("Location")
		storage.Delete("userAuth")
		storage.Delete("LoginPass")
		storage.Set("UserAuth", loginPass)
		ctx.Response.Header.SetCanonical(strLocation, []byte("/"))
		ctx.Response.SetStatusCode(fasthttp.StatusFound)
	} else {
		strLocation := []byte("Location")
		// TODO yeah, i know. don't mention
		ctx.Response.Header.Set("AuthFailed", "wrong credentials")
		ctx.Response.Header.SetCanonical(strLocation, []byte("/login"))
		ctx.Response.SetStatusCode(fasthttp.StatusFound)
	}
}

func (h *Auth) CheckCsrf(ctx *fasthttp.RequestCtx, token []byte) bool {
	storage := h.SessionStorage(ctx)
	t := storage.Get("csrf")
	if t == nil {
		return false
	}
	switch t.(type) {
	case []byte:
		return bytes.Equal(t.([]byte), token)
	case string:
		return util.Equals(t.(string), string(token))
	}
	return false
}

func (h *Auth) Session() *session.Session {
	return h.session
}

func (h *Auth) SessionStorage(ctx *fasthttp.RequestCtx) *session.Store {
	storage, _ := h.Session().Get(ctx)
	return storage
}

// Auth post
// TODO add throttle (middleware)
func (h *Auth) Auth(ctx *fasthttp.RequestCtx) {
	// login, password, remember, csrfToken
	login := ctx.PostArgs().Peek("username")
	pass := ctx.PostArgs().Peek("password")
	csrf := ctx.PostArgs().Peek("csrf")
	if util.NotEmpty(login) && util.NotEmpty(pass) && util.NotEmpty(csrf) &&
		h.userExists(login, pass) && h.CheckCsrf(ctx, csrf) {
		storage, _ := h.Session().Get(ctx)
		storage.Set("userAuth", util.RandStr(util.DefRndLen))
		storage.Set("LoginPass", fmt.Sprintf("%s:%s", login, pass))
		defer h.Session().Save(ctx, storage)
		strLocation := []byte("Location")
		ctx.Response.Header.SetCanonical(strLocation, []byte("/totp"))
		ctx.Response.SetStatusCode(fasthttp.StatusFound)
	} else {
		strLocation := []byte("Location")
		// TODO yeah, i know. don't mention
		ctx.Response.Header.Set("AuthFailed", "wrong credentials")
		ctx.Response.Header.SetCanonical(strLocation, []byte("/login"))
		ctx.Response.SetStatusCode(fasthttp.StatusFound)
	}
}

// Login GET,  Return form data
func (h *Auth) Login(ctx *fasthttp.RequestCtx) {
	sessionStorage, _ := h.session.Get(ctx)
	eblo := sessionStorage.Get("eblo")
	if eblo == nil {
		sessionStorage.Set("eblo", 1)
	} else {
		sessionStorage.Set("eblo", reflect.ValueOf(eblo).Int() + 1)
	}
	// TODO prob non-required when placed on middleware level
	defer h.Session().Save(ctx, sessionStorage)
	p := &templates.Login{CsrfToken: sessionStorage.Get("csrf").(string)}
	templates.WritePageTemplate(ctx, p)
}

func (h *Auth) Totp(ctx *fasthttp.RequestCtx) {
	sessionStorage, _ := h.session.Get(ctx)
	p := &templates.Totp{CsrfToken: sessionStorage.Get("csrf").(string)}
	templates.WritePageTemplate(ctx, p)
}
